package routers

import (
	"first_static_analiz/internal/handlers/api"
	"first_static_analiz/internal/handlers/web"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	gin.SetMode(gin.DebugMode)
	r := gin.Default()

	r.LoadHTMLGlob("web/templates/*")
	r.Static("/static", "./web/static")

	webGroup := r.Group("/")
	{
		webGroup.GET("/", web.HomePage)
		webGroup.GET("/select-analysis/:filename", web.SelectAnalysis)
	}

	apiGroup := r.Group("/api/v1")
	{
		apiGroup.POST("/upload-file", api.UploadFile)
		apiGroup.POST("/analyze", api.AnalyzeHandler)
		apiGroup.POST("/generate-report", api.GenerateReport)
	}

	return r
}
