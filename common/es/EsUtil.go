package es

import (
	"XcxcVideo/common/define"
	"XcxcVideo/common/models"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"log"
	"math"
	"strconv"
)

var esClient = Init()

func Init() *elasticsearch.Client {
	// 初始化 Elasticsearch 客户端
	esClient, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9201", // 替换为你的 Elasticsearch 地址
		},
	})
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	// 检查 Elasticsearch 集群健康状态
	res, err := esClient.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	log.Println(res)
	return esClient
}
func AddUser(user models.UserVo) {
	indexName := "xcxc_user"
	res, err := esClient.Indices.Exists([]string{indexName})
	if err != nil {
		log.Fatalf("Error marshaling document: %s", err)
	}
	if res.StatusCode == 404 {
		// 如果索引不存在，则创建它
		res, err = esClient.Indices.Create(indexName)
		if err != nil {
			log.Fatalf("Error creating index: %s", err)
		}
		defer res.Body.Close()
		fmt.Printf("Index %s created\n", indexName)
	} else {
		fmt.Printf("Index %s already exists\n", indexName)
	}
	userJson, err := json.Marshal(user)
	if err != nil {
		log.Printf("Error marshaling document: %s", err)
	}
	_, err = esClient.Index(indexName, bytes.NewReader(userJson), esClient.Index.WithDocumentID(strconv.Itoa(user.Id)))
	if err != nil {
		log.Printf("Error marshaling document: %s", err)
	}

}
func AddSearchWord(content string) {
	searchWord := models.EsSearchWord{}
	searchWord.Content = content
	searchWordJson, err := json.Marshal(searchWord)
	if err != nil {
		log.Fatalf("Error marshaling document: %s", err)
	}
	_, err = esClient.Index(define.ES_SEARCH_WORD, bytes.NewReader(searchWordJson))
	if err != nil {
		log.Fatalf("Error indexing document: %s", err)
	}

}

func GetSearchWord(keyword string) []string {
	indexName := "xcxc_search_word"
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"should": []map[string]interface{}{
					{
						"match": map[string]interface{}{
							"content": keyword,
						},
					},
					{
						"match_phrase_prefix": map[string]interface{}{
							"content": keyword,
						},
					},
				},
			},
		},
	}

	// 将查询转换为 JSON
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	// 执行搜索请求
	res, err := esClient.Search(
		esClient.Search.WithContext(context.Background()),
		esClient.Search.WithIndex(indexName),
		esClient.Search.WithBody(&buf),
		esClient.Search.WithPretty(),
	)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	// 解析搜索结果
	if res.IsError() {
		log.Fatalf("Error response: %s", res.String())
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	searchWordList := []string{}
	// 打印匹配的结果和高亮部分
	fmt.Printf("Found %d hits:\n", int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)))
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		hitMap := hit.(map[string]interface{})
		source := hitMap["_source"].(map[string]interface{})
		searchWordList = append(searchWordList, source["content"].(string))
	}
	return searchWordList
}

func AddSearchVideo(videoVo models.VideoVo) {
	indexName := define.ES_VIDEO
	videoVoJson, err := json.Marshal(videoVo)
	if err != nil {
		log.Printf("Error marshaling document: %s", err)
	}
	_, err = esClient.Index(indexName, bytes.NewReader(videoVoJson), esClient.Index.WithDocumentID(strconv.Itoa(videoVo.Vid)))
	if err != nil {
		log.Printf("Error marshaling document: %s", err)
	}
}

