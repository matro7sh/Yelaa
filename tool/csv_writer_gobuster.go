package tool

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
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
	regex, _ := regexp.Compile(`\[(.*?)\]|\((.*?)\)`)
	for index := range lineArray {
		if regex.MatchString(lineArray[index]) {
			all_sub_string := regex.FindStringSubmatch(lineArray[index])
			if all_sub_string[1] == "" {
				lineArray[index] = all_sub_string[2]
			} else {
				lineArray[index] = all_sub_string[1]
			}
		}
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
		lineArray = moveSizeAndStatusToSameArrayCase(lineArray)
		lineArray = reformatLine(lineArray)
		data.WriteString(strings.Join(lineArray, ";") + ";\n")
	}
}
