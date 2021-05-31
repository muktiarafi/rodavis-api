package entity

import "time"

type Report struct {
	ID           int       `json:"id"`
	UserID       int       `json:"-"`
	ReporterName string    `json:"reporterName"`
	Status       string    `json:"status"`
	ImageURL     string    `json:"imageUrl"`
	Classes      []string  `json:"classes"`
	Note         string    `json:"note"`
	Address      string    `json:"address"`
	Location     *Location `json:"location"`
	DateReported time.Time `json:"dateReported"`
}

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}
