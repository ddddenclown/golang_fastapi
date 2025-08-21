package analytics

import (
	"testing"
	"time"
)

func BenchmarkService_calculateOSA(b *testing.B) {
	service := NewService()
	
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	

	events := make([]StockEvent, 1000)
	for i := 0; i < 1000; i++ {
		events[i] = StockEvent{
			Time:  startDate.Add(time.Duration(i) * time.Hour),
			Start: float64(i % 100),
			End:   float64((i + 1) % 100),
		}
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.calculateOSA(events, startDate, endDate)
	}
}

func BenchmarkService_calculateABCClassification(b *testing.B) {
	service := NewService()
	

	items := make([]ItemAnalyticsResult, 10000)
	for i := 0; i < 10000; i++ {
		items[i] = ItemAnalyticsResult{
			Sales: float64(i + 1),
		}
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		itemsCopy := make([]ItemAnalyticsResult, len(items))
		copy(itemsCopy, items)
		service.calculateABCClassification(itemsCopy)
	}
}

func BenchmarkService_parseDateTime(b *testing.B) {
	service := NewService()
	dateStr := "02.01.2024 15:04:05"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.parseDateTime(dateStr)
	}
}

func BenchmarkService_splitIntoChunks(b *testing.B) {
	service := NewService()
	service.SetWorkers(8)
	

	stockData := make([]StockItem, 100000)
	for i := 0; i < 100000; i++ {
		stockData[i] = StockItem{
			НоменклатураКод: string(rune(i % 1000)),
		}
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.splitIntoChunks(stockData)
	}
}

func BenchmarkService_processChunk(b *testing.B) {
	service := NewService()
	
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	

	chunk := Chunk{
		Items: make([]StockItem, 1000),
		Index: 0,
	}
	
	for i := 0; i < 1000; i++ {
		chunk.Items[i] = StockItem{
			НоменклатураКод:    string(rune(i % 100)),
			Период:             "01.01.2024 12:00:00",
			НачальныйОстаток:   float64(i % 100),
			КонечныйОстаток:    float64((i + 1) % 100),
			СтатьяРасходов:     "Порча на складах (94)",
		}
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.processChunk(chunk, startDate, endDate)
	}
}

func BenchmarkService_processDataParallel(b *testing.B) {
	service := NewService()
	service.SetWorkers(8)
	
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	

	stockData := make([]StockItem, 10000)
	salesData := make([]SalesItem, 1000)
	
	for i := 0; i < 10000; i++ {
		stockData[i] = StockItem{
			НоменклатураКод:    string(rune(i % 1000)),
			Период:             "01.01.2024 12:00:00",
			НачальныйОстаток:   float64(i % 100),
			КонечныйОстаток:    float64((i + 1) % 100),
		}
	}
	
	for i := 0; i < 1000; i++ {
		salesData[i] = SalesItem{
			Код:         string(rune(i % 1000)),
			Количество:  float64(i % 10),
			Сумма:       float64(i * 100),
		}
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.processDataParallel(stockData, salesData, startDate, endDate)
	}
}


func BenchmarkService_processDataParallel_2Workers(b *testing.B) {
	service := NewService()
	service.SetWorkers(2)
	benchmarkProcessDataParallel(service, b)
}

func BenchmarkService_processDataParallel_4Workers(b *testing.B) {
	service := NewService()
	service.SetWorkers(4)
	benchmarkProcessDataParallel(service, b)
}

func BenchmarkService_processDataParallel_8Workers(b *testing.B) {
	service := NewService()
	service.SetWorkers(8)
	benchmarkProcessDataParallel(service, b)
}

func BenchmarkService_processDataParallel_16Workers(b *testing.B) {
	service := NewService()
	service.SetWorkers(16)
	benchmarkProcessDataParallel(service, b)
}

func benchmarkProcessDataParallel(service *Service, b *testing.B) {
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	

	stockData := make([]StockItem, 5000)
	salesData := make([]SalesItem, 500)
	
	for i := 0; i < 5000; i++ {
		stockData[i] = StockItem{
			НоменклатураКод:    string(rune(i % 1000)),
			Период:             "01.01.2024 12:00:00",
			НачальныйОстаток:   float64(i % 100),
			КонечныйОстаток:    float64((i + 1) % 100),
		}
	}
	
	for i := 0; i < 500; i++ {
		salesData[i] = SalesItem{
			Код:         string(rune(i % 1000)),
			Количество:  float64(i % 10),
			Сумма:       float64(i * 100),
		}
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.processDataParallel(stockData, salesData, startDate, endDate)
	}
}
