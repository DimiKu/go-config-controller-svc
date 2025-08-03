package repos

import (
	"fmt"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"os"
)

type FileRepo struct {
	localPath string
	log       *zap.Logger
}

func NewFileRepo(localPath string, log *zap.Logger) *FileRepo {
	return &FileRepo{localPath: localPath, log: log}
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
		f.log.Error("Failed to decode yaml YAML: %v", zap.Error(err))
	}

	return configMap, nil
}
