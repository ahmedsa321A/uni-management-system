package Registration

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"university-management/backend/models"
)
func RegisterStudent(path string ) error {
	file, err := os.OpenFile(path)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}
	for i, row := range records {
		if i == 0 {
			continue
		}

		age, err := strconv.Atoi(row[1])
		if err != nil {
			fmt.Println("Error converting age:", err)
			continue
		}

		studen := models.Student{
			Name:  row[0],
			Age:   age,
			Email: row[2],
		}
	}

}
