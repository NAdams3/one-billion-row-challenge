package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type LocationStats struct {
	min, max, sum float64
	count         int
}

func processMeasurements(count int) error {

	out := make(map[string]*LocationStats, 0)
	keys := make([]string, 0)

	// fileBytes, err := os.ReadFile(fmt.Sprintf(inFileName, count))
	// if err != nil {
	// 	return err
	// }
	//
	// rows := strings.Split(string(fileBytes), "\n")
	//

	file, err := os.Open(fmt.Sprintf(inFileName, count))
	if err != nil {
		return err
	}

	chunks := make(chan string, 1)
	searchString := "\n"
	searchBytes := []byte(searchString)
	searchLen := len(searchBytes)

	go populateChannel(file, chunks, searchBytes, searchLen)

	for chunk := range chunks {
		rows := strings.Split(chunk, "\n")

		for _, row := range rows {
			/* 	fmt.Println(row) */
			columns := strings.Split(row, ";")
			if len(columns) < 2 {
				continue
			}

			location := columns[0]
			temperature, err := strconv.ParseFloat(columns[1], 64)
			if err != nil {
				return err
			}

			stats := out[location]

			if stats == nil {
				out[location] = &LocationStats{
					min:   temperature,
					max:   temperature,
					sum:   temperature,
					count: 1,
				}
				keys = append(keys, location)
				continue
			}

			if temperature < stats.min {
				stats.min = temperature
			}

			if temperature > stats.max {
				stats.max = temperature
			}

			stats.sum += temperature
			stats.count++

		}
	}

	sort.Strings(keys)

	outFile, err := os.Create(fmt.Sprintf(outFileName, count))
	if err != nil {
		return err
	}
	defer outFile.Close()

	for _, key := range keys {

		outStats := out[key]

		outRow := fmt.Sprintf("%v;%.1f;%.1f;%.1f\n", key, outStats.min, outStats.sum/float64(outStats.count), outStats.max)
		outFile.Write([]byte(outRow))
	}

	return nil

}

func populateChannel(file *os.File, ch chan string, searchBytes []byte, searchLen int) {
	scanner := bufio.NewScanner(file)
	scanner.Split(splitAtLast(searchBytes, searchLen))
	for scanner.Scan() {
		ch <- scanner.Text()
	}
	close(ch)
}

func splitAtLast(searchBytes []byte, searchLen int) func(data []byte, atEOF bool) (advance int, token []byte, err error) {
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		dataLen := len(data)

		if atEOF && dataLen == 0 {
			return 0, nil, nil
		}

		if i := bytes.LastIndex(data, searchBytes); i >= 0 {
			return i + searchLen, data[0:i], nil
		}

		if atEOF {
			return dataLen, data, nil
		}

		return 0, nil, nil

	}
}
