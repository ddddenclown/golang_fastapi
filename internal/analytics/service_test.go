package analytics

import (
	"os"
	"testing"
	"time"
)

func TestService_NewService(t *testing.T) {
	service := NewService()
	
	if service.workers <= 0 {
		t.Fatalf("Expected positive number of workers, got %d", service.workers)
	}
}

func TestService_SetWorkers(t *testing.T) {
	service := NewService()
	
	service.SetWorkers(8)
	if service.workers != 8 {
		t.Fatalf("Expected 8 workers, got %d", service.workers)
	}
	
	service.SetWorkers(0)
	if service.workers != 8 {
		t.Fatalf("Expected workers to remain 8, got %d", service.workers)
	}
}

func TestService_parseDateTime(t *testing.T) {
	service := NewService()
	
	testCases := []struct {
		input    string
		expected bool
	}{
		{"02.01.2006 15:04:05", true},
		{"02.01.2006 15:04", true},
		{"02.01.2006", true},
		{"invalid-date", false},
		{"", false},
	}
	
	for _, tc := range testCases {
		result := service.parseDateTime(tc.input)
		if tc.expected && result == nil {
			t.Fatalf("Expected valid date for %s, got nil", tc.input)
		}
		if !tc.expected && result != nil {
			t.Fatalf("Expected nil for %s, got %v", tc.input, result)
		}
	}
}

func TestService_calculateOSA(t *testing.T) {
	service := NewService()
	
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
	
	osa := service.calculateOSA([]StockEvent{}, startDate, endDate)
	if osa != 0.0 {
		t.Fatalf("Expected OSA 0.0 for empty events, got %f", osa)
	}
	
	events := []StockEvent{
		{
			Time:  time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			Start: 10.0,
			End:   8.0,
		},
		{
			Time:  time.Date(2024, 1, 1, 18, 0, 0, 0, time.UTC),
			Start: 8.0,
			End:   6.0,
		},
	}
	
	osa = service.calculateOSA(events, startDate, endDate)
	if osa < 0 || osa > 100 {
		t.Fatalf("Expected OSA between 0 and 100, got %f", osa)
	}
}

func TestService_calculateABCClassification(t *testing.T) {
	service := NewService()
	
	items := []ItemAnalyticsResult{}
	service.calculateABCClassification(items)
	
	items = []ItemAnalyticsResult{
		{Sales: 1000.0},
		{Sales: 500.0},
		{Sales: 200.0},
		{Sales: 100.0},
	}
	
	service.calculateABCClassification(items)
	
	for _, item := range items {
		if item.ABC == "" {
			t.Fatalf("Expected ABC classification for item with sales %f", item.Sales)
		}
		if item.ABC != "A" && item.ABC != "B" && item.ABC != "C" {
			t.Fatalf("Expected ABC classification to be A, B, or C, got %s", item.ABC)
		}
	}
}

func TestService_splitIntoChunks(t *testing.T) {
	service := NewService()
	service.SetWorkers(2)
	
	stockData := []StockItem{
		{НоменклатураКод: "1001"},
		{НоменклатураКод: "1002"},
		{НоменклатураКод: "1003"},
		{НоменклатураКод: "1004"},
	}
	
	chunks := service.splitIntoChunks(stockData)
	
	if len(chunks) != 2 {
		t.Fatalf("Expected 2 chunks, got %d", len(chunks))
	}
	
	totalItems := 0
	for _, chunk := range chunks {
		totalItems += len(chunk.Items)
	}
	
	if totalItems != len(stockData) {
		t.Fatalf("Expected %d total items, got %d", len(stockData), totalItems)
	}
}

func TestService_processChunk(t *testing.T) {
	service := NewService()
	
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
	
	chunk := Chunk{
		Items: []StockItem{
			{
				НоменклатураКод:    "1001",
				Период:             "01.01.2024 12:00:00",
				НачальныйОстаток:   10.0,
				КонечныйОстаток:    8.0,
				СтатьяРасходов:     "Порча на складах (94)",
			},
		},
		Index: 0,
	}
	
	result := service.processChunk(chunk, startDate, endDate)
	
	if len(result.Events) == 0 {
		t.Fatal("Expected events to be processed")
	}
	
	if len(result.Losses) == 0 {
		t.Fatal("Expected losses to be calculated")
	}
}

func TestService_GetItemAnalytics_InvalidDates(t *testing.T) {
	service := NewService()
	
	req := &ItemAnalyticsRequest{
		Token:      "test-token",
		StartDate:  "invalid-date",
		FinishDate: "01.01.2024",
	}
	
	_, err := service.GetItemAnalytics(req)
	if err == nil {
		t.Fatal("Expected error for invalid start date")
	}
	
	req.StartDate = "01.01.2024"
	req.FinishDate = "invalid-date"
	
	_, err = service.GetItemAnalytics(req)
	if err == nil {
		t.Fatal("Expected error for invalid finish date")
	}
}

func createTestDataFiles(t *testing.T) {
	if err := os.MkdirAll("routes", 0755); err != nil {
		t.Fatalf("Failed to create routes directory: %v", err)
	}
	
	stockData := `[
		{
			"НоменклатураКод": "1001",
			"Номенклатура": "Товар 1",
			"Родитель": "Группа 1",
			"Период": "01.01.2024 00:00:00",
			"НачальныйОстаток": 10,
			"КонечныйОстаток": 8
		}
	]`
	
	err := os.WriteFile("routes/stock_dump.json", []byte(stockData), 0644)
	if err != nil {
		t.Fatalf("Failed to create test stock file: %v", err)
	}
	
	salesData := `[
		{
			"Код": "1001",
			"Номенклатура": "Товар 1",
			"Количество": 5,
			"Сумма": 500
		}
	]`
	
	err = os.WriteFile("routes/sales_dump.json", []byte(salesData), 0644)
	if err != nil {
		t.Fatalf("Failed to create test sales file: %v", err)
	}
}

func cleanupTestDataFiles(t *testing.T) {
	os.Remove("routes/stock_dump.json")
	os.Remove("routes/sales_dump.json")
}

func TestService_GetItemAnalytics_Integration(t *testing.T) {
	createTestDataFiles(t)
	t.Cleanup(func() {
		cleanupTestDataFiles(t)
	})
	
	service := NewService()
	
	req := &ItemAnalyticsRequest{
		Token:      "test-token",
		StartDate:  "01.01.2024",
		FinishDate: "31.01.2024",
	}
	
	response, err := service.GetItemAnalytics(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	if response == nil {
		t.Fatal("Expected non-nil response")
	}
	
	if response.Total < 0 {
		t.Fatalf("Expected non-negative total, got %d", response.Total)
	}
}
