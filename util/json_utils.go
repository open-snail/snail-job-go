package util

import (
	"encoding/json"
	"path/filepath"
)

func ToByteArr(v interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

func ToObj(jsonData []byte, v any) error {
	err := json.Unmarshal(jsonData, &v)
	if err != nil {
		return err
	}

	return nil
}

func TrimProjectPath(fullPath, projectRoot string) string {
	relativePath, err := filepath.Rel(projectRoot, fullPath)
	if err != nil {
		// 如果出错，直接返回原路径
		return fullPath
	}
	return relativePath
}
