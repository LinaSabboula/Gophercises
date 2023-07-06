package parser

import (
	"encoding/json"
	"log"
)

type Chapter struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []struct {
		Text string `json:"text"`
		Arc  string `json:"arc"`
	} `json:"options"`
}

func ParseData(data []byte) (map[string]Chapter, error) {
	var parsedData map[string]Chapter
	err := json.Unmarshal(data, &parsedData)
	if err != nil {
		log.Fatal(err)
	}
	return parsedData, err
}
