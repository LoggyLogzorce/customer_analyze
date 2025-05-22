package api

import (
	"first_static_analiz/internal/handlers/api/analyze"
	"first_static_analiz/internal/handlers/api/generate_pdf"
	"first_static_analiz/internal/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"path/filepath"
)

type AnalysisRequest struct {
	Filename           string   `form:"filename" json:"filename"`
	Demografi          []string `form:"demografi[]" json:"demografi"`
	Finances           []string `form:"finances[]" json:"finances"`
	BehavioralAnalysis []string `form:"behavioral_analysis[]" json:"behavioral_analysis"`
	Visualizations     []string `form:"visualization[]" json:"visualizations"`
}

func AnalyzeHandler(c *gin.Context) {
	var req AnalysisRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request format"})
		return
	}

	fmt.Println(req)

	// 1. Чтение файла
	filePath := filepath.Join("storage", "uploaded_files", req.Filename)
	customers, code, err := readCustomerFile(filePath)
	if err != nil {
		c.JSON(code, gin.H{"error": fmt.Sprintf("Ошибка чтения файла: %v", err)})
		return
	}

	demResult := gin.H{}
	behResult := gin.H{}
	finResult := gin.H{}
	vizResult := gin.H{}

	// 2. Выполнение анализов
	if len(req.Demografi) > 0 {
		demResult = analyze.DemografiAnalyze(customers, req.Demografi)
	} else {
		demResult = nil
	}

	if len(req.BehavioralAnalysis) > 0 {
		behResult = analyze.CustomersAnalyze(customers, req.BehavioralAnalysis)
	} else {
		behResult = nil
	}

	if len(req.Finances) > 0 {
		finResult = analyze.FinanceAnalyze(customers, req.Finances)
	} else {
		finResult = nil
	}

	if len(req.Visualizations) > 0 {
		for _, v := range req.Visualizations {
			switch v {
			case "gender_pie":
				if demResult["gender_dist"] != nil {
					vizResult["gender_pie"] = demResult["gender_dist"]
				} else {
					vizResult = analyze.DemografiAnalyze(customers, req.Visualizations)
				}
			case "age_histogram":
				if vizResult["age_group"] == nil {
					vizResult = analyze.DemografiAnalyze(customers, req.Visualizations)
				}
			}
		}
	} else {
		vizResult = nil
	}

	result := gin.H{
		"filename": req.Filename,
		"status":   "analyzed",
	}

	if demResult != nil {
		result["demografi"] = demResult
	}

	if behResult != nil {
		result["behavioral_analysis"] = behResult
	}

	if finResult != nil {
		result["finances"] = finResult
	}

	if vizResult != nil {
		result["visualizations"] = vizResult
	}

	fmt.Println("Result:", result)

	c.JSON(http.StatusOK, result)
}

func GenerateReport(c *gin.Context) {
	var reportData model.ReportData
	if err := c.ShouldBindJSON(&reportData); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	bufPdf, err := generate_pdf.GeneratePDF(reportData)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.Data(200, "application/pdf; charset=utf-8", bufPdf.Bytes())
	//c.JSON(200, gin.H{})
}

func UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "No file uploaded"})
		return
	}

	newFileName := generateUniqueName(file.Filename)
	savePath := "./storage/uploaded_files/" + newFileName

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(500, gin.H{"error": "Failed to save file"})
		return
	}

	c.JSON(200, gin.H{
		"message":  "File uploaded successfully",
		"filename": newFileName,
	})
}

func generateUniqueName(original string) string {
	ext := filepath.Ext(original)
	return uuid.New().String() + ext
}
