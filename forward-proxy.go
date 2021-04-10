package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/elazarl/goproxy"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

type AppConfig struct {
	Hosts         map[string]HostConfig
	ListenAddress string
	Verbose       bool
}

type HostConfig struct {
	OverrideHost bool
	TargetHost   string
}

func (appConfig *AppConfig) LoadFromArgument() {
	flag.StringVar(&appConfig.ListenAddress, "addr", "0.0.0.0:8282", "Proxy server listen address")
	flag.BoolVar(&appConfig.Verbose, "v", false, "Verbose")

	flag.Parse()
}

func (appConfig *AppConfig) LoadFromJSONFile(jsonFile string) {
	if _, err := os.Stat(jsonFile); err == nil {
		jsonData := readFromFile(jsonFile)
		jsonError := json.Unmarshal(jsonData, appConfig)

		if jsonError != nil {
			log.Println(jsonError)
		}
	}
}

func (appConfig *AppConfig) FindHostConfig(host string) *HostConfig {
	if appConfig.Hosts == nil {
		return nil
	}

	if val, exists := appConfig.Hosts[host]; exists {
		return &val
	}

	return nil
}

func readFromFile(file string) []byte {
	jsonFile, err := os.Open(file)

	if err != nil {
		return []byte("")
	}

	jsonData, _ := ioutil.ReadAll(jsonFile)

	return jsonData
}

func detectPort(url url.URL) int {
	port := url.Port()

	if port != "" {
		portNumber, _ := strconv.Atoi(port)

		return portNumber
	}

	if url.Scheme == "https" {
		return 443
	}

	return 80
}

func patchRequest(request *http.Request, config *HostConfig) {
	if config.OverrideHost {
		request.Header.Set("Host", request.Host)
	}

	if config.TargetHost != "" {
		request.URL.Host = config.TargetHost
	}
}

func main() {
	const VERSION = "0.0.1"
	const configFileName = "forward-proxy-config.json"

	appConfig := AppConfig{}
	appConfig.LoadFromArgument()
	path, _ := os.Getwd()
	configFile := path + "/" + configFileName
	appConfig.LoadFromJSONFile(configFile)

	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = appConfig.Verbose

	proxy.OnRequest().DoFunc(func(request *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		hostConfig := appConfig.FindHostConfig(request.URL.Host)

		if hostConfig != nil {
			patchRequest(request, hostConfig)

			if appConfig.Verbose {
				log.Println("Host", request.URL.Host)
				log.Println("Request", request)
				log.Println("Config", hostConfig)
			}

			return request, nil
		}

		if request.URL.Port() == "" {
			port := detectPort(*request.URL)
			hostWithPort := fmt.Sprintf("%s:%d", request.URL.Host, port)
			hostConfig = appConfig.FindHostConfig(hostWithPort)

			if hostConfig != nil {
				patchRequest(request, hostConfig)

				if appConfig.Verbose {
					log.Println("Host", hostWithPort)
					log.Println("Request", request)
					log.Println("Config", hostConfig)
				}

				return request, nil
			}
		}

		if appConfig.Verbose {
			log.Println("Original Request", request)
		}

		return request, nil
	})

	log.Println("F0rward Pr0xy " + VERSION + " | Listening on: " + appConfig.ListenAddress)
	log.Fatal(http.ListenAndServe(appConfig.ListenAddress, proxy))
}
