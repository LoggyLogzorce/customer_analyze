package model

import "time"

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

	Finances struct {
		AvgOrder float64 `json:"avg_order"`
		Median   float64 `json:"median"`
	} `json:"finances"`

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
