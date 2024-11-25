package util

import (
	"encoding/json"
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
