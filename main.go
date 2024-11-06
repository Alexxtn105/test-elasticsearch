package main

import (
	"log"
	"net/http"
	"test-elasticsearch/controllers"
	"test-elasticsearch/models"
)

func main() {
	//fmt.Println("Hello, World!")
	addr := ":8087"

	//создаем пустую БД
	models.ConnectDatabase()
	models.DBMigrate()

	//подключаем ElasticSearch
	models.ESClientConnection()
	//создаем индекс в ElasticSearch
	models.ESCreateIndexIfNotExist()

	//создаем роутер
	mux := http.NewServeMux()

	// В хендлере будем использовать параметры пути ({id}). Введены в go 1.22.
	// Можно посмотреть видео: https://www.youtube.com/watch?v=H7tbjKFSg58&t=48s
	mux.HandleFunc("GET /blogs/{id}", controllers.BlogsShow)
	mux.HandleFunc("POST /blogs/index", controllers.BlogsBuildSearchIndex)

	//запускаем сервер
	log.Printf("server is listening at %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
