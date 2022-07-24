package tool

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func remove(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}

func lineSkipper(data *os.File) *os.File {
	r := bufio.NewReader(data)
	for {
		_, err := r.ReadBytes('\n')
		if err == nil {
			break
		}
	}
	return data
}

func reformatLine(lineArray []string) []string {
	for index := range lineArray {
		lineArray[index] = strings.Replace(lineArray[index], "(", "", -1)
		lineArray[index] = strings.Replace(lineArray[index], ")", "", -1)
		lineArray[index] = strings.Replace(lineArray[index], "[", "", -1)
		lineArray[index] = strings.Replace(lineArray[index], "]", "", -1)
	}
	return lineArray
}

func moveSizeAndStatusToSameArrayCase(lineArray []string) []string {
	index := 0
	for index < len(lineArray) {
		if strings.Contains(lineArray[index], "Status") || strings.Contains(lineArray[index], "Size") {
			lineArray[index] = strings.Join(lineArray[index:index+2], " ")
			lineArray = remove(lineArray, index+1)
		}
		index++
	}
	return lineArray
}

func CsvWriterGobuster() {
	data, err := os.OpenFile("scan_data.csv", os.O_RDWR, 0666)
	if err != nil {
		fmt.Print(err)
		return
	}
	data = lineSkipper(data)
	log, err := os.ReadFile("scan_log_gobuster.txt")
	if err != nil {
		fmt.Print(err)
		return
	}
	out := strings.Split(string(log[:]), "\n")
	for line, str := range out {
		if !strings.Contains(str, "200") {
			continue
		}
		lineArray := strings.Fields(out[line])
		lineArray = reformatLine(lineArray)
		lineArray = moveSizeAndStatusToSameArrayCase(lineArray)
		data.WriteString(strings.Join(lineArray, ";"))
		data.WriteString(";\n")
	}
}
