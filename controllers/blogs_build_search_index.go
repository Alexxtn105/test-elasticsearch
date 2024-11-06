package controllers

import (
	"net/http"
	"test-elasticsearch/models"
)

func BlogsBuildSearchIndex(w http.ResponseWriter, r *http.Request) {
	// берем все блоги
	blogs := models.BlogsAll()

	// добавляем все блоги в индекс
	for _, blog := range *blogs {
		blog.AddToIndex()
	}
}
