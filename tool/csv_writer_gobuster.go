package tool

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/CMEPW/Yelaa/helper"
)

var regex = regexp.MustCompile(`\[(.*?)\]|\((.*?)\)`)

func remove(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}

func reformatLine(lineArray []string) []string {
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

func CsvWriterGobuster(g *GoBuster) {
	data, err := os.Create(helper.YelaaPath + "/gobuster/scan_data-" +
		time.Now().Format("2006-01-02_15-04-05") + ".csv")
	if err != nil {
		fmt.Print(err)
		return
	}
	defer data.Close()
	data.WriteString("Url:;" + g.optDir.URL + "\nMethod:;" + g.optDir.Method +
		"\nThreads:;" + fmt.Sprintf("%d\n", g.opts.Threads) + "Wordlist:;" +
		g.opts.Wordlist + "\nOutput File name:;" + g.opts.OutputFilename +
		"\nTimeout:;" + g.optDir.Timeout.String() + "\n")
	if err != nil {
		fmt.Print(err)
		return
	}
	log, err := os.ReadFile(g.opts.OutputFilename)
	if err != nil {
		fmt.Print(err)
		return
	}
	out := strings.Split(string(log[:]), "\n")
	for line, str := range out {
		if !strings.Contains(str, "Status: 200") {
			continue
		}
		lineArray := strings.Fields(out[line])
		lineArray = moveSizeAndStatusToSameArrayCase(lineArray)
		lineArray = reformatLine(lineArray)
		data.WriteString(strings.Join(lineArray, ";") + "\n")
	}
}
