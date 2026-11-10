package library

import (
	"os"
	
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func Add(path string, key string, value string) error {
	currentLibrary, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	newLibrary, _ := sjson.SetBytes(currentLibrary, key, value)
	err = os.WriteFile(path, newLibrary, 0o755)
	return err
}

func Read(path string) (map[string]gjson.Result, error) {
	currentLibrary, err := os.ReadFile(path)
	if err != nil {
		return map[string]gjson.Result{}, err
	}
	libraryMap := gjson.GetBytes(currentLibrary, "@this").Map()
	return libraryMap, nil
}
