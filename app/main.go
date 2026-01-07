package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	// открываем файл .txt
	queries, err := os.Open("links.txt")
	if err != nil {
		fmt.Println("[ERROR]: ошибка открытия файла")
		return
	}
	// закрываем тело запроса
	defer queries.Close()
	// вызываем функцию readLines
	lines, err := readLines(queries)
	if err != nil {
		fmt.Println("[ERROR]: ошибка вызова функции - readLines")
		return
	}
	//вызываем функцию search
	err = search(lines)
	if err != nil {
		fmt.Println("[ERROR]: ошибка вызова функции - search")
		return
	}

}

func readLines(file *os.File) ([]string, error) {
	// переменная куда мы будем складывать строки
	var lines []string
	// создаем новый объект скана
	scanner := bufio.NewScanner(file)

	// тело сканнера, которое получает данные с .txt
	for scanner.Scan() {
		// чтение каждой строки полученной из .Scan()
		line := scanner.Text()
		// удаление пробелов, табуляции и переноса строки
		line = strings.TrimSpace(line)
		// проверка - пустая ли строка
		// если строка не пустая, добавляем ее в slice lines
		if line != "" {
			lines = append(lines, line)
		}
	}
	// проверяем ошибки сканера после цикла
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	// возвращаем итоговое значение без ошибок
	return lines, nil
}

func search(lines []string) error {
	// открываем файл для записи итоговых путей
	result, err := os.Create("result.txt")
	if err != nil {
		fmt.Println("[ERROR]: ошибка создания файла result.txt")
		return err
	}
	// закрываем тело после создания
	defer result.Close()

	writer := bufio.NewWriter(result)

	root := `C:\`

	for _, query := range lines {
		err := walkDir(root, query, writer)
		if err != nil {
			fmt.Println("[ERROR]: ошибка поиска для запроса", query, err)
		}
	}
	writer.Flush()

	return nil
}

func walkDir(dir string, query string, writer *bufio.Writer) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println("[ERROR]: ошибка чтения пути")
		return err
	}

	for _, entry := range entries {
		fullPath := dir + "\\" + entry.Name()

		if entry.Name() == query {
			if _, err := writer.WriteString(query + " | " + fullPath + "\n"); err != nil {
				return err
			}
		} else if strings.HasPrefix(entry.Name(), query) {
			if _, err := writer.WriteString(query + " | " + fullPath + "\n"); err != nil {
				return err
			}
		}

		if entry.IsDir() {
			if err := walkDir(fullPath, query, writer); err != nil {
				fmt.Println("[ERROR]: ошибка при обходе папок", fullPath, err)
			}
		}
	}

	return nil
}
