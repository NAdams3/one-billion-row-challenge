package main

import (
	"fmt"
	"math/rand"
	"os"
)

func generateMeasurements(count int) error {

	fileName := fmt.Sprintf(inFileName, count)

	_, err := os.Open(fileName)
	if err == nil {
		fmt.Printf("Measurements file with %v rows already exists.\n", count)
		return nil
	}

	stdDev := 7.78
	mean := 13.21

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	for range count {

		cityIndex := rand.Intn(len(cities))
		temperature := rand.NormFloat64()*stdDev + mean

		row := fmt.Sprintf("%v;%.1f\n", cities[cityIndex], temperature)

		_, err = file.Write([]byte(row))
		if err != nil {
			return err
		}

	}

	return nil

}
