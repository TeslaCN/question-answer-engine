package es

import (
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v7"
	"log"
	"question-answer-engine/cmd/qa/app/config"
	"question-answer-engine/internal/app/qa/bo"
)

var Client *elasticsearch.Client

func init() {
	cfg := elasticsearch.Config{
		Addresses: config.Config.Elasticsearch.Hosts,
	}
	var err error
	Client, err = elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	// Get cluster info
	var r map[string]interface{}
	res, err := Client.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	// Check response status
	if res.IsError() {
		log.Fatalf("Error: %s", res.String())
	}
	// Deserialize the response into a map.
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	// Print client and server version numbers.
	log.Printf("Client: %s", elasticsearch.Version)
	log.Printf("Server: %s", r["version"].(map[string]interface{})["number"])
}

type IndexResponse struct {
	Index         string `json:"_index"`
	Type          string `json:"_type"`
	Id            string `json:"_id"`
	Version       int    `json:"_version"`
	Result        string `json:"result"`
	ForcedRefresh bool   `json:"forced_refresh"`
	Shards        struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
	SeqNo       uint64 `json:"_seq_no"`
	PrimaryTerm int    `json:"_primary_term"`
}

type QaQueryResponse struct {
	Took     uint64 `json:"took"`
	TimedOut bool   `json:"timed_out"`
	Hits     *struct {
		Total struct {
			Value    uint64 `json:"value"`
			Relation string `json:"relation"`
		} `json:"total"`
		MaxScore float64 `json:"max_score"`
		Hits     []*struct {
			Index     string             `json:"_index"`
			Type      string             `json:"_type"`
			Id        string             `json:"_id"`
			Score     float64            `json:"_score"`
			Source    *bo.QuestionAnswer `json:"_source"`
			Highlight *struct {
				Question []string `json:"question"`
				Answer   []string `json:"answer"`
				Tag      []string `json:"tag"`
			} `json:"highlight"`
		} `json:"hits"`
	} `json:"hits"`
}
