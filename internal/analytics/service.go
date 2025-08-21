package analytics

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"
)

type Service struct {
	workers int
}

func NewService() *Service {
	workers := runtime.NumCPU()
	if workers < 2 {
		workers = 2
	}
	return &Service{
		workers: workers,
	}
}

func (s *Service) SetWorkers(workers int) {
	if workers > 0 {
		s.workers = workers
	}
}

func (s *Service) GetItemAnalytics(req *ItemAnalyticsRequest) (*AnalyticsResponse, error) {
	startTime := time.Now()
	
	startDate, err := time.Parse("02.01.2006", req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date format: %w", err)
	}
	
	finishDate, err := time.Parse("02.01.2006", req.FinishDate)
	if err != nil {
		return nil, fmt.Errorf("invalid finish date format: %w", err)
	}
	
	finishDate = finishDate.Add(24 * time.Hour)
	
	stockData, err := s.loadStockData()
	if err != nil {
		return nil, fmt.Errorf("failed to load stock data: %w", err)
	}
	
	salesData, err := s.loadSalesData()
	if err != nil {
		return nil, fmt.Errorf("failed to load sales data: %w", err)
	}
	
	log.Printf("Loaded %d stock items and %d sales items", len(stockData), len(salesData))
	log.Printf("Date range: %s to %s", startDate.Format("02.01.2006"), finishDate.Format("02.01.2006"))
	
	items, err := s.processDataParallel(stockData, salesData, startDate, finishDate)
	if err != nil {
		return nil, fmt.Errorf("failed to process data: %w", err)
	}
	
	s.calculateABCClassification(items)
	
	sort.Slice(items, func(i, j int) bool {
		return items[i].Sales > items[j].Sales
	})
	
	processingTime := time.Since(startTime)
	log.Printf("Analytics processing completed in %v", processingTime)
	log.Printf("Generated %d analytics items", len(items))
	
	return &AnalyticsResponse{
		Items: items,
		Total: len(items),
	}, nil
}

func (s *Service) loadStockData() ([]StockItem, error) {
	file, err := os.Open("routes/stock_dump.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	
	var items []StockItem
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, err
	}
	
	return items, nil
}

func (s *Service) loadSalesData() ([]SalesItem, error) {
	file, err := os.Open("routes/sales_dump.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	
	var items []SalesItem
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, err
	}
	
	return items, nil
}

func (s *Service) processDataParallel(stockData []StockItem, salesData []SalesItem, startDate, finishDate time.Time) ([]ItemAnalyticsResult, error) {
	chunks := s.splitIntoChunks(stockData)
	
	results := make(chan ProcessedChunk, len(chunks))
	chunkChan := make(chan Chunk, len(chunks))
	var wg sync.WaitGroup
	
	for i := 0; i < s.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for chunk := range chunkChan {
				processed := s.processChunk(chunk, startDate, finishDate)
				results <- processed
			}
		}()
	}
	
	go func() {
		for _, chunk := range chunks {
			chunkChan <- chunk
		}
		close(chunkChan)
	}()
	
	go func() {
		wg.Wait()
		close(results)
	}()
	
	allEvents := make(map[string][]StockEvent)
	allLosses := make(map[string]float64)
	
	for result := range results {
		for code, events := range result.Events {
			allEvents[code] = append(allEvents[code], events...)
		}
		for code, loss := range result.Losses {
			allLosses[code] += loss
		}
	}
	
	salesByCode := make(map[string]float64)
	salesQtyByCode := make(map[string]float64)
	nameByCode := make(map[string]string)
	groupByCode := make(map[string]string)
	
	for _, item := range salesData {
		code := strings.TrimSpace(item.Код)
		if code == "" {
			continue
		}
		
		salesByCode[code] += item.Сумма
		salesQtyByCode[code] += item.Количество
		nameByCode[code] = item.Номенклатура
	}
	
	for _, item := range stockData {
		code := strings.TrimSpace(item.НоменклатураКод)
		if code != "" {
			groupByCode[code] = item.Родитель
		}
	}
	
	priceByCode := make(map[string]float64)
	for code, qty := range salesQtyByCode {
		if qty > 0 {
			priceByCode[code] = salesByCode[code] / qty
		}
	}
	
	var items []ItemAnalyticsResult
	for code, totalSales := range salesByCode {
		price := priceByCode[code]
		lossQty := allLosses[code]
		lossAmount := lossQty * price
		
		osa := s.calculateOSA(allEvents[code], startDate, finishDate)
		
		name := nameByCode[code]
		if name == "" {
			name = code
		}
		
		group := groupByCode[code]
		if group == "" {
			group = "Без группы 🤔"
		}
		
		lossPercent := 0.0
		if totalSales > 0 {
			lossPercent = (lossAmount / totalSales) * 100
		}
		
		items = append(items, ItemAnalyticsResult{
			Name:         name,
			Code:         code,
			Group:        group,
			Sales:        math.Round(totalSales*100) / 100,
			Loss:         math.Round(lossAmount*100) / 100,
			LossOfProfit: math.Round(lossPercent*1000) / 1000,
			OSA:          osa,
		})
	}
	
	return items, nil
}

