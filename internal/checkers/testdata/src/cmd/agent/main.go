package main

import "os"

func main() {
	os.Exit(1) // want "os.Exit should not be used in main.go"
}
