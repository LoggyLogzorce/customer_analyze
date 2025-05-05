package analyze

import (
	"first_static_analiz/internal/model"
	"github.com/gin-gonic/gin"
	"sort"
	"time"
)

// CustomersAnalyze Комбинированная функция анализа
func CustomersAnalyze(customers []model.Customer, params []string) gin.H {
	result := gin.H{}

	for _, param := range params {
		switch param {
		case "veterans":
			result["veterans"] = analyzeVeterans(customers)["veterans"]
		case "newbies":
			result["newbies"] = analyzeNewbies(customers)["newbies"]
		case "vips":
			result["vips"] = analyzeVIPs(customers)
		}
	}

	return result
}

// Анализ "Ветераны" (регистрация до 2023)
func analyzeVeterans(customers []model.Customer) gin.H {
	veterans := gin.H{
		"veterans": 0,
	}
	cutoffDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	for _, c := range customers {
		if c.RegDate.Before(cutoffDate) {
			//veterans = append(veterans, c)
			veterans["veterans"] = veterans["veterans"].(int) + 1
		}
	}
	return veterans
}

// Анализ "Новички" (регистрация в 2024)
func analyzeNewbies(customers []model.Customer) gin.H {
	newbies := gin.H{
		"newbies": 0,
	}
	cutoffDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	for _, c := range customers {
		if c.RegDate.After(cutoffDate) {
			//newbies = append(newbies, c)
			newbies["newbies"] = newbies["newbies"].(int) + 1
		}
	}
	return newbies
}

// Анализ VIP-клиенты (чек выше 75-го перцентиля)
func analyzeVIPs(customers []model.Customer) gin.H {
	checks := make([]float64, len(customers))
	for i, c := range customers {
		checks[i] = c.AvgOrder
	}

	sort.Float64s(checks)
	n := len(checks)
	percentile75 := checks[(n*75)/100]

	//var vips []model.Customer
	vips := gin.H{
		"Count":        0,
		"percentile75": percentile75,
	}

	for _, c := range customers {
		if c.AvgOrder > percentile75 {
			//vips = append(vips, c)
			vips["Count"] = vips["Count"].(int) + 1
		}
	}
	return vips
}
