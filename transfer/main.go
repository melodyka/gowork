package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
)

// 读取CSV文件并获取指定列和行的数据
func readCSV(filename string) ([][]string, []string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	header := strings.Split(rows[0][0], " ")
	header[0] = strings.TrimSpace(header[0])

	data := make([][]string, len(rows)-1)
	for i := 1; i < len(rows); i++ {
		row := strings.Split(rows[i][0], " ")
		for j := range row {
			row[j] = strings.TrimSpace(row[j])
		}
		data[i-1] = row
	}

	return data, header
}

func main() {
	filename := "grid.csv" // CSV文件路径
	//path := "/path/to/csv" // CSV文件所在路径
	data, header := readCSV(filename)

	hh := []string{"N"}
	for _, h := range header[0][1:] {
		hh = append(hh, string(h))
	}

	l := make([][]string, len(data))
	for i, row := range data {
		l[i] = make([]string, len(row))
		for j, element := range row {
			l[i][j] = strings.TrimSpace(element)
		}
	}

	// 示例用法
	index := []int{1, 2, 3, 4, 5, 6, 7}
	df := make([][]string, len(index))
	for i, idx := range index {
		df[i] = l[idx-1][1:]
	}

	// 打印 DataFrame
	for _, row := range df {
		fmt.Println(row)
	}

	m := "A"
	n := 3
	ri := n - 1
	ci := int([]rune(m)[0]) - int('A')

	fmt.Println(df[ri][ci])
}
