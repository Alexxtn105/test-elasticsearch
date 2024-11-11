package models

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"gorm.io/gorm"
	"log"
	"strconv"
)

type Blog struct {
	gorm.Model
	Title   string `gorm:"size:255"`
	Content string `gorm:"type:text"`
}

func BlogsAll() *[]Blog {
	var blogs []Blog
	DB.Where("deleted_at is NULL").Order("updated_at desc").Find(&blogs)
	return &blogs
}

func BlogsFind(id uint64) *Blog {
	var blog Blog
	DB.Where("id = ?", id).First(&blog)
	return &blog
}

func (b *Blog) AddToIndex() {
	//fmt.Printf("%#v\n", b)
	// Создаем структуру документа, которая будет добавлена в индекс
	// Эта структура содержит текстовые поля модели Blog с соответствующими данными
	document := struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}{
		b.Title,
		b.Content,
	}

	fmt.Println(document.Title)
	fmt.Println(document.Content)

	// преобразуем документ в JSON
	data, err := json.Marshal(document)
	fmt.Println(data)
	if err != nil {
		log.Fatalf("ошибка маршалинга документа: %s", err)
	}

	// создаем новый запрос
	req := esapi.IndexRequest{
		Index:      SearchIndex,
		DocumentID: strconv.Itoa(int(b.ID)), // ИД документа в ElasticSearch должен быть типа string
		Body:       bytes.NewReader(data),
		Refresh:    "true", // обновление будет происходить немедленно после операции, документ сразу будет доступен для поиска
	}

	//fmt.Println(b.ID)
	//fmt.Println(req)
	//fmt.Println(ESClient)

	//выполняем запрос
	resp, err := req.Do(context.Background(), ESClient)
	if err != nil {
		log.Fatalf("ошибка выполнения запроса: %s", err)
	}

	// обязательно закрываем тело запроса
	defer resp.Body.Close()

	log.Printf("Индексирован документ %s в индекс %s\n", resp.String(), SearchIndex)
}

// BlogSearch Возвращает результат запроса по поиску
func BlogSearch(searchQuery string) *[]Blog {
	var buf bytes.Buffer

	// Создаем понятный для ElasticSearch запрос (это специально сформированная хеш-таблица)
	query := map[string]any{
		"query": map[string]any{ // query - мапа, специфичная для языка запросов elasticsearch
			"multi_match": map[string]any{ // multi_match - специальный запрос по поиску в нескольких полях
				"query":  searchQuery,                  // это ключевое поле для поиска (что ищем)
				"fields": []string{"title", "content"}, // тут описываем поля, в которых искать (где конкретно ищем)
			},
		},
	}

	// Преобразуем запрос в JSON, который понимает ElasticSearch
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	// Выполняем запрос
	res, err := ESClient.Search(
		ESClient.Search.WithIndex(SearchIndex),
		ESClient.Search.WithBody(&buf),
	)

	// Обязательно закрываем тело ответа
	defer res.Body.Close()

	if err != nil || res.IsError() {
		return nil
	}

	// декодируем результат поиска из JSON в мапу
	var r map[string]any
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil
	}

	// извлекаем ИД записей, удовлетворяющих условию поиска из тела запроса
	var ids []uint
	if hits, ok := r["hits"].(map[string]any); ok {
		if hitsHits, ok := hits["hits"].([]any); ok {
			for _, hit := range hitsHits {
				if hitMap, ok := hit.(map[string]any); ok {
					if idStr, ok := hitMap["_id"].(string); ok {
						id, _ := strconv.Atoi(idStr)
						ids = append(ids, uint(id))
					}
				}
			}
		}
	}

	// по извлеченным из ElasticSearch ИД находим соответствующие записи в БД
	var blogs []Blog

	DB.Where("deleted_at IS NULL").Where("id in ?", ids).Order("updated_at desc").Find(&blogs)

	return &blogs
}
