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
    if s == nil {
		return false
	}
	for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}

// Check if a slice of int contains a value
func ContainsString(s []string, e string) bool {
	if s == nil {
		return false
	}
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}
// Remove all occurences of the element in the slice
func RemoveElement(s []string, element string) []string {
	list := make([]string,0)
	for _,e := range s {
		if e != element {
			list = append(list, e)
		}
	}
	return list;
}
