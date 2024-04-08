package pkg

import (
	"YadroProject/types"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
)

type Storage interface {
	Save(string, interface{}) error
	Read(string) error
}

type Data struct {
	JSONPath string
}

func NewData(s string) *Data {
	return &Data{
		JSONPath: s,
	}
}

func (d *Data) Save(filename string, comics interface{}) error {

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Создаем новый JSON encoder, который пишет в файл
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ") // Устанавливаем отступы для красивого форматирования

	// Записываем данные в файл
	if err := encoder.Encode(comics); err != nil {
		return err
	}

	// Если необходимо добавить разделители между объектами JSON, можно записать разделитель в файл
	if _, err := file.Write([]byte("\n")); err != nil {
		return err
	}

	return nil
}

func (d *Data) Read(filepath string) error {
	jsonData, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return err
	}

	// Десериализуем данные из JSON
	var comics map[string]types.Num
	if err := json.Unmarshal(jsonData, &comics); err != nil {
		fmt.Printf("Error unmarshalling JSON: %v\n", err)
		return err
	}

	// Получаем и сортируем ключи для гарантии порядка обхода
	var keys []string
	for k := range comics {
		keys = append(keys, k)
	}
	sort.Strings(keys) // Это важно, если ключи не являются строго упорядоченными числами

	// Выводим первые n объектов
	for i, k := range keys {
		if i >= d.cfg.Xkcd.DbSize {
			break
		}
		fmt.Printf("=============\n")
		fmt.Printf("ID: %s\n URL: %s\n Keywords: %v\n", k, comics[k].URL, comics[k].Keywords)
		fmt.Printf("=============\n")
	}
	if err != nil {
		log.Fatalf("Error reading database: %v", err)
	}
	return nil
}
