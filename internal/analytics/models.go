package analytics

import (
	"time"
)

type ItemAnalyticsRequest struct {
	Token      string `json:"token"`
	StartDate  string `json:"StartDate"`
	FinishDate string `json:"FinishDate"`
}

type StockItem struct {
	НоменклатураКод    string  `json:"НоменклатураКод"`
	Номенклатура       string  `json:"Номенклатура"`
	Родитель           string  `json:"Родитель"`
	Период             string  `json:"Период"`
	НачальныйОстаток   float64 `json:"НачальныйОстаток"`
	КонечныйОстаток    float64 `json:"КонечныйОстаток"`
	СтатьяРасходов     string  `json:"СтатьяРасходов,omitempty"`
}

type SalesItem struct {
	Код         string  `json:"Код"`
	Номенклатура string  `json:"Номенклатура"`
	Количество  float64 `json:"Количество"`
	Сумма       float64 `json:"Сумма"`
}

type StockEvent struct {
	Time  time.Time `json:"time"`
	Start float64   `json:"start"`
	End   float64   `json:"end"`
}

type ItemAnalyticsResult struct {
	Name         string  `json:"Name"`
	Code         string  `json:"Code"`
	Group        string  `json:"Group"`
	Sales        float64 `json:"Sales"`
	Loss         float64 `json:"Loss"`
	LossOfProfit float64 `json:"LossOfProfit"`
	OSA          float64 `json:"OSA"`
	ABC          string  `json:"ABC"`
}

type AnalyticsResponse struct {
	Items []ItemAnalyticsResult `json:"items"`
	Total int                   `json:"total"`
}

type Chunk struct {
	Items []StockItem
	Index int
}

type ProcessedChunk struct {
	Events map[string][]StockEvent
	Losses map[string]float64
	Index  int
}