func UpdateVideoStatus(vid int, status int) {
	indexName := "xcxc_video"
	updateBody := map[string]interface{}{
		"doc": map[string]interface{}{
			"status": status,
		},
	}
	searchWordJson, err := json.Marshal(updateBody)
	if err != nil {
		log.Printf("Error marshaling document: %s", err)
	}
	req := esapi.UpdateRequest{
		Index:      indexName,
		DocumentID: strconv.Itoa(vid),
		Body:       bytes.NewReader(searchWordJson),
		Refresh:    "true",
	}
	res, err := req.Do(context.Background(), esClient)
	if err != nil {
		log.Printf("Error marshaling document: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("Error marshaling document: %s", res.IsError())
	} else {
		var response map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
			log.Printf("Error marshaling document: %s", err)
		}
		log.Printf("Update response: %v", response)
	}
}
func GetSearchVideo(keyword string, pageNum int, pageSize int, onlyPass bool) []int {
	indexName := define.ES_VIDEO
	query := map[string]interface{}{
		"query": map[string]interface{}{
			//"fuzzy": map[string]interface{}{
			//	"title": map[string]interface{}{
			//		"value":     keyword, // 模糊查询的关键词
			//		"fuzziness": "AUTO",  // 自动模糊级别，支持 0、1、2 或 "AUTO"
			//	},
			//},
			"match": map[string]interface{}{
				"title": keyword,
			},
		},
		"highlight": map[string]interface{}{
			"fields": map[string]interface{}{
				"title": map[string]interface{}{},
			},
		},
	}

	// 将查询转换为 JSON
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	// 执行搜索请求
	res, err := esClient.Search(
		esClient.Search.WithContext(context.Background()),
		esClient.Search.WithIndex(indexName),
		esClient.Search.WithBody(&buf),
		esClient.Search.WithPretty(),
	)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	// 解析搜索结果
	if res.IsError() {
		log.Fatalf("Error response: %s", res.String())
	}

	var r struct {
		Hits struct {
			Total struct {
				Value int `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source    models.VideoVo      `json:"_source"`
				Highlight map[string][]string `json:"highlight"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}

	//videoVoList := make([]models.VideoVo, 0)
	fmt.Printf("Found %d hits:\n", r.Hits.Total.Value)
	ids := []int{}
	if onlyPass {
		for _, hit := range r.Hits.Hits {
			video := hit.Source
			if video.Status == 1 {
				ids = append(ids, video.Vid)
			}
		}
	} else {
		for _, hit := range r.Hits.Hits {
			video := hit.Source
			ids = append(ids, video.Vid)
		}
	}
	endIndex := math.Min(float64((pageNum+1)*pageSize), float64(len(ids)))
	ids = ids[(pageNum-1)*pageSize : int(endIndex)]
	//for _, hit := range r.Hits.Hits {
	//	video := hit.Source
	//if highlights, ok := hit.Highlight["title"]; ok && len(highlights) > 0 {
	//	video.Title = highlights[0] // 将第一个高亮片段存入结构体
	//}
	//videoVoList = append(videoVoList, video)

	//}

	return ids
}

func GetSearchVideoCount(keyword string) int {
	indexName := define.ES_VIDEO
	query := map[string]interface{}{
		"query": map[string]interface{}{
			//"fuzzy": map[string]interface{}{
			//	"title": map[string]interface{}{
			//		"value":     keyword, // 模糊查询的关键词
			//		"fuzziness": "AUTO",  // 自动模糊级别，支持 0、1、2 或 "AUTO"
			//	},
			//},
			"match": map[string]interface{}{
				"title": keyword,
			},
		},
		"highlight": map[string]interface{}{
			"fields": map[string]interface{}{
				"title": map[string]interface{}{},
			},
		},
	}

	// 将查询转换为 JSON
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	// 执行搜索请求
	res, err := esClient.Search(
		esClient.Search.WithContext(context.Background()),
		esClient.Search.WithIndex(indexName),
		esClient.Search.WithBody(&buf),
		esClient.Search.WithPretty(),
	)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	// 解析搜索结果
	if res.IsError() {
		log.Fatalf("Error response: %s", res.String())
	}

	var r struct {
		Hits struct {
			Total struct {
				Value int `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source    models.VideoVo      `json:"_source"`
				Highlight map[string][]string `json:"highlight"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}

	fmt.Printf("Found %d hits:\n", r.Hits.Total.Value)
	return r.Hits.Total.Value
}
