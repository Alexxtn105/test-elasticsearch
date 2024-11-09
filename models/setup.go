package models

import (
	"context"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

var ESClient *elasticsearch.Client

const SearchIndex = "blogs"

func ConnectDatabase() {
	database, err := gorm.Open(sqlite.Open("storage/database.db"), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

	DB = database
}

func DBMigrate() {
	DB.AutoMigrate(&Blog{})
}

// ESClientConnection Создание подключения к Elasticsearch (порт по умолчанию 9200)
func ESClientConnection() {
	client, err := elasticsearch.NewDefaultClient()
	//caCert := []byte{0x71, 0xcd, 0x9e, 0x8d, 0xc0, 0x51, 0x05, 0x56, 0x3b, 0x35, 0x56, 0xc6, 0x35, 0x9e, 0xcd, 0x37, 0xb0, 0x25, 0xd0, 0x54, 0xe7, 0x00, 0xf8, 0x53, 0x54, 0xed, 0x5a, 0x6c, 0x43, 0xed, 0xb9, 0xb1}
	//cfg := elasticsearch.Config{
	//	Addresses: []string{
	//		"https://localhost:9200",
	//		"https://localhost:9300",
	//	},
	//	//CACert:   caCert,
	//	//Username: "elastic",
	//	//Password: "nOWlxL2hpKsVyD4LVF=2",
	//}
	//client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	res, _ := client.Info()

	log.Println(res)

	ESClient = client
}

// ESCreateIndexIfNotExist Создание индекса
func ESCreateIndexIfNotExist() {
	// сперва проверяем, существует ли индекс
	_, err := esapi.IndicesExistsRequest{
		Index: []string{SearchIndex},
	}.Do(context.Background(), ESClient)

	// Если возникает ошибка, значит индекса не существует. Его необходимо создать
	if err != nil {
		ESClient.Indices.Create(SearchIndex)
	}
}
