package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"io/ioutil"
	"log"
	"question-answer-engine/cmd/qa/app/config"
	"question-answer-engine/internal/app/qa/es"
	"strings"
)

var index = config.Config.QaIndex

func Put(o interface{}) string {
	b, e := json.Marshal(o)
	if e != nil {
		log.Println(e)
	}
	request := esapi.IndexRequest{
		Index:   index,
		Body:    bytes.NewReader(b),
		Refresh: "true",
	}
	response, e := request.Do(context.Background(), es.Client)
	if e != nil {
		log.Println(e)
	}
	b, _ = ioutil.ReadAll(response.Body)
	indexResponse := &es.IndexResponse{}
	_ = json.Unmarshal(b, indexResponse)
	log.Printf("%s\n", string(b))
	return indexResponse.Id
}

func MultiMatchQa(q string, fieldsWeight map[string]float64, size int, preTag string, postTag string) interface{} {
	fields := []string{"question", "answer", "tag"}
	builder := &strings.Builder{}
	for i, field := range fields {
		fieldTemplate := `"%s^%f"`
		weight := 1.0
		if w, exist := fieldsWeight[field]; exist {
			weight = w
		}
		builder.WriteString(fmt.Sprintf(fieldTemplate, field, weight))
		if i < len(fields)-1 {
			builder.WriteString(", ")
		}
	}
	formattedFields := builder.String()
	requestBodyTemplate := `{
  "query": {
    "multi_match": {
      "fields": [%s],
      "query": "%s"
    }
  },
  "highlight": {
    "pre_tags": "%s",
    "post_tags": "%s",
    "fields": {
      "question": {},
      "answer": {}
    }
  },
  "size": %d
}
`
	requestBody := fmt.Sprintf(requestBodyTemplate, formattedFields, q, preTag, postTag, size)
	response, e := es.Client.Search(
		es.Client.Search.WithIndex(index),
		es.Client.Search.WithContext(context.Background()),
		es.Client.Search.WithBody(strings.NewReader(requestBody)),
	)
	if e != nil {
		log.Println(e)
	}
	b, e := ioutil.ReadAll(response.Body)
	qaQueryResponse := &es.QaQueryResponse{}
	_ = json.Unmarshal(b, qaQueryResponse)
	hits := qaQueryResponse.Hits.Hits
	log.Println(string(b))
	return hits
}
