package urlshort

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"net/http"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {

	return func(writer http.ResponseWriter, request *http.Request) {
		url, ok := pathsToUrls[request.URL.Path]
		if ok {
			http.Redirect(writer, request, url, http.StatusPermanentRedirect)
		} else {
			fallback.ServeHTTP(writer, request)
		}
	}
}

//
//type YAMLInput struct {
//	YAML struct {
//		YAML []YAML `yaml:"path,url"`
//	}
//}

type YAML struct {
	Path string `yaml:"path"`
	Url  string `yaml:"url"`
}

type JSON struct {
	path string
	url  string
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}
	pathMap := buildYamlMap(parsedYaml)
	return MapHandler(pathMap, fallback), nil
}
func JSONHandler(json []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedJson, err := parseJSON(json)
	if err != nil {
		return nil, err
	}
	pathMap := buildJsonMap(parsedJson)
	return MapHandler(pathMap, fallback), nil
}
func parseJSON(decodedJson []byte) ([]JSON, error) {
	var parsedJson []JSON
	err := json.Unmarshal(decodedJson, &parsedJson)
	return parsedJson, err
}
func parseYAML(yml []byte) ([]YAML, error) {
	var parsedYAML []YAML
	err := yaml.Unmarshal(yml, &parsedYAML)
	return parsedYAML, err
}
func buildYamlMap(parsedYAML []YAML) map[string]string {
	yamlMap := make(map[string]string)
	for _, v := range parsedYAML {
		yamlMap[v.Path] = v.Url
	}
	return yamlMap
}

func buildJsonMap(parsedJSON []JSON) map[string]string {
	yamlMap := make(map[string]string)
	for _, v := range parsedJSON {
		yamlMap[v.path] = v.url
	}
	return yamlMap
}