func (s *Service) splitIntoChunks(data []StockItem) []Chunk {
	chunkSize := (len(data) + s.workers - 1) / s.workers
	var chunks []Chunk
	
	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize
		if end > len(data) {
			end = len(data)
		}
		
		chunks = append(chunks, Chunk{
			Items: data[i:end],
			Index: i / chunkSize,
		})
	}
	
	return chunks
}

func (s *Service) processChunk(chunk Chunk, startDate, finishDate time.Time) ProcessedChunk {
	events := make(map[string][]StockEvent)
	losses := make(map[string]float64)
	
	for _, item := range chunk.Items {
		code := strings.TrimSpace(item.НоменклатураКод)
		if code == "" {
			continue
		}
		
		dt := s.parseDateTime(item.Период)
		if dt == nil {
			continue
		}
		
		events[code] = append(events[code], StockEvent{
			Time:  *dt,
			Start: item.НачальныйОстаток,
			End:   item.КонечныйОстаток,
		})
		
		if item.СтатьяРасходов == "Порча на складах (94)" {
			diff := item.НачальныйОстаток - item.КонечныйОстаток
			if diff > 0 {
				losses[code] += diff
			}
		}
	}
	
	return ProcessedChunk{
		Events: events,
		Losses: losses,
		Index:  chunk.Index,
	}
}

func (s *Service) parseDateTime(dateStr string) *time.Time {
	formats := []string{
		"02.01.2006 15:04:05",
		"02.01.2006 15:04",
		"02.01.2006",
	}
	
	for _, format := range formats {
		if dt, err := time.Parse(format, dateStr); err == nil {
			return &dt
		}
	}
	
	return nil
}

func (s *Service) calculateOSA(events []StockEvent, startDate, finishDate time.Time) float64 {
	if len(events) == 0 {
		return 0.0
	}
	
	sort.Slice(events, func(i, j int) bool {
		return events[i].Time.Before(events[j].Time)
	})
	
	balance := events[0].Start
	current := startDate
	availHours := 0.0
	
	for _, event := range events {
		t := event.Time
		if t.Before(startDate) {
			balance = event.End
			continue
		}
		if t.After(finishDate) {
			break
		}
		
		if balance > 0 {
			availHours += t.Sub(current).Hours()
		}
		
		balance = event.End
		current = t
	}
	
	if current.Before(finishDate) && balance > 0 {
		availHours += finishDate.Sub(current).Hours()
	}
	
	totalHours := finishDate.Sub(startDate).Hours()
	if totalHours <= 0 {
		return 0.0
	}
	
	return math.Round((availHours/totalHours)*10000) / 100
}

func (s *Service) calculateABCClassification(items []ItemAnalyticsResult) {
	if len(items) == 0 {
		return
	}
	
	totalSales := 0.0
	for _, item := range items {
		totalSales += item.Sales
	}
	
	if totalSales <= 0 {
		return
	}
	
	cumulative := 0.0
	for i := range items {
		cumulative += items[i].Sales
		share := (cumulative / totalSales) * 100
		
		if share <= 80 {
			items[i].ABC = "A"
		} else if share <= 95 {
			items[i].ABC = "B"
		} else {
			items[i].ABC = "C"
		}
	}
}
