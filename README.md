# go-forward-proxy

[![CI](https://github.com/zerosuxx/go-forward-proxy/workflows/CI/badge.svg)](https://github.com/zerosuxx/go-forward-proxy/actions?query=workflow%3ACI)

## Show available arguments
```
forward-proxy -h
```

## Install
```
make build
```

## Run
```
make run
```

## Usage

### Example Config (forward-proxy-config.json)
```
{
  "hosts": {
    "my-project.local": {
      "overrideHost": true,
      "targetHost": "localhost:8080"
    }
  }
}
```

### Example curl command
```
curl -x localhost:8282 http://my-project.local # forwarded to http://localhost:8080 (Host: my-project.local)
```
