package tool

import (
	"fmt"
	"os"
	"strings"
)

func remove(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}

func CsvWriterGobuster(record []string) {
	data, err := os.Create("scan_data.csv")
	if err != nil {
		fmt.Print(err)
	}
	defer data.Close()
	gobusterOut := strings.Split(record[2], "\n")
	line := 4
	for line < len(gobusterOut)-5 {
		lineArray := strings.Fields(gobusterOut[line])
		lineArray = remove(lineArray, 0)
		for !strings.Contains(lineArray[0], ":") {
			lineArray = remove(lineArray, 0)
		}
		gobusterOut[line] = strings.Join(lineArray, ";")
		data.WriteString(gobusterOut[line])
		data.WriteString("\n")
		line++
	}
	data.WriteString("log:;")
	gobusterOut = strings.Split(record[3], "\n")
	for _, line := range gobusterOut {
		data.WriteString(line)
		data.WriteString("\n")
	}
}
