package util

import (
	"log"
)

// Protect a function call from panic
func Protect(g func()) {
	defer func() {
		// executes normally even if there is a panic
		if err := recover(); err != nil {
			log.Println("run time panic: %v", err)
		}
	}()
	g() // possible runtime-error
}

// Check if a slice of int contains a value
func Contains(s []int, e int) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}
