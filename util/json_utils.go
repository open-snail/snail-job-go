package util

import (
	"encoding/json"
)

func ToByteArr(v interface{}) []byte {
	jsonData, err := json.Marshal(v)
	if err != nil {
		return nil
	}

	return jsonData
}

func ToObj(jsonData []byte, v any) {
	err := json.Unmarshal(jsonData, &v)
	if err != nil {
		return
	}
}
