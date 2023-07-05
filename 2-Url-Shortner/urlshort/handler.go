package urlshort

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"log"
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

type Data struct {
	Path string
	Url  string
}

func DataHandler(data []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedData, err := parseData(data)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedData)
	return MapHandler(pathMap, fallback), nil
}
func DBHandler(db *bolt.DB, fallback http.Handler) http.HandlerFunc {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	return func(writer http.ResponseWriter, request *http.Request) {
		bucket := viper.Get("BUCKET").(string)
		err = db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(bucket))
			value := b.Get([]byte(request.URL.Path))
			http.Redirect(writer, request, string(value), http.StatusPermanentRedirect)
			return nil
		})
		if err != nil {
			fallback.ServeHTTP(writer, request)
		}
	}
}

func parseData(data []byte) ([]Data, error) {
	var parsedData []Data
	err := json.Unmarshal(data, &parsedData)
	if err != nil {
		err = yaml.Unmarshal(data, &parsedData)
	}
	return parsedData, err
}

func buildMap(parsedData []Data) map[string]string {
	dataMap := make(map[string]string)
	for _, v := range parsedData {
		dataMap[v.Path] = v.Url
	}
	return dataMap
}
