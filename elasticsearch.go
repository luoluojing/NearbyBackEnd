package main

import (
	"context"

	"github.com/olivere/elastic/v7"
)

const (
	ES_URL = "http://10.138.0.2:9200"
)

func readFromES(query elastic.Query, index string) (*elastic.SearchResult, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(ES_URL),
		elastic.SetBasicAuth("luoluojing", "123456"))
	if err != nil {
		return nil, err
	}

	searchResult, err := client.Search().
		Index(index).            // search in index "twitter"
		Query(query).            // specify the query
		Pretty(true).            // pretty print request and response JSON
		Do(context.Background()) // execute
	if err != nil {
		return nil, err
	}

	return searchResult, nil
}

func saveToES(i interface{}, index string, id string) error {
	//为什么存interface而不是直接存post
	//saveToESz应该支持各种数据存储。
	//用Post就把参数传递写死了，无法更改。
	client, err := elastic.NewClient(
		elastic.SetURL(ES_URL),
		elastic.SetBasicAuth("luoluojing", "123456"))
	if err != nil {
		return err
	}

	_, err = client.Index().
		Index(index).
		Id(id).
		BodyJson(i).
		Do(context.Background())
	return err
}

func deleteFromES(query elastic.Query, index string) error {
	client, err := elastic.NewClient(
		elastic.SetURL(ES_URL),
		elastic.SetBasicAuth("luoluojing", "123456"))
	if err != nil {
		return err
	}

	_, err = client.DeleteByQuery().
		Index(index).
		Query(query).
		Pretty(true).
		Do(context.Background())

	return err
}
