package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"
)

// ToJsonString converte uma interface{} para uma string JSON.
func ToJsonString(data interface{}) string {
	if data == nil {
		return ""
	}
	// Se já for string, retorna direto (útil para corpos não-JSON)
	if str, ok := data.(string); ok {
		return str
	}
	b, err := json.Marshal(data)
	if err != nil {
		log.Printf("ERROR: Failed to marshal data to JSON: %v", err)
		// Retorna uma representação de erro JSON em vez de panic
		return `{"error": "internal serialization error"}`
	}
	return string(b)
}

func ConvertMapToJSON(input map[interface{}]interface{}) (string, error) {
	// Primeiro converter para map[string]interface{}
	stringMap := make(map[string]interface{})

	for key, value := range input {
		strKey := fmt.Sprintf("%v", key) // Converte qualquer tipo de chave para string
		stringMap[strKey] = value
	}

	// Serializar para JSON
	jsonBytes, err := json.Marshal(stringMap)
	if err != nil {
		return "", fmt.Errorf("erro ao serializar JSON: %w", err)
	}

	return string(jsonBytes), nil
}

func ConvertAnyMapToJSON(input interface{}) (string, error) {
	// Se for uma string, verifique se é JSON válido
	if jsonStr, ok := input.(string); ok {
		var obj interface{}
		// Tenta deserializar para verificar se é JSON válido
		if err := json.Unmarshal([]byte(jsonStr), &obj); err == nil {
			// Se for JSON válido, retorne a string original sem serializar novamente
			return jsonStr, nil
		}
	}

	// Continua o processamento normal para outros tipos
	if m, ok := ConvertToMapStringInterface(input); ok {
		input = m
	}

	jsonBytes, err := json.Marshal(input)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return string(jsonBytes), nil
}

func normalizeValue(value interface{}) (interface{}, error) {
	switch v := value.(type) {
	case map[interface{}]interface{}:
		result := make(map[string]interface{})
		for key, val := range v {
			strKey := fmt.Sprintf("%v", key)
			normalizedVal, err := normalizeValue(val)
			if err != nil {
				return nil, err
			}
			result[strKey] = normalizedVal
		}
		return result, nil
	case []interface{}:
		var result []interface{}
		for _, item := range v {
			normalizedItem, err := normalizeValue(item)
			if err != nil {
				return nil, err
			}
			result = append(result, normalizedItem)
		}
		return result, nil
	default:
		return v, nil
	}
}

// SetNestedValue define um valor em um mapa aninhado usando um path estilo dot-notation.
// Ex: SetNestedValue(myMap, "user.address.city", "New York")
func SetNestedValue(data map[string]interface{}, path string, value interface{}) error {
	keys := strings.Split(path, ".")
	lastKey := keys[len(keys)-1]
	current := data

	for i := 0; i < len(keys)-1; i++ {
		key := keys[i]
		next, ok := current[key]
		if !ok {
			// Cria mapa intermediário se não existir
			newMap := make(map[string]interface{})
			current[key] = newMap
			current = newMap
		} else {
			// Verifica se o valor existente é um mapa
			if nextMap, ok := next.(map[string]interface{}); ok {
				current = nextMap
			} else {
				return fmt.Errorf("path element '%s' in '%s' is not a map (type: %T)", key, path, next)
			}
		}
	}

	// Define o valor final
	current[lastKey] = value
	return nil
}

// GetNestedValue obtém um valor de um mapa/struct aninhado usando dot-notation.
func GetNestedValue(data interface{}, path string) (interface{}, bool) {
	keys := strings.Split(path, ".")
	current := data

	for _, key := range keys {
		if current == nil {
			return nil, false
		}

		currentValue := reflect.ValueOf(current)
		// Dereferencia ponteiros
		if currentValue.Kind() == reflect.Ptr {
			currentValue = currentValue.Elem()
		}

		switch currentValue.Kind() {
		case reflect.Map:
			// Assume map[string]interface{} ou similar
			mapKeys := currentValue.MapKeys()
			found := false
			for _, mapKey := range mapKeys {
				if mapKey.Kind() == reflect.String && mapKey.String() == key {
					current = currentValue.MapIndex(mapKey).Interface()
					found = true
					break
				}
			}
			if !found {
				return nil, false // Chave não encontrada no mapa
			}
		case reflect.Struct:
			fieldValue := currentValue.FieldByName(key)
			if !fieldValue.IsValid() {
				// Tenta buscar por tag JSON se o campo não for encontrado diretamente
				fieldValue = findFieldByJSONTag(currentValue, key)
				if !fieldValue.IsValid() {
					return nil, false // Campo não encontrado na struct
				}
			}
			// Verifica se o campo é exportado antes de chamar Interface()
			if !fieldValue.CanInterface() {
				log.Printf("WARN: Field '%s' in path '%s' is unexported in struct %s", key, path, currentValue.Type().Name())
				return nil, false
			}
			current = fieldValue.Interface()
		default:
			// Não é possível navegar mais fundo
			return nil, false
		}
	}
	return current, true
}

// findFieldByJSONTag busca um campo em uma struct pela sua tag JSON.
func findFieldByJSONTag(structValue reflect.Value, jsonTag string) reflect.Value {
	structType := structValue.Type()
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		tag := field.Tag.Get("json")
		// Considera a tag exata ou a parte antes da vírgula (ex: "name,omitempty")
		tagParts := strings.Split(tag, ",")
		if len(tagParts) > 0 && tagParts[0] == jsonTag {
			return structValue.Field(i)
		}
	}
	return reflect.Value{} // Retorna valor inválido se não encontrado
}

// ToStringMap converte interface{} (mapa ou struct) para map[string]string.
func ToStringMap(data interface{}) map[string]string {
	m := make(map[string]string)
	if data == nil {
		return m
	}
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Map {
		iter := v.MapRange()
		for iter.Next() {
			key := fmt.Sprintf("%v", iter.Key().Interface())
			val := fmt.Sprintf("%v", iter.Value().Interface())
			m[key] = val
		}
	}
	// Adicionar conversão de struct se necessário
	return m
}

// GetEnvVarsMap e GetSecretsMap (podem ficar aqui ou em config/utils se mais específicos)
func GetEnvVarsMap() map[string]string {
	// Implementar leitura de variáveis de ambiente relevantes
	log.Println("Placeholder: Reading relevant ENV variables...")
	// Exemplo: pegar todas as vars com um prefixo específico?
	return map[string]string{"EXAMPLE_API_KEY": "dummy_key_from_util"}
}

func GetSecretsMap() map[string]string {
	// Implementar leitura de secrets (ex: AWS Secrets Manager)
	log.Println("Placeholder: Reading relevant Secrets...")
	return map[string]string{"DB_PASSWORD": "dummy_secret_from_util"}
}
