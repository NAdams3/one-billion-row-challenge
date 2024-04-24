package main

import (
	"fmt"
	"os"
	"runtime/pprof"
	"runtime/trace"
	"sort"
	"strconv"
	"strings"
)

func main() {

	profFile, err := os.Create("performance.prof")
	if err != nil {
		panic(err)
	}
	defer profFile.Close()

	err = pprof.StartCPUProfile(profFile)
	if err != nil {
		panic(err)
	}
	defer pprof.StopCPUProfile()

	traceFile, err := os.Create("trace.out")
	if err != nil {
		panic(err)
	}
	defer traceFile.Close()

	err = trace.Start(traceFile)
	if err != nil {
		panic(err)
	}
	defer trace.Stop()

	err = printMinMeanMax()
	if err != nil {
		panic(err)
	}

}

type LocationStats struct {
	min, max, sum float64
	count         int
}

func printMinMeanMax() error {

	out := make(map[string]*LocationStats, 0)
	keys := make([]string, 0)

	fileBytes, err := os.ReadFile("measurements.txt")
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

	outFile, err := os.Create("out.txt")
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
