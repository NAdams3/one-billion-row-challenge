package main

import (
	"fmt"
	"math/rand"
	"os"
)

func generateMeasurements(count int) error {

	stdDev := 7.78
	mean := 13.21

	file, err := os.Create(fmt.Sprintf(inFileName, count))
	if err != nil {
		return err
	}
	defer file.Close()

	for range count {

		cityIndex := rand.Intn(len(cities))
		temperature := rand.NormFloat64()*stdDev + mean

		row := fmt.Sprintf("%v;%v\n", cities[cityIndex], fmt.Sprintf("%.1f", temperature))

		_, err = file.Write([]byte(row))
		if err != nil {
			return err
		}

	}

	return nil

}
