package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"log/slog"
	"os"

	"golang.org/x/text/encoding/charmap"
)

func prepareCSV() *bytes.Reader {
	file, err := os.ReadFile("TimeTable_04_10_2026.csv")
	if err != nil {
		slog.Error("failed to read file", "error", err)
		return nil
	}
	decoder := charmap.Windows1251.NewDecoder()
	decoded, err := decoder.Bytes(file)
	if err != nil {
		slog.Error("failed to decode file", "error", err)
		return nil
	}
	replaced := bytes.ReplaceAll(decoded, []byte("\r"), []byte("\n"))
	fixedCSV := bytes.NewReader(replaced)
	return fixedCSV
}

func readCSVFile(fixedCSV *bytes.Reader) [][]string {
	r := csv.NewReader(fixedCSV)
	records, err := r.ReadAll()
	if err != nil {
		slog.Error("failed to parse CSV", "error", err)
		return nil
	}
	return records
}

func main() {
	data := prepareCSV()
	if data == nil {
		slog.Error("Import error")
		return
	}
	records := readCSVFile(data)
	if records == nil {
		slog.Error("Import error")
		return
	}
	records = records[1:]
	fmt.Println(records)
}
