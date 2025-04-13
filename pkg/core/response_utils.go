package core

import (
	"fmt"
	"strconv"
	"strings"
)

// convertHeaders transforma o mapa de headers do YAML para o formato de resposta
// Suporta:
// - Valores string diretos
// - Números convertidos para string
// - Booleanos convertidos para string
// - Slices unidos com vírgula (para multi-value)
func convertHeaders(headersConfig interface{}) map[string]string {
	headers := make(map[string]string)

	if headersMap, ok := headersConfig.(map[string]interface{}); ok {
		for key, value := range headersMap {
			switch v := value.(type) {
			case string:
				headers[key] = v
			case bool:
				headers[key] = strconv.FormatBool(v)
			case int, int32, int64, float32, float64:
				headers[key] = fmt.Sprintf("%v", v)
			case []interface{}: // Para valores múltiplos (ex: Set-Cookie)
				var strValues []string
				for _, item := range v {
					strValues = append(strValues, fmt.Sprintf("%v", item))
				}
				headers[key] = strings.Join(strValues, ", ")
			default:
				headers[key] = fmt.Sprintf("%v", v)
			}
		}
	}

	// Garante Content-Type padrão se não especificado
	if _, exists := headers["Content-Type"]; !exists {
		headers["Content-Type"] = "application/json"
	}

	return headers
}
