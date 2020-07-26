package csv

import (
	"io"
	"encoding/csv"
	"os"
)

type CsvDataSet struct {
	Columns []string;
	Records [][]string;
}

func (c *CsvDataSet) ForEachRow(processRow func(row map[string]string)) {
	for _, record := range c.Records {
		recordMap := make(map[string]string)
		for columnIndex, column := range c.Columns {
			recordMap[column] = record[columnIndex]
		}
		processRow(recordMap)
	}
}

func ReadDataSet(filePath string) (*CsvDataSet, error) {
	file, fileOpenErr := os.Open(filePath)
	if fileOpenErr != nil {
		return nil, fileOpenErr
	}
	fileCsvReader := csv.NewReader(file)

	var columns []string
	var records [][]string = make([][]string, 0)

	var nextRow []string
	var nextRowErr error
	if nextRow, nextRowErr = fileCsvReader.Read(); nextRowErr == nil {
		columns = nextRow
		nextRow, nextRowErr = fileCsvReader.Read()
		for nextRowErr == nil {
			records = append(records, nextRow)
			nextRow, nextRowErr = fileCsvReader.Read()
		}
	}

	if nextRowErr != nil && nextRowErr != io.EOF {
		return nil, nextRowErr
	} else {
		return &CsvDataSet {
			columns,
			records,
		}, nil
	}
}
