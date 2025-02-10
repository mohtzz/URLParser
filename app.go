package main

/*вывод времени выполнения программы в конце, после каждой операции
прога общается с юзером */

import (
	"fmt"
	"io"
	"os"
	// "net/http"
	// "io"
	// "time"
)

func main() {

	// filePath := flag.String("src", "", "путь к файлу")
	// directoryPath := flag.String("dst", "", "путь к директории")
	// flag.Parse()
	file, err := os.Open("sites.txt")
	if err != nil {
		fmt.Println("не удалось открыть файл")
		os.Exit(1)
	}
	defer file.Close()

	// эта часть с чтением файла вообще не понятна
	data := make([]byte, 64)

	for {
		n, err := file.Read(data)
		if err == io.EOF {
			break
		}
		fmt.Println(string(data[:n]))
	}
}
