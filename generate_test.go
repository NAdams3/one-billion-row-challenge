package main

import (
	"fmt"
	"testing"
)

func TestGenerate(t *testing.T) {

	err := generate(20)
	if err != nil {
		t.Errorf("expected nil, got err: %v \n", err)
	}

}

func BenchmarkGenerate(b *testing.B) {

	count := 10000
	fmt.Printf("Generate %v lines \n", count)

	generate(count)

}
