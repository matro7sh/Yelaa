package tool

import (
	"fmt"
	"os"
	"strings"
)

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
		gobusterOut[line] = strings.Join(lineArray, ";")
		data.WriteString(gobusterOut[line])
		data.WriteString("\n")
		line++
	}
	data.WriteString("log:;")
	gobusterOut = strings.Split(record[3], "\n")
	for nb, line := range gobusterOut {
		if nb != 0 {
			data.WriteString(";")
		}
		data.WriteString(line)
		data.WriteString("\n")
	}
}
