package scraper

import (
	"archive/zip"
	"encoding/csv"
	"io"
	"log/slog"
)

func createCsv(w io.Writer, data [][]string) error {

	csvWriter := csv.NewWriter(w)
	defer csvWriter.Flush()

	for _, row := range data {
		if err := csvWriter.Write(row); err != nil {
			return err
		}

	}

	return nil
}

func createMultipleCsv(w io.Writer, dataMap map[string][][]string) {
	zipWriter := zip.NewWriter(w)
	defer func() {
		if err := zipWriter.Close(); err != nil {
			slog.Error("Failed to close zip writer", "error", err)
		}
	}()

	for fileName, data := range dataMap {
		file, err := zipWriter.Create(fileName)
		if err != nil {
			slog.Error(err.Error())
		}
		createCsv(file, data)
	}
}
