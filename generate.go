package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"sync"
)

func writerWorker(file *os.File, chunk string, wg *sync.WaitGroup, mutex *sync.Mutex) error {

	mutex.Lock()

	_, err := file.Write([]byte(chunk))
	if err != nil {
		return err
	}

	mutex.Unlock()

	wg.Done()

	return nil
}

func generateMeasurements(count int, force bool) error {

	chunkSize := int(math.Floor(math.Pow(float64(count), .333)))
	/* numChunks := int(math.Ceil(float64(count) / float64(chunkSize))) */

	stdDev := 7.78
	mean := 13.21

	fileName := fmt.Sprintf(inFileName, count)

	_, err := os.Open(fileName)
	if err == nil && !force {
		fmt.Printf("Measurements file with %v rows already exists.\n", count)
		return nil
	}

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	var wg sync.WaitGroup
	var mutex sync.Mutex

	rows := ""
	rowsCount := 0

	for range count {

		cityIndex := rand.Intn(len(cities))
		temperature := rand.NormFloat64()*stdDev + mean
		rows += fmt.Sprintf("%v;%.1f\n", cities[cityIndex], temperature)
		rowsCount++

		if rowsCount == chunkSize {
			wg.Add(1)
			go writerWorker(file, rows, &wg, &mutex)
			rows = ""
			rowsCount = 0
		}

	}

	wg.Add(1)
	go writerWorker(file, rows, &wg, &mutex)

	wg.Wait()

	return nil

}
