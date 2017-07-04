package main

import (
	"os"
	"bufio"
	"flag"
	"fmt"
)

// Returns a list of lines in a file
// It is used to get the Twitter keys for the account
// and to get list of filter words for tweet selection
func processKeyFile(keyFile string) []string {
	// Open file
	file, err := os.Open(keyFile)
	if err != nil {
		panic(err)
	}
	// Defer occurs only after the function ends
	// which makes sense, considering it closes the file
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)

	// Gets each line from the file, using the scanner, and appends it to the array
	for scanner.Scan(){
		lines = append(lines, scanner.Text())
	}
	return lines
}

// This is a generic function to write the contents of an array to a file
// The function will by default overwrite an existing file of the same name
func writeTextFile(fileName string, contents []string) bool {
	// Open the file
	file, err := os.Create(fileName)

	if err != nil {
		fmt.Println("ERROR: Could not create file")
		return false
	}

	// This will close the line once the function ends
	defer file.Close()

	for _,line := range contents{
		file.WriteString(line + "\n")
	}

	// Flushes writes to stable storage
	file.Sync()

	return true
}

// Compares two slices and returns an array with the differences
// Aims at finding unique elements in the "left" slice
func compareSlices(left []string, right []string ) []string {
	var result []string

	for _, leftElement := range left {
		match := false
		for _, rightElement := range right {
			if leftElement == rightElement {
				match = true
				break
			}
		}
		if match == false {
			result = append(result, leftElement)
		}
	}

	return result
}

// Parses arguments and returns a map
func get_commandline_args() map[string]string {
	// Declare command line parameters and their default values

	// The master account has control over what the bot does
	masterNamePtr := flag.String("master", "luisdconejo", "the name of the master account")
	// myname is the screen name of the servant account
	mynamePtr := flag.String("servant", "eran_marno", "screen name of the servant account")

	flag.Parse()

	// Create an empty map (similar to a Python dictionary)
	cmdLineArgs := map[string]string{}
	cmdLineArgs["masterName"] = *masterNamePtr
	cmdLineArgs["servantName"] = *mynamePtr
	return cmdLineArgs
}


