package internal

import (
	"log"
	"time"
)

func TimeFunction(name string) func() {
	start := time.Now()
	return func() {
		log.Printf("%s took: %v\n", name, time.Since(start))
	}
}
