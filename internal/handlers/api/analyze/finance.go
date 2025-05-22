package analyze

import (
	"first_static_analiz/internal/model"
	"github.com/gin-gonic/gin"
	"math"
	"sort"
)

func FinanceAnalyze(customers []model.Customer, params []string) gin.H {
	result := gin.H{}

	for _, param := range params {
		switch param {
		case "income_stats":
			result["avg_income"] = AvgIncome(customers)
		case "median_stats":
			result["median"] = CalcMedian(customers)
		}
	}

	return result
}

func CalcMedian(customers []model.Customer) float64 {
	//var median float64
	var arrAvgOrder []float64

	for _, customer := range customers {
		arrAvgOrder = append(arrAvgOrder, customer.AvgOrder)
	}

	sort.Float64s(arrAvgOrder)

	length := len(arrAvgOrder)

	if length == 0 {
		return 0
	}

	if length%2 != 0 {
		return arrAvgOrder[length/2]
	}

	return (arrAvgOrder[(length-1)/2] + arrAvgOrder[length/2]) / 2
}

func AvgIncome(customers []model.Customer) float64 {
	var avgincome float64
	for _, customer := range customers {
		avgincome += customer.AvgOrder
	}

	length := float64(len(customers))

	return RoundFloat(avgincome/length, 2)
}

func RoundFloat(x float64, decimals int) float64 {
	pow := math.Pow(10, float64(decimals))
	return math.Round(x*pow) / pow
}
