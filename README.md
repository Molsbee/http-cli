# http-cli
Simple CLI for performing HTTP requests it supports the GET, PUT, POST, DELETE, HEAD methods.  
It also supports executing a yaml file which contains a list of requests that should be performed.  
The requests are executed in order, and it supports parsing fields of the response to be used in subsequent requests.
The parse value is a jq syntax path.

Example YAML file 
```yaml
requests:
  - name: "get countries"
    method: GET
    url: "https://restcountries.eu/rest/v2/all"
    parse:
      countryName: ".[0].name"
  - name: "get country"
    method: GET
    url: "https://restcountries.eu/rest/v2/name/{{ .countryName }}"
```
