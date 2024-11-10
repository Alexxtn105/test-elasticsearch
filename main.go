package main

import (
	"log"
	"net/http"
	"test-elasticsearch/controllers"
	"test-elasticsearch/models"
)

func main() {
	// Адрес клиента
	addr := ":8087"

	// Создаем пустую БД
	models.ConnectDatabase()
	models.DBMigrate()

	// Подключаем ElasticSearch
	models.ESClientConnection()

	// Создаем индекс в ElasticSearch
	models.ESCreateIndexIfNotExist()

	// Создаем роутер
	mux := http.NewServeMux()

	// В хендлере будем использовать параметры пути ({id}). Введены в go 1.22.
	// Можно посмотреть видео: https://www.youtube.com/watch?v=H7tbjKFSg58&t=48s
	mux.HandleFunc("GET /blogs/{id}", controllers.BlogsShow)
	mux.HandleFunc("POST /blogs/index", controllers.BlogsBuildSearchIndex)

	// Запускаем сервер
	log.Printf("server is listening at %s", addr)
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		log.Fatal()
	}
}
