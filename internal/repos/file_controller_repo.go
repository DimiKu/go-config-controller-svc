package repos

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type FileRepo struct {
	localPath string
}

func NewFileRepo(localPath string) *FileRepo {
	return &FileRepo{localPath: localPath}
}

func (f *FileRepo) GetValuesFromFile(filePath string) (map[string]map[string]interface{}, error) {
	configMap := make(map[string]map[string]interface{})
	path := fmt.Sprintf("%s/%s/%s.yaml", f.localPath, filePath, filePath)
	file, err := os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&configMap)
	if err != nil {
		log.Fatalf("Failed to decode yaml YAML: %v", err)
	}

	return configMap, nil
}
