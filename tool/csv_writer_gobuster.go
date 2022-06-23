package tool

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

func remove(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}

func changePrefix(line string) string {
	rUp, _ := regexp.Compile(`+\[[0-9]+m([A-Z]+)++\[0m`)
	rLow, _ := regexp.Compile(`+\[[0-9]+m([a-z]+)++\[0m`)
	word := strings.Fields(line)
	for index := range word {
		if index >= len(word) {
			break
		}
		if rUp.MatchString(word[index]) {
			tmp := rUp.FindStringSubmatch(word[index])
			word = remove(word, index)
			word = append(word[:index+1], word[index:]...)
			word[index] = tmp[1]
		}
		if rLow.MatchString(word[index]) {
			tmp := rLow.FindStringSubmatch(word[index])
			word = remove(word, index)
			word = append(word[:index+1], word[index:]...)
			word[index] = tmp[1]
		}
	}
	return strings.Join(word, " ")
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
	for nb, line := range gobusterOut {
		if nb != 0 {
			data.WriteString(";")
		}
		line = changePrefix(line)
		data.WriteString(line)
		data.WriteString("\n")
	}
}
