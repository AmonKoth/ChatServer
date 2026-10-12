package main

import "fmt"

func WriteupError(err error) {
	if err != nil {
		fmt.Println("Error Occurred: ", err)
	}
}
