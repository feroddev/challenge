package messaging

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github/feroddev/challengeV3/internal/core"
)

func (p *Producer) encryptSensitiveData(data interface{}) (interface{}, error) {
	if p.encryptor == nil {
		return data, nil
	}

	switch v := data.(type) {
	case core.Gyroscope:
		return p.encryptGyroscopeData(v)
	case core.GPS:
		return p.encryptGPSData(v)
	case core.Photo:
		return p.encryptPhotoData(v)
	default:
		return data, nil
	}
}

func (p *Producer) encryptGyroscopeData(data core.Gyroscope) (core.Gyroscope, error) {
	encryptedDeviceID, err := p.encryptor.Encrypt(data.DeviceID)
	if err != nil {
		return data, fmt.Errorf("erro ao criptografar DeviceID: %w", err)
	}

	data.DeviceID = encryptedDeviceID
	return data, nil
}

func (p *Producer) encryptGPSData(data core.GPS) (core.GPS, error) {
	encryptedDeviceID, err := p.encryptor.Encrypt(data.DeviceID)
	if err != nil {
		return data, fmt.Errorf("erro ao criptografar DeviceID: %w", err)
	}

	data.DeviceID = encryptedDeviceID
	return data, nil
}

func (p *Producer) encryptPhotoData(data core.Photo) (core.Photo, error) {
	encryptedDeviceID, err := p.encryptor.Encrypt(data.DeviceID)
	if err != nil {
		return data, fmt.Errorf("erro ao criptografar DeviceID: %w", err)
	}

	encryptedPhoto, err := p.encryptor.Encrypt(data.Photo)
	if err != nil {
		return data, fmt.Errorf("erro ao criptografar Photo: %w", err)
	}

	data.DeviceID = encryptedDeviceID
	data.Photo = encryptedPhoto
	return data, nil
}

func (p *Producer) encryptMap(data map[string]interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	
	for key, value := range data {
		if key == "device_id" || key == "photo" {
			strValue, ok := value.(string)
			if ok {
				encrypted, err := p.encryptor.Encrypt(strValue)
				if err != nil {
					return nil, fmt.Errorf("erro ao criptografar campo %s: %w", key, err)
				}
				result[key] = encrypted
			} else {
				result[key] = value
			}
		} else if mapValue, ok := value.(map[string]interface{}); ok {
			encryptedMap, err := p.encryptMap(mapValue)
			if err != nil {
				return nil, err
			}
			result[key] = encryptedMap
		} else if sliceValue, ok := value.([]interface{}); ok {
			encryptedSlice, err := p.encryptSlice(sliceValue)
			if err != nil {
				return nil, err
			}
			result[key] = encryptedSlice
		} else {
			result[key] = value
		}
	}
	
	return result, nil
}

func (p *Producer) encryptSlice(data []interface{}) ([]interface{}, error) {
	result := make([]interface{}, len(data))
	
	for i, value := range data {
		if mapValue, ok := value.(map[string]interface{}); ok {
			encryptedMap, err := p.encryptMap(mapValue)
			if err != nil {
				return nil, err
			}
			result[i] = encryptedMap
		} else if sliceValue, ok := value.([]interface{}); ok {
			encryptedSlice, err := p.encryptSlice(sliceValue)
			if err != nil {
				return nil, err
			}
			result[i] = encryptedSlice
		} else {
			result[i] = value
		}
	}
	
	return result, nil
}

func (p *Producer) encryptGenericData(data interface{}) (interface{}, error) {
	if data == nil {
		return nil, nil
	}
	
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar dados para criptografia: %w", err)
	}
	
	var mapData map[string]interface{}
	if err := json.Unmarshal(jsonData, &mapData); err != nil {
		return nil, fmt.Errorf("erro ao deserializar dados para criptografia: %w", err)
	}
	
	encryptedMap, err := p.encryptMap(mapData)
	if err != nil {
		return nil, err
	}
	
	// Se o tipo original era um struct, tenta converter de volta para o mesmo tipo
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Struct {
		encryptedJSON, err := json.Marshal(encryptedMap)
		if err != nil {
			return encryptedMap, nil
		}
		
		newVal := reflect.New(val.Type()).Interface()
		if err := json.Unmarshal(encryptedJSON, newVal); err != nil {
			return encryptedMap, nil
		}
		
		return reflect.ValueOf(newVal).Elem().Interface(), nil
	}
	
	return encryptedMap, nil
}
