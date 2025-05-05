package model

import "time"

type Customer struct {
	ID          int       `csv:"id" json:"id"`
	Age         int       `csv:"возраст" json:"age"`
	Gender      string    `csv:"пол" json:"gender"`
	Region      string    `csv:"регион" json:"region"`
	RegDate     time.Time `csv:"дата регистрации" json:"reg_date"`
	OrdersCount int       `csv:"количество заказов" json:"orders_count"`
	AvgOrder    float64   `csv:"средний чек" json:"avg_order"`
}
