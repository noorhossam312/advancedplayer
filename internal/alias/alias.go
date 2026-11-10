package alias

import (
	"os"
	
	"github.com/tidwall/sjson"
	"github.com/tidwall/gjson"
)

func Add(path string, key string, value string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	newContent, _ := sjson.SetBytes(file, key, value)
	return os.WriteFile(path, newContent, 0o755)
}

func Read(path string) (map[string]gjson.Result, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return map[string]gjson.Result{}, err
	}
	content := gjson.GetBytes(file, "@this")
	return content.Map(), nil
}

