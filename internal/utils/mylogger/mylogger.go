package mylogger

import (
	"log"
	"runtime"
	"strconv"
)

func PrintGoroutines(text string) {
	log.Println(text, "Goroutines: " + strconv.Itoa(runtime.NumGoroutine()))
}
