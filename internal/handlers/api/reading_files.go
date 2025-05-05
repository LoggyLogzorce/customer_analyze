package api

import (
	"encoding/csv"
	"first_static_analiz/internal/model"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Чтение файла (CSV/Excel)
func readCustomerFile(path string) ([]model.Customer, error) {
	ext := filepath.Ext(path)

	switch ext {
	case ".csv":
		return readCSV(path)
	//case ".xlsx", ".xls":
	//	return readExcel(path)
	default:
		return nil, fmt.Errorf("unsupported file format: %s", ext)
	}
}

// Пример обработки CSV
func readCSV(path string) ([]model.Customer, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	r := csv.NewReader(file)
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	var customers []model.Customer
	for i, record := range records {
		if i == 0 { // Пропускаем заголовок
			continue
		}

		if len(record) != 7 {
			return nil, fmt.Errorf("invalid record length in line %d", i+1)
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

	return customers, nil
}
