package generate_pdf

import (
	"bytes"
	"encoding/base64"
	"first_static_analiz/internal/model"
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"net/http"
	"os"
	"strings"
)

const (
	pageWidth    = 210 // Ширина A4 в мм
	pageHeight   = 297 // Высота A4 в мм
	marginTop    = 10  // Верхний отступ
	marginBottom = 20  // Нижний отступ (место для колонтитула)
	marginLeft   = 15  // Левый отступ
)

func GeneratePDF(reportData model.ReportData) (bytes.Buffer, error) {
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
	if len(reportData.Demografi.GenderDistribution) != 0 || len(reportData.Demografi.AgeGroups) != 0 || len(reportData.Demografi.TopRegions) != 0 {
		pdf.SetFont("dejavu", "", 16)
		pdf.SetY(yPos)
		pdf.Cell(40, 10, "Демография")
		yPos += 10

		pdf.SetFont("dejavu", "", 12)
		if len(reportData.Demografi.GenderDistribution) != 0 {
			pdf.SetY(yPos)
			pdf.Cell(40, 10, "Распределение по полу:")
			yPos += 7
			for k, v := range reportData.Demografi.GenderDistribution {
				pdf.SetY(yPos)
				pdf.Cell(40, 10, fmt.Sprintf("%s: %d", k, v))
				yPos += 7
			}
			yPos += 5
		}

		if len(reportData.Demografi.AgeGroups) != 0 {
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
		}

		if len(reportData.Demografi.TopRegions) != 0 {
			pdf.SetY(yPos)
			pdf.Cell(40, 10, "Топ регионов:")
			yPos += 7
			for _, v := range reportData.Demografi.TopRegions {
				pdf.SetY(yPos)
				pdf.Cell(40, 10, fmt.Sprintf("%s: %d", v.Name, v.Count))
				yPos += 7
			}
			yPos += 5
		}
	}

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
		yPos += 5
	}

	if reportData.Finances.AvgOrder != 0 || reportData.Finances.Median != 0 {
		pdf.SetFont("dejavu", "", 16)
		pdf.SetY(yPos)
		pdf.Cell(40, 10, "Финансовые метрики")
		yPos += 10

		pdf.SetFont("dejavu", "", 12)
		if reportData.Finances.AvgOrder != 0 {
			pdf.SetY(yPos)
			pdf.Cell(40, 10, fmt.Sprintf("Средний доход: %f", reportData.Finances.AvgOrder))
			yPos += 7
		}

		if reportData.Finances.Median != 0 {
			pdf.SetY(yPos)
			pdf.Cell(40, 10, fmt.Sprintf("Медиана по среднему чеку: %f", reportData.Finances.Median))
			yPos += 7
		}
		yPos += 5
	}

	// Изображения
	if reportData.Visualizations.GenderPie != "" {
		yPos = checkPageBreak(pdf, yPos, 60)
		if err := addSafeImage(pdf, reportData.Visualizations.GenderPie, marginLeft, yPos); err == nil {
			yPos += 130 // Отступ после диаграммы
		}
	}

	if reportData.Visualizations.AgeHistogram != "" {
		yPos = checkPageBreak(pdf, yPos, 60)
		if err := addSafeImage(pdf, reportData.Visualizations.AgeHistogram, marginLeft, yPos); err == nil {
			yPos += 130 // Отступ после диаграммы
		}
	}

	// Сохранение
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return bytes.Buffer{}, err
	}
	return buf, nil
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

	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

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
