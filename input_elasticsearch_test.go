package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestElasticSearchConfigRange(t *testing.T) {
	config := InputElasticSearchConfig{
		Address:   []string{"127.0.0.1", "127.0.0.2"},
		Index:     "cwl-raw-2021.07.03",
		FromDate:  time.Date(2021, 8, 3, 0, 0, 0, 0, time.Local),
		ToDate:    time.Date(2021, 8, 4, 0, 0, 0, 0, time.Local),
		Includes:  nil,
		Transport: nil,
	}

	assert.Equal(t, 60*24, config.Range())

}

func TestES(t *testing.T) {
	const layout = "cwl-raw-2006.01.02"
	quit := make(chan struct{})
	//date := time.Now().AddDate(0,-1,1)
	date := time.Now().AddDate(0, 0, -1)
	indexName := date.Format(layout)
	transport := &http.Transport{
		ResponseHeaderTimeout: time.Second * 5,
		MaxIdleConns:          10,
	}

	y, m, d := date.Date()
	config := InputElasticSearchConfig{
		Address:   []string{"https://vpc-sl-logstrg-orange-prd-q76s3uteh4ooxa3r4brwce2yau.ap-northeast-2.es.amazonaws.com"},
		Index:     indexName,
		FromDate:  time.Date(y, m, d, 0, 0, 0, 0, time.Local),
		ToDate:    time.Date(y, m, d, 0, 1, 0, 0, time.Local),
		Includes:  []string{},
		Transport: transport,
	}

	input := NewElasticsearchInput("", &config)
	for {
		msg, err := input.PluginRead()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(msg.Meta))
		fmt.Println(string(msg.Data))
		fmt.Println(strings.Repeat("-", 80))
	}

	_ = input
	<-quit

	log.Println("TestES end")

}
