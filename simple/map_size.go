package main

import (
	"fmt"
	"runtime"
)

func main() {
	v := struct{}{}

	m := make(map[int]struct{})
	for i := 0; i < 10000; i++ {
		m[i] = v
	}

	// Manual GC
	runtime.GC()
	printMemStats("Deleted 10K keys result")

	// set is nil recycle
	m = nil
	runtime.GC()
	printMemStats("Set is nil result")

}

func printMemStats(msg string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("%v: Assign Mem = %vKB, GC count = %v\n", msg, m.Alloc/1024, m.NumGC)
}
