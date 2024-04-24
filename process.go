package main

import (
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

	fileBytes, err := os.ReadFile(fmt.Sprintf(inFileName, count))
	if err != nil {
		return err
	}

	rows := strings.Split(string(fileBytes), "\n")

	for _, row := range rows {
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
