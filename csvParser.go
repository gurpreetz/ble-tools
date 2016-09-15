package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

// csvReadFile reads the specified csv file
func csvReadFile(fileName string) (map[string]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error opening file \n\t", err)
		return nil, err
	}
	defer file.Close()

	lines, err := csv.NewReader(file).ReadAll()
	if err != nil {
		fmt.Println("Error reading file \n\t", err)
		return nil, err
	}

	uuidNames := make(map[string]string, len(lines))

	for _, line := range lines {
		uuidNames[line[1]] = line[0]
	}

	return uuidNames, nil
}
