package web

import (
	"github.com/gin-gonic/gin"
	"os"
)

func HomePage(c *gin.Context) {
	c.HTML(200, "index.html", gin.H{
		"title": "Анализ клиентской базы",
	})
}

func SelectAnalysis(c *gin.Context) {
	fileName := c.Param("filename")
	filePath := "./storage/uploaded_files/" + fileName

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(400, gin.H{
			"error": "file is not exist",
		})
		return
	}

	c.HTML(200, "select_analysis.html", gin.H{
		"title":    "Выбор категорий анализа",
		"filename": fileName,
	})
}
