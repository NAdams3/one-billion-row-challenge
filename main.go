package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/pprof"
	"runtime/trace"
)

const inFileName = "measurements-%v.txt"
const outFileName = "out-%v.txt"
const profFileName = "%v-%v-profile.prof"
const traceFileName = "%v-%v-trace.out"

func main() {

	generate := flag.Bool("generate", false, "Generates a measurements file with count lines instead of processing")
	count := flag.Int("count", 0, "The length of which measurements file to process.")
	flag.Parse()

	if *count > 1_000_000_000 {
		fmt.Println("Bro 1 billion isn't enough for you?")
		return
	}

	action := "process"
	if *generate {
		action = "generate"
	}

	profFile, err := os.Create(fmt.Sprintf(profFileName, action, *count))
	if err != nil {
		panic(err)
	}
	defer profFile.Close()

	err = pprof.StartCPUProfile(profFile)
	if err != nil {
		panic(err)
	}
	defer pprof.StopCPUProfile()

	traceFile, err := os.Create(fmt.Sprintf(traceFileName, action, *count))
	if err != nil {
		panic(err)
	}
	defer traceFile.Close()

	err = trace.Start(traceFile)
	if err != nil {
		panic(err)
	}
	defer trace.Stop()

	if *generate {
		fmt.Printf("Generating measurements file with %v rows.\n", *count)
		err = generateMeasurements(*count)
		if err != nil {
			panic(err)
		}
		fmt.Println("Generation successful!")
	} else {
		fmt.Printf("Processing measurements file with %v rows.\n", *count)
		err = processMeasurements(*count)
		if err != nil {
			panic(err)
		}
		fmt.Println("Measurements have been processed successfully!")
	}

}
