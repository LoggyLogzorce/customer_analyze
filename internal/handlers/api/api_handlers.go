package api

import (
	"bytes"
	"encoding/base64"
	"first_static_analiz/internal/handlers/api/analyze"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jung-kurt/gofpdf"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	pageWidth    = 210 // Ширина A4 в мм
	pageHeight   = 297 // Высота A4 в мм
	marginTop    = 10  // Верхний отступ
	marginBottom = 20  // Нижний отступ (место для колонтитула)
	marginLeft   = 15  // Левый отступ
)

type ReportData struct {
	Meta struct {
		ReportID    string    `json:"report_id"`
		GeneratedAt time.Time `json:"generated_at"`
		Filename    string    `json:"filename"`
		TotalUsers  uint      `json:"total_users"`
	} `json:"meta"`

	Demografi struct {
		GenderDistribution map[string]int `json:"gender_distribution"`
		AgeGroups          map[string]int `json:"age_groups"`
		TopRegions         []struct {
			Name  string `json:"name"`
			Count int    `json:"count"`
		} `json:"top_regions"`
	} `json:"demografi"`

	Behavioral struct {
		Veterans int `json:"veterans"`
		Newbies  int `json:"newbies"`
		VIPs     struct {
			Count      int     `json:"count"`
			Percentile float64 `json:"percentile"`
		} `json:"vips"`
	} `json:"behavioral"`

	//Financials struct {
	//	Income struct {
	//		Mean   float64 `json:"mean"`
	//		Median float64 `json:"median"`
	//		Mode   float64 `json:"mode"`
	//		StdDev float64 `json:"std_dev"`
	//	} `json:"income"`
	//	Spending struct {
	//		Mean   float64 `json:"mean"`
	//		Median float64 `json:"median"`
	//		Mode   float64 `json:"mode"`
	//		StdDev float64 `json:"std_dev"`
	//	} `json:"spending"`
	//	CheckSizeDistribution map[string]int `json:"check_size_distribution"`
	//} `json:"financials"`

	Visualizations struct {
		GenderPie    string `json:"gender_pie"`    // Base64 encoded PNG
		AgeHistogram string `json:"age_histogram"` // Base64 encoded PNG
		// IncomeVsSpendingScatter string `json:"income_vs_spending_scatter"`
	} `json:"visualizations"`

	//Additional struct {
	//	Notes       string `json:"notes"`
	//	AnalystName string `json:"analyst_name"`
	//	Segments    []struct {
	//		Name        string `json:"name"`
	//		Description string `json:"description"`
	//		Count       int    `json:"count"`
	//	} `json:"segments"`
	//} `json:"additional"`
}

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

	// 1. Чтение файла
	filePath := filepath.Join("storage", "uploaded_files", req.Filename)
	customers, err := readCustomerFile(filePath)
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("failed to read file: %v", err)})
		return
	}

	demResult := gin.H{}
	behResult := gin.H{}
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
	// [gender_pie, age_histogram] len=2 len
	if len(req.Visualizations) > 0 {
		//vizResult = analyze.DemografiAnalyze(customers, req.Visualizations)
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

	//if len(req.Finances) > 0 {
	//	fulResult = analyzeDemografi(customers, req.Demografi)
	//}

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

	if vizResult != nil {
		result["visualizations"] = vizResult
	}

	fmt.Println(result)

	c.JSON(http.StatusOK, result)
}

