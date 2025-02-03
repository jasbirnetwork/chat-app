package tests

import (
	"bytes"
	"encoding/csv"
	"strings"
	"testing"
)

func parseCSVQuote(csvData string) (string, error) {
	reader := csv.NewReader(bytes.NewBufferString(csvData))
	records, err := reader.ReadAll()
	if err != nil {
		return "", err
	}
	if len(records) < 2 {
		return "", err
	}

	headers := records[0]
	values := records[1]

	var symbolIdx, closeIdx int = -1, -1
	for i, h := range headers {
		switch strings.ToLower(h) {
		case "symbol":
			symbolIdx = i
		case "close":
			closeIdx = i
		}
	}
	if symbolIdx == -1 || closeIdx == -1 {
		return "", err
	}
	quote := values[closeIdx]
	if quote == "N/D" || quote == "" {
		return "", nil
	}
	stockSymbol := strings.ToUpper(values[symbolIdx])
	return stockSymbol + " quote is $" + quote + " per share", nil
}

func TestParseCSVQuote_Valid(t *testing.T) {
	csvData := `Symbol,Date,Time,Open,High,Low,Close,Volume
aapl.us,2025-01-31,22:00:09,150.00,152.00,149.00,151.25,1000000`
	expected := "AAPL.US quote is $151.25 per share"
	result, err := parseCSVQuote(csvData)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestParseCSVQuote_Invalid(t *testing.T) {
	csvData := `Symbol,Date,Time,Open,High,Low,Close,Volume
aapl.us,2025-01-31,22:00:09,150.00,152.00,149.00,N/D,1000000`
	result, err := parseCSVQuote(csvData)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if result != "" {
		t.Errorf("Expected empty result for invalid stock, got '%s'", result)
	}
}
