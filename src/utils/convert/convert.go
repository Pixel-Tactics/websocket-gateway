package convert_utils

import (
	"encoding/json"
	"errors"
)

var ErrInvalidData = errors.New("invalid data")

func MapToObject(data map[string]interface{}, obj any) error {
	encoded, err := json.Marshal(data)
	if err != nil {
		return ErrInvalidData
	}

	err = json.Unmarshal(encoded, obj)
	if err != nil {
		return ErrInvalidData
	}

	return nil
}

func ObjectToMap(obj any) (map[string]interface{}, error) {
	encoded, err := json.Marshal(obj)
	if err != nil {
		return nil, ErrInvalidData
	}

	var data map[string]interface{}
	err = json.Unmarshal(encoded, &data)
	if err != nil {
		return nil, ErrInvalidData
	}

	return data, nil
}
