package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"test-elasticsearch/controllers"
	"test-elasticsearch/models"
)

var addr = ":8087"

func main() {
	// Адрес клиента
	r := gin.Default()
	r.Use(gin.Logger())

	// подгружаем шаблоны страниц
	r.LoadHTMLGlob("templates/**/*")

	// Создаем пустую БД
	models.ConnectDatabase()
	models.DBMigrate()

	// Подключаем ElasticSearch
	models.ESClientConnection()

	// Создаем индекс в ElasticSearch
	models.ESCreateIndexIfNotExist()

	r.GET("/blogs", controllers.BlogsIndex)
	r.GET("/blogs/:id", controllers.BlogsShow)
	r.POST("/blogs/index", controllers.BlogsBuildSearchIndex)

	log.Printf("server is listening at %s", addr)
	r.Run(addr)
	//region //Стандартный роутер
	//// Создаем роутер
	//mux := http.NewServeMux()
	//
	//// В хендлере будем использовать параметры пути ({id}). Введены в go 1.22.
	//// Можно посмотреть видео: https://www.youtube.com/watch?v=H7tbjKFSg58&t=48s
	////mux.HandleFunc("GET /blogs", controllers.BlogsIndex)
	//mux.HandleFunc("GET /blogs/{id}", controllers.BlogsShow)
	//mux.HandleFunc("POST /blogs/index", controllers.BlogsBuildSearchIndex)
	//
	//// Запускаем сервер
	//log.Printf("server is listening at %s", addr)
	//err := http.ListenAndServe(addr, mux)
	//if err != nil {
	//	log.Fatal()
	//}
	//endregion
}
