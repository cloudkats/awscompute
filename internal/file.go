package internal

import (
	"encoding/csv"
	"os"
	"path/filepath"
)

type FileWriter struct {
}

func (fw FileWriter) WriteToFile(path string, input [][]string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(filepath.Dir(path), 0700)
		if err != nil {
			return err
		}
	}
	data := [][]string{
		{
			"AccountID", "AccountAlias", "Type", "CPU", "Memory",
		},
	}
	for _, i := range input {
		data = append(data, i)
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	for _, value := range data {
		err = writer.Write(value)
		if err != nil {
			return err
		}
	}
	return nil
}
