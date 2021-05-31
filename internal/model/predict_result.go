package model

type PredictResult struct {
	ImageUrl string   `json:"imageUrl"`
	Classes  []string `json:"classes"`
	Score    float64  `json:"score"`
}