func GeneratePDF(c *gin.Context) {
	var reportData ReportData
	if err := c.ShouldBindJSON(&reportData); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.AddUTF8Font("dejavu", "", "./storage/fonts/dejavu-fonts-ttf-2.37/ttf/DejaVuSansCondensed.ttf")
	pdf.SetFont("dejavu", "", 20)

	// Текущая позиция Y
	yPos := 10.0

	// Заголовок с кириллицей
	pdf.SetY(yPos)
	pdf.Cell(40, 10, "Аналитический отчет")
	yPos += 15

	// Meta данные
	pdf.SetFont("dejavu", "", 12)
	pdf.SetY(yPos)
	pdf.Cell(40, 10, fmt.Sprintf("ReportID: %s", reportData.Meta.ReportID))
	yPos += 7

	pdf.SetY(yPos)
	pdf.Cell(40, 10, fmt.Sprintf("Сгенерирован: %s", reportData.Meta.GeneratedAt))
	yPos += 7

	pdf.SetY(yPos)
	pdf.Cell(40, 10, fmt.Sprintf("Файл: %s", reportData.Meta.Filename))
	yPos += 7

	if reportData.Meta.TotalUsers != 0 {
		pdf.SetY(yPos)
		pdf.Cell(40, 10, fmt.Sprintf("Всего клиентов: %d", reportData.Meta.TotalUsers))
		yPos += 7
	}

	yPos += 5

	// Демография
	pdf.SetFont("dejavu", "", 16)
	pdf.SetY(yPos)
	pdf.Cell(40, 10, "Демография")
	yPos += 10

	pdf.SetFont("dejavu", "", 12)
	pdf.SetY(yPos)
	pdf.Cell(40, 10, "Распределение по полу:")
	yPos += 7
	for k, v := range reportData.Demografi.GenderDistribution {
		pdf.SetY(yPos)
		pdf.Cell(40, 10, fmt.Sprintf("%s: %d", k, v))
		yPos += 7
	}

	yPos += 5

	pdf.SetY(yPos)
	pdf.Cell(40, 10, "Возрастные группы:")
	yPos += 7
	for k, v := range reportData.Demografi.AgeGroups {
		if k == "Count" {
			continue
		}
		pdf.SetY(yPos)
		pdf.Cell(40, 10, fmt.Sprintf("%s: %d", k, v))
		yPos += 7
	}

	yPos += 5

	pdf.SetY(yPos)
	pdf.Cell(40, 10, "Топ регионов:")
	yPos += 7
	for _, v := range reportData.Demografi.TopRegions {
		pdf.SetY(yPos)
		pdf.Cell(40, 10, fmt.Sprintf("%s: %d", v.Name, v.Count))
		yPos += 7
	}

	yPos += 5

	// Поведенческий анализ
	if reportData.Behavioral.Veterans != 0 || reportData.Behavioral.Newbies != 0 || reportData.Behavioral.VIPs.Count != 0 {
		pdf.SetFont("dejavu", "", 16)
		pdf.SetY(yPos)
		pdf.Cell(40, 10, "Поведенческий анализ")
		yPos += 10

		pdf.SetFont("dejavu", "", 12)
		if reportData.Behavioral.Veterans != 0 {
			pdf.SetY(yPos)
			pdf.Cell(40, 10, fmt.Sprintf("Ветераны: %d (регистрация до 2023 года)",
				reportData.Behavioral.Veterans))
			yPos += 7
		}
		if reportData.Behavioral.Newbies != 0 {
			pdf.SetY(yPos)
			pdf.Cell(40, 10, fmt.Sprintf("Новички: %d (регистрация в 2025 года)",
				reportData.Behavioral.Newbies))
			yPos += 7
		}
		if reportData.Behavioral.VIPs.Count != 0 {
			pdf.SetY(yPos)
			pdf.Cell(40, 10, fmt.Sprintf("VIP-клиенты: %d (перцентиль: %v) (средний чек выше 75-го перцентиля)",
				reportData.Behavioral.VIPs.Count, reportData.Behavioral.VIPs.Percentile))
			yPos += 7
		}
	}

	// Изображения
	yPos = checkPageBreak(pdf, yPos, 60)
	if err := addSafeImage(pdf, reportData.Visualizations.GenderPie, marginLeft, yPos); err == nil {
		yPos += 130 // Отступ после диаграммы
	}
	// Вторая диаграмма (возрастные группы)
	yPos = checkPageBreak(pdf, yPos, 60)
	if err := addSafeImage(pdf, reportData.Visualizations.AgeHistogram, marginLeft, yPos); err == nil {
		yPos += 130 // Отступ после диаграммы
	}

	// Сохранение
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.Data(200, "application/pdf; charset=utf-8", buf.Bytes())
	//c.JSON(200, gin.H{})
}

func addSafeImage(pdf *gofpdf.Fpdf, base64Image string, x, y float64) error {
	// Проверка пустой строки
	if base64Image == "" {
		return fmt.Errorf("empty base64 string")
	}

	// Очистка Base64
	cleanedBase64 := cleanBase64(base64Image)

	// Проверка длины (минимальная длина валидного Base64 - 4 символа)
	if len(cleanedBase64) < 4 {
		return fmt.Errorf("base64 string too short")
	}

	// Декодирование
	imgData, err := base64.StdEncoding.DecodeString(cleanedBase64)
	if err != nil {
		return fmt.Errorf("base64 decode error: %v", err)
	}

	// Определение формата изображения
	contentType := http.DetectContentType(imgData)
	var tempExt string
	switch contentType {
	case "image/png":
		tempExt = "*.png"
	case "image/jpeg":
		tempExt = "*.jpg"
	default:
		return fmt.Errorf("unsupported image format: %s", contentType)
	}

	// Создание временного файла
	tempFile, err := os.CreateTemp("./storage/reports", tempExt)
	if err != nil {
		return fmt.Errorf("temp file creation error: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Запись изображения
	if _, err := tempFile.Write(imgData); err != nil {
		return fmt.Errorf("file write error: %v", err)
	}

	// Добавление в PDF
	pdf.ImageOptions(
		tempFile.Name(),
		x,
		y,
		120,
		0,
		false,
		gofpdf.ImageOptions{ImageType: strings.TrimPrefix(tempExt, "*."), ReadDpi: true},
		0,
		"",
	)

	return nil
}

func checkPageBreak(pdf *gofpdf.Fpdf, currentY, elementHeight float64) float64 {
	var availableHeight float64
	availableHeight = pageHeight - marginTop - marginBottom
	if currentY+elementHeight > availableHeight {
		pdf.AddPage()
		currentY = marginTop
		return currentY // Возвращаем новую Y-позицию
	}
	return currentY
}

func cleanBase64(data string) string {
	// Удаляем возможный префикс data:image/png;base64,
	parts := strings.SplitN(data, ",", 2)
	if len(parts) == 2 {
		return parts[1]
	}
	return data
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
