package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"test-elasticsearch/models"
)

func BlogsShow(w http.ResponseWriter, r *http.Request) {

	// Извлекаем ИД из параметров пути
	//idStr := r.URL.Path[len("/blogs/"):] // Старый вариант  (до версии 1.22)
	idStr := r.PathValue("id") //сразу берем параметр из пути (стало доступно в go версии 1.22)
	// конвертим строку в int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Ivalid blog ID", http.StatusBadRequest)
		return
	}
	//log.Println(idStr)

	blog := models.BlogsFind(uint64(id))

	// преобразуем в слайс байтов
	blogBytes, _ := json.Marshal(blog)

	// имитируем длительный процесс
	//time.Sleep(time.Second * 2)

	// посылаем ответ пользователю
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(blogBytes); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}

}
