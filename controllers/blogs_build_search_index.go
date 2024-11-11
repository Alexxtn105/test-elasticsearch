package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"test-elasticsearch/models"
)

func BlogsBuildSearchIndex(c *gin.Context) {
	// берем все блоги
	blogs := models.BlogsAll()

	// добавляем все блоги в индекс
	for _, blog := range *blogs {
		blog.AddToIndex()
	}
}

func BlogsBuildSearchIndex_StabndardHttpClient(w http.ResponseWriter, r *http.Request) {
	// берем все блоги
	blogs := models.BlogsAll()

	// добавляем все блоги в индекс
	for _, blog := range *blogs {
		blog.AddToIndex()
	}
}
