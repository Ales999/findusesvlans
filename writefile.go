package main

import "os"

// WriteOutFile - создать и записать в файл
func WriteOutFile(lines []string, fileName string) {

	f, err := os.Create(fileName) // Create or trunk file
	if err != nil {
		panic(err)
	}
	defer f.Close()

	for _, line := range lines {
		_, err := f.WriteString(line)
		if err != nil {
			panic(err)
		}
	}
}
