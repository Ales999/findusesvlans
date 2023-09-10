package main

import "os"

// WriteOutFile - создать и записать в файл
func WriteOutFile(lines []string, fileName string, forceoverwrite bool) {

	var fileptr *os.File
	var err error

	if forceoverwrite { //
		fileptr, err = os.Create(fileName) // Create or trunk file
	} else { // Если флаг не задан то будем добавлять
		fileptr, err = os.OpenFile(fileName, os.O_CREATE|os.O_APPEND, 0644)
	}
	if err != nil {
		panic(err)
	}
	defer fileptr.Close()

	// Записываем файлы
	for _, line := range lines {
		_, err := fileptr.WriteString(line)
		if err != nil {
			panic(err)
		}
	}
}
