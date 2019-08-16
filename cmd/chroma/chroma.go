package main

import (
	"fmt"
	"os"

	"github.com/phR0ze/chroma/internal/chroma"
)

func main() {
	r := chroma.New()
	if err := r.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
