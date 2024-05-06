package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"sort"
)

type LocationStats struct {
	min, max, sum float64
	count         int
}

func processMeasurements(count int) error {

	out := make(map[string]*LocationStats, 1000)
	keys := make([]string, 0, 1000)

	file, err := os.Open(fmt.Sprintf(inFileName, count))
	if err != nil {
		return err
	}

	searchString := "\n"
	searchBytes := []byte(searchString)
	searchLen := len(searchBytes)

	columnSplit := ";"
	columnSplitBytes := []byte(columnSplit)
	columnSplitLen := len(columnSplitBytes)

	scanner := bufio.NewScanner(file)
	scanner.Split(splitAtLast(searchBytes, searchLen))

	bufMax := 1024 * 1024
	buf := make([]byte, bufMax)
	scanner.Buffer(buf, bufMax)

	for scanner.Scan() {
		chunk := scanner.Bytes()

		for len(chunk) > 0 {

			endIndex := bytes.Index(chunk, searchBytes)
			row := chunk
			if endIndex != -1 {
				row = chunk[:endIndex]
				chunk = chunk[endIndex+searchLen:]
			} else {
				chunk = []byte("")
			}

			splitIndex := bytes.Index(row, columnSplitBytes)
			if splitIndex == -1 || len(row[splitIndex+columnSplitLen:]) <= 0 {
				continue
			}

			location := string(row[:splitIndex])
			fmt.Printf("temp: %v \n", string(row[splitIndex+columnSplitLen:]))
			temperature := float64FromBytes(row[splitIndex+columnSplitLen:])

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

	err = writeOut(keys, out, count)
	if err != nil {
		return err
	}

	return nil

}

func writeOut(keys []string, out map[string]*LocationStats, count int) error {

	sort.Strings(keys)

	outFile, err := os.Create(fmt.Sprintf(outFileName, count))
	if err != nil {
		return err
	}
	defer outFile.Close()

	for _, key := range keys {

		outStats := out[key]

		outRow := fmt.Sprintf("%v;%.1f;%.1f;%.1f\n",
			key,
			outStats.min,
			outStats.sum/float64(outStats.count),
			outStats.max,
		)
		outFile.Write([]byte(outRow))
	}

	return nil

}

func float64FromBytes(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	return math.Float64frombits(bits)
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
