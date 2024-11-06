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
