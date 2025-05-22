package api

import (
	"encoding/csv"
	"first_static_analiz/internal/model"
	"fmt"
	"github.com/extrame/xls"
	"github.com/xuri/excelize/v2"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Чтение файла (CSV/Excel)
func readCustomerFile(path string) ([]model.Customer, int, error) {
	ext := filepath.Ext(path)

	switch ext {
	case ".csv":
		return readCSV(path)
	case ".xlsx":
		return readXLSX(path)
	case ".xls":
		return readXLS(path)
	default:
		return nil, 400, fmt.Errorf("unsupported file format: %s", ext)
	}
}

func readCSV(path string) ([]model.Customer, int, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, 500, fmt.Errorf("ошибка открытия файла %v", err)
	}
	defer file.Close()

	r := csv.NewReader(file)
	records, err := r.ReadAll()
	if err != nil {
		return nil, 500, fmt.Errorf("ошибка чтения %v", err)
	}

	return parseFileRows(records)
}

func readXLSX(path string) ([]model.Customer, int, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, 500, fmt.Errorf("ошибка открытия файла %v", err)
	}
	defer f.Close()

	// Получаем первый лист
	sheet := f.GetSheetName(0)
	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, 500, fmt.Errorf("ошибка считывания первого листа %v", err)
	}

	return parseFileRows(rows)
}

func readXLS(path string) ([]model.Customer, int, error) {
	xlFile, err := xls.Open(path, "utf-8")
	if err != nil {
		return nil, 500, fmt.Errorf("ошибка открытия файла %v", err)
	}

	sheet := xlFile.GetSheet(0)
	if sheet == nil {
		return nil, 400, fmt.Errorf("empty sheet")
	}

	var rows [][]string
	for i := 0; i <= int(sheet.MaxRow); i++ {
		row := sheet.Row(i)
		var cells []string
		for j := 0; j < 7; j++ { // Ожидаем 7 колонок
			cells = append(cells, row.Col(j))
		}
		rows = append(rows, cells)
	}

	return parseFileRows(rows)
}

func parseFileRows(rows [][]string) ([]model.Customer, int, error) {
	var customers []model.Customer
	for i, record := range rows {
		if i == 0 { // Пропускаем заголовок
			continue
		}

		if len(record) != 7 {
			return nil, 400, fmt.Errorf("ошибка в строении файла, строка №%d", i+1)
		}

		customer := model.Customer{}

		// Парсинг ID
		customer.ID, _ = strconv.Atoi(record[0])

		// Парсинг возраста
		customer.Age, _ = strconv.Atoi(record[1])

		// Парсинг пола (приводим к нижнему регистру)
		customer.Gender = strings.ToLower(strings.TrimSpace(record[2]))

		// Регион
		customer.Region = strings.TrimSpace(record[3])

		// Дата регистрации
		customer.RegDate, _ = time.Parse("2006-01-02", record[4])

		// Количество заказов
		customer.OrdersCount, _ = strconv.Atoi(record[5])

		// Средний чек
		customer.AvgOrder, _ = strconv.ParseFloat(record[6], 64)

		customers = append(customers, customer)
	}

	return customers, 200, nil
}
