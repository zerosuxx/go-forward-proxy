package main

import (
	"encoding/json"
	"flag"
	"github.com/elazarl/goproxy"
	"github.com/zerosuxx/go-escher-proxy/pkg/config"
	"github.com/zerosuxx/go-escher-proxy/pkg/handler"
	"io"
	"log"
	"net/http"
	"os"
)

type AppConfig struct {
	Hosts         map[string]HostConfig
	ListenAddress string
	TargetUrl     string
	Verbose       bool
}

type HostConfig struct {
	OverrideScheme 	bool
	OverrideHost 	bool
	TargetHost   	string
	UserAgent	 	string
}

func (appConfig *AppConfig) LoadFromArgument() {
	flag.StringVar(&appConfig.ListenAddress, "addr", "0.0.0.0:8282", "Proxy server listen address")
	flag.StringVar(&appConfig.TargetUrl, "url", "", "Target url")
	flag.BoolVar(&appConfig.Verbose, "v", false, "Verbose output")

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

	jsonData, _ := io.ReadAll(jsonFile)

	return jsonData
}

func patchRequest(request *http.Request, config *HostConfig) {
	if config.OverrideScheme {
		request.URL.Scheme = "https"
	}

	if config.OverrideHost {
		request.Header.Set("Host", request.Host)
	}

	if config.TargetHost != "" {
		request.URL.Host = config.TargetHost
	}

	if config.UserAgent != "" {
		request.Header.Set("User-Agent", config.UserAgent)
	}
}

var Version = "development"

func main() {
	const configFileName = "forward-proxy-config.json"

	appConfig := AppConfig{}
	appConfig.LoadFromArgument()
	path, _ := os.Getwd()
	configFile := path + "/" + configFileName
	appConfig.LoadFromJSONFile(configFile)

	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = appConfig.Verbose

	escherProxyConfig := config.AppConfig{}
	escherProxyConfig.Verbose = appConfig.Verbose

	proxy.NonproxyHandler = http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		webRequestHandler := handler.WebRequest{
			AppConfig: escherProxyConfig,
			Client:    &http.Client{},
		}

		if appConfig.TargetUrl != "" {
			request.Header.Set("X-Target-Url", appConfig.TargetUrl)
		}

		webRequestHandler.Handle(request, responseWriter)
	})
	proxy.OnRequest().DoFunc(func(request *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		host := request.URL.Host
		hostConfig := appConfig.FindHostConfig(host)

		if appConfig.Verbose {
			log.Println("Host", host)
			log.Println("Request", request)
			log.Println("Config", hostConfig)
		}

		if hostConfig != nil {
			patchRequest(request, hostConfig)

			return request, nil
		}

		if appConfig.Verbose {
			log.Println("Request is not patched", request)
		}

		return request, nil
	})

	log.Println("F0rward Pr0xy " + Version + " | Listening on: " + appConfig.ListenAddress)
	if appConfig.Verbose {
		log.Println("ConfigFile", configFile)
	}
	log.Fatalln(http.ListenAndServe(appConfig.ListenAddress, proxy))
}
