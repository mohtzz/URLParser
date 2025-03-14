package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func main() {
	start := time.Now()
	defer func() {
		progTime := time.Since(start)
		fmt.Println("Время выполнения программы: ", progTime)
	}()

	inputFilePath := flag.String("src", "", "input")
	resultDir := flag.String("dst", "result", "output")
	flag.Parse()
	fmt.Println("Флаги считаны")

	// проверка на корректность флагов

	if *inputFilePath == "" || *resultDir == "" {
		fmt.Println("Не указаны флаги/один из флагов(src/dst)")
		return
	}

	// чтение файла с помощью функции
	sites, err := readSitesFromFile(*inputFilePath)
	if err != nil {
		fmt.Println("Ошибка чтения файла: ", err)
		return
	}

	// создание директории при ее отсутствии

	if _, err := os.Stat(*resultDir); os.IsNotExist(err) {
		fmt.Println("Директория не была найдена, создание директории...")
		// 0777 - полный доступ для всех
		os.MkdirAll(*resultDir, 0777)
	}

	// get-запросы и обработка всех сайтов из файла
	wg := sync.WaitGroup{}
	for _, site := range sites {
		wg.Add(1)
		go func(site string) {
			defer wg.Done()
			err := processSite(site, *resultDir)
			if err != nil {
				fmt.Println("Ошибка обработки адреса: ", err)
			}
		}(site)
	}
	wg.Wait()
}

// readSitesFromFile - функция, считывающая адреса с файла и записывающая в массив для дальнейшей работы
func readSitesFromFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	fmt.Println("Файл открыт, чтение...")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	//в созданный срез записываем append'ом адреса с файла
	var sites []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}
		sites = append(sites, scanner.Text())
	}
	return sites, nil
}

// processSite - функция для обработки адреса(get-запрос и превращение в .html-файл)
func processSite(site string, resultDir string) error {
	resp, err := http.Get(site)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// проверка статус-кода и форматирование
	if resp.StatusCode == http.StatusOK {
		filename := filepath.Join(resultDir, strings.ReplaceAll(site, "/", "_")+".html")
		return saveHTML(filename, resp.Body)
	} else {
		fmt.Println("Сайт", site, "вернул статус-код: ", resp.StatusCode)
		return nil
	}
}

// saveHTML - функция, создающая файл и записывающая в него тело ответа get-запроса
func saveHTML(filename string, body io.Reader) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, body)
	if err != nil {
		return err
	}
	fmt.Println("Файл", filename, "сохранен.")
	return nil
}
