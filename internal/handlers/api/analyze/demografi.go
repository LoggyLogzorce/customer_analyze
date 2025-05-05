package analyze

import (
	"first_static_analiz/internal/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"sort"
)

type RegionStat struct {
	Name  string
	Count int
}

func DemografiAnalyze(customers []model.Customer, params []string) gin.H {
	var result, allAgeGroup, ageGroup, regionStat = gin.H{}, gin.H{}, gin.H{}, gin.H{}

	for _, param := range params {
		switch param {
		case "gender_dist":
			result["gender_dist"] = genderDistribution(customers)
		case "18-25":
			allAgeGroup = ageDistribution(customers)
			ageGroup["18-25"] = allAgeGroup["18-25"]
			ageGroup["Count"] = allAgeGroup["Count"]
		case "26-35":
			if allAgeGroup["26-35"] == nil {
				allAgeGroup = ageDistribution(customers)
			}
			ageGroup["26-35"] = allAgeGroup["26-35"]
			ageGroup["Count"] = allAgeGroup["Count"]
		case "36-45":
			if allAgeGroup["36-45"] == nil {
				allAgeGroup = ageDistribution(customers)
			}
			ageGroup["36-45"] = allAgeGroup["36-45"]
			ageGroup["Count"] = allAgeGroup["Count"]
		case "46+":
			if allAgeGroup["46+"] == nil {
				allAgeGroup = ageDistribution(customers)
			}
			ageGroup["46+"] = allAgeGroup["46+"]
			ageGroup["Count"] = allAgeGroup["Count"]
		case "age_histogram":
			if allAgeGroup["18-25"] == nil {
				allAgeGroup = ageDistribution(customers)
			}
			ageGroup = allAgeGroup
		case "gender_pie":
			result["gender_pie"] = genderDistribution(customers)
		case "top5":
			regionStat["top5"], _ = GetTopRegions(customers, 5)
		case "top10":
			regionStat["top10"], _ = GetTopRegions(customers, 10)
		}
	}

	if ageGroup["18-25"] != nil || ageGroup["26-35"] != nil || ageGroup["36-45"] != nil || ageGroup["46+"] != nil {
		result["age_group"] = ageGroup
	}

	if regionStat["top5"] != nil || regionStat["top10"] != nil {
		result["region_stat"] = regionStat
	}

	return result
}

// Распределение по полу
func genderDistribution(customers []model.Customer) gin.H {
	counts := make(gin.H)

	for _, c := range customers {
		if val, exists := counts[c.Gender]; exists {
			counts[c.Gender] = val.(int) + 1
		} else {
			counts[c.Gender] = 1
		}
	}

	return counts
}

// Распределение по возрастам
func ageDistribution(customers []model.Customer) gin.H {
	ageGroups := gin.H{
		"18-25": 0,
		"26-35": 0,
		"36-45": 0,
		"46+":   0,
		"Count": 0,
	}

	for _, c := range customers {
		switch {
		case c.Age <= 25:
			ageGroups["18-25"] = ageGroups["18-25"].(int) + 1
			ageGroups["Count"] = ageGroups["Count"].(int) + 1
		case c.Age <= 35:
			ageGroups["26-35"] = ageGroups["26-35"].(int) + 1
			ageGroups["Count"] = ageGroups["Count"].(int) + 1
		case c.Age <= 45:
			ageGroups["36-45"] = ageGroups["36-45"].(int) + 1
			ageGroups["Count"] = ageGroups["Count"].(int) + 1
		default:
			ageGroups["46+"] = ageGroups["46+"].(int) + 1
			ageGroups["Count"] = ageGroups["Count"].(int) + 1
		}
	}

	return ageGroups
}

// GetTopRegions TOP регионов
func GetTopRegions(customers []model.Customer, topNumber int) ([]RegionStat, error) {
	if len(customers) == 0 {
		return nil, fmt.Errorf("empty customers list")
	}

	// Считаем количество клиентов по регионам
	regionCounts := make(map[string]int)
	for _, customer := range customers {
		regionCounts[customer.Region]++
	}

	// Конвертируем в слайс для сортировки
	var regions []RegionStat
	for region, count := range regionCounts {
		regions = append(regions, RegionStat{
			Name:  region,
			Count: count,
		})
	}

	// Сортируем регионы:
	// 1. По убыванию количества клиентов
	// 2. По алфавиту если количество одинаковое
	sort.Slice(regions, func(i, j int) bool {
		if regions[i].Count == regions[j].Count {
			return regions[i].Name < regions[j].Name
		}
		return regions[i].Count > regions[j].Count
	})

	// Возвращаем топ N результатов
	if len(regions) > topNumber {
		return regions[:topNumber], nil
	}
	return regions, nil
}
