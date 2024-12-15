package service

import (
	repository "a21hc3NpZ25tZW50/repository/fileRepository"
	"encoding/csv"
	"errors"
	"strings"
)

type FileService struct {
	Repo *repository.FileRepository
}

func (s *FileService) ProcessFile(fileContent string) (map[string][]string, error) {
	// Parse the CSV content from the file
	reader := csv.NewReader(strings.NewReader(fileContent))
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, errors.New("failed to parse CSV file")
	}

	if len(rows) == 0 {
		return nil, errors.New("empty CSV file")
	}

	// Create a map to store the parsed data
	data := make(map[string][]string)

	// Assume the first row contains headers
	headers := rows[0]
	for i := 1; i < len(rows); i++ {
		for j, value := range rows[i] {
			header := headers[j]
			data[header] = append(data[header], value)
		}
	}

	// Return the parsed data
	return data, nil
}
