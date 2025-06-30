package messaging

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/xeipuuv/gojsonschema"
	"go.uber.org/zap"
)

type SchemaValidator struct {
	schemas map[string]*gojsonschema.Schema
	logger  *zap.Logger
}

func NewSchemaValidator(logger *zap.Logger) *SchemaValidator {
	return &SchemaValidator{
		schemas: make(map[string]*gojsonschema.Schema),
		logger:  logger,
	}
}

func (v *SchemaValidator) RegisterSchema(topic string, schemaJSON string) error {
	schemaLoader := gojsonschema.NewStringLoader(schemaJSON)
	schema, err := gojsonschema.NewSchema(schemaLoader)
	if err != nil {
		return fmt.Errorf("erro ao carregar schema para o topico %s: %w", topic, err)
	}

	v.schemas[topic] = schema
	v.logger.Info("Schema registrado com sucesso", zap.String("topic", topic))
	return nil
}

func (v *SchemaValidator) ValidateMessage(topic string, message interface{}) error {
	schema, exists := v.schemas[topic]
	if !exists {
		return fmt.Errorf("schema nao registrado para o topico %s", topic)
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("erro ao serializar mensagem: %w", err)
	}

	documentLoader := gojsonschema.NewStringLoader(string(jsonData))
	result, err := schema.Validate(documentLoader)
	if err != nil {
		return fmt.Errorf("erro ao validar mensagem: %w", err)
	}

	if !result.Valid() {
		var errorMessages string
		for _, desc := range result.Errors() {
			errorMessages += fmt.Sprintf("- %s\n", desc)
		}
		return fmt.Errorf("mensagem invalida para o topico %s:\n%s", topic, errorMessages)
	}

	return nil
}

func (v *SchemaValidator) RegisterDefaultSchemas() {
	// Schema para dados do giroscópio
	gyroscopeSchema := `{
		"type": "object",
		"required": ["x", "y", "z", "timestamp", "device_id"],
		"properties": {
			"id": {"type": "integer"},
			"x": {"type": "number"},
			"y": {"type": "number"},
			"z": {"type": "number"},
			"timestamp": {"type": "string", "format": "date-time"},
			"device_id": {"type": "string", "minLength": 1},
			"created_at": {"type": "string", "format": "date-time"},
			"updated_at": {"type": "string", "format": "date-time"}
		}
	}`
	v.RegisterSchema("telemetry.gyroscope", gyroscopeSchema)

	// Schema para dados de GPS
	gpsSchema := `{
		"type": "object",
		"required": ["latitude", "longitude", "timestamp", "device_id"],
		"properties": {
			"id": {"type": "integer"},
			"latitude": {"type": "number"},
			"longitude": {"type": "number"},
			"timestamp": {"type": "string", "format": "date-time"},
			"device_id": {"type": "string", "minLength": 1},
			"created_at": {"type": "string", "format": "date-time"},
			"updated_at": {"type": "string", "format": "date-time"}
		}
	}`
	v.RegisterSchema("telemetry.gps", gpsSchema)

	// Schema para dados de foto
	photoSchema := `{
		"type": "object",
		"required": ["photo", "timestamp", "device_id"],
		"properties": {
			"id": {"type": "integer"},
			"photo": {"type": "string", "minLength": 1},
			"timestamp": {"type": "string", "format": "date-time"},
			"device_id": {"type": "string", "minLength": 1},
			"recognized": {"type": "boolean"},
			"similarity": {"type": "number"},
			"created_at": {"type": "string", "format": "date-time"},
			"updated_at": {"type": "string", "format": "date-time"}
		}
	}`
	v.RegisterSchema("telemetry.photo", photoSchema)
}

func ValidateStruct(data interface{}) error {
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		
		tag := fieldType.Tag.Get("validate")
		if tag == "required" && isEmptyValue(field) {
			return fmt.Errorf("campo %s é obrigatório", fieldType.Name)
		}
	}
	
	return nil
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Struct:
		if t, ok := v.Interface().(time.Time); ok {
			return t.IsZero()
		}
		return false
	default:
		return v.IsNil()
	}
}
