package main

import (
	"mime/multipart"
	"reflect"

	"github.com/olivere/elastic/v7"
)

// POJO
// '' 代表raw string
// 首字母大写是public, 首字母小写是private
const (
	POST_INDEX = "post"
)

type Post struct {
	Id      string `json:"id"`
	User    string `json:"user"`
	Message string `json:"message"`
	Url     string `json:"url"`
	Type    string `json:"type"`
}

func searchPostsByUser(user string) ([]Post, error) {
	query := elastic.NewTermQuery("user", user)
	searchResult, err := readFromES(query, POST_INDEX)
	if err != nil {
		return nil, err
	}

	return getPostFromSearchResult(searchResult), nil
}

func searchPostsByKeywords(keywords string) ([]Post, error) {
	query := elastic.NewMatchQuery("message", keywords)
	//多个关键词匹配成功才返回。
	//不提供关键字返回所有。
	query.Operator("AND")
	if keywords == "" {
		query.ZeroTermsQuery("all")
	}

	searchResult, err := readFromES(query, POST_INDEX)
	if err != nil {
		return nil, err
	}
	return getPostFromSearchResult(searchResult), nil
}

func getPostFromSearchResult(searchResult *elastic.SearchResult) []Post {
	var ptype Post
	var posts []Post

	for _, item := range searchResult.Each(reflect.TypeOf(ptype)) {
		p := item.(Post)
		posts = append(posts, p)
	}
	return posts
}

//9200端口不经过go程序，直接进入elastic读取。
//8080端口多了节点，postman query go程序读取elastic -> 因为1. 前端处理所有请求就慢。 2. 前端读取数据库暴露了一些信息，

func savePost(post *Post, file multipart.File) error {
	//use pointer 为了减少开销。
	mediaLink, err := saveToGCS(file, post.Id)
	if err != nil {
		return err
	}
	post.Url = mediaLink

	return saveToES(post, POST_INDEX, post.Id)
}

func deletePost(id string, user string) error {
	query := elastic.NewBoolQuery()
	query.Must(elastic.NewTermQuery("id", id))
	query.Must(elastic.NewTermQuery("user", user))

	return deleteFromES(query, POST_INDEX)
}
