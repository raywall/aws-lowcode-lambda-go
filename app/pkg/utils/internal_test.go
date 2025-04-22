// internal_test.go
package utils

import (
	"reflect"
	"testing"
)

// --- Testes para convertToMap ---

func TestConvertToMap(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name:    "Input já é map[string]interface{}",
			input:   map[string]interface{}{"key": "value", "number": 123},
			want:    map[string]interface{}{"key": "value", "number": 123},
			wantErr: false,
		},
		{
			name:    "Input é string JSON válida",
			input:   `{"key": "value", "nested": {"num": 456}}`,
			want:    map[string]interface{}{"key": "value", "nested": map[string]interface{}{"num": 456.0}}, // JSON numbers são float64 por padrão
			wantErr: false,
		},
		{
			name:    "Input é string JSON inválida",
			input:   `{"key": "value",}`, // JSON inválido com vírgula extra
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Input é tipo não suportado (int)",
			input:   12345,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Input é tipo não suportado (slice)",
			input:   []string{"a", "b"},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertToMap(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertToMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertToMap() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// --- Testes para extractTemplateVars ---

func TestExtractTemplateVars(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "String com uma variável",
			input: "Bearer ${step.requestAccessToken.outputVar}",
			want:  []string{"step.requestAccessToken.outputVar"},
		},
		{
			name:  "String com múltiplas variáveis",
			input: "/users/${input.body.userId}/product/${input.body.productSku}",
			want:  []string{"input.body.userId", "input.body.productSku"},
		},
		{
			name:  "String sem variáveis",
			input: "https://sts.company-hostname.com/token",
			want:  []string{}, // Espera slice vazio, não nil
		},
		{
			name:  "String com variável no início e fim",
			input: "${env:VAR_A} some text ${env:VAR_B}",
			want:  []string{"env:VAR_A", "env:VAR_B"},
		},
		{
			name:  "Variáveis adjacentes",
			input: "${var1}${var2}",
			want:  []string{"var1", "var2"},
		},
		{
			name:  "String com chaves não balanceadas ou vazias",
			input: "Test ${var1} and ${var2 and ${}",
			want:  []string{"var1"}, // Só extrai a válida
		},
		{
			name:  "String vazia",
			input: "",
			want:  []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractTemplateVars(tt.input)
			// Comparar slices vazios vs nil pode ser complicado, normalizar para comparação
			if len(got) == 0 && len(tt.want) == 0 {
				// Considerar iguais se ambos estão vazios
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractTemplateVars() = %v, want %v", got, tt.want)
			}
		})
	}
}

// --- Testes para getDeepValue ---

func TestGetDeepValue(t *testing.T) {
	// Simula o contexto de variáveis que o framework teria
	vars := map[string]interface{}{
		"input": map[string]interface{}{
			"source": "apiGatewayProxy",
			"body": map[string]interface{}{
				"userId":     "user-123",
				"productSku": "sku-abc",
			},
			"headers": map[string]interface{}{
				"Content-Type": "application/json",
			},
		},
		"step": map[string]interface{}{
			"validateRequestBody": map[string]interface{}{
				"status": "success",
			},
			"getUserDetails": map[string]interface{}{
				"outputVar": map[string]interface{}{
					"statusCode": 200,
					"body": map[string]interface{}{
						"id":     "user-123",
						"name":   "John Doe",
						"email":  "john.doe@example.com",
						"status": "active",
						"address": map[string]interface{}{ // Nível extra de profundidade
							"street": "123 Main St",
							"city":   "Anytown",
						},
					},
				},
			},
		},
		"env": map[string]interface{}{
			"TABLE_NAME": "MyDynamoTable",
		},
		"this": map[string]interface{}{
			"response": "placeholder for response", // Apenas para teste
		},
	}

	tests := []struct {
		name string
		path string
		data interface{} // Usar o 'vars' definido acima na maioria dos casos
		want interface{}
	}{
		{
			name: "Caminho simples nível 1",
			path: "input.source",
			data: vars,
			want: "apiGatewayProxy",
		},
		{
			name: "Caminho aninhado",
			path: "input.body.userId",
			data: vars,
			want: "user-123",
		},
		{
			name: "Caminho mais profundo",
			path: "step.getUserDetails.outputVar.body.address.city",
			data: vars,
			want: "Anytown",
		},
		{
			name: "Caminho para um mapa inteiro",
			path: "input.body",
			data: vars,
			want: map[string]interface{}{"userId": "user-123", "productSku": "sku-abc"},
		},
		{
			name: "Caminho para chave inexistente no final",
			path: "input.body.nonExistentKey",
			data: vars,
			want: nil,
		},
		{
			name: "Caminho com chave inexistente no meio",
			path: "step.nonExistentStep.outputVar",
			data: vars,
			want: nil,
		},
		{
			name: "Caminho com chave inválida (tentando acessar subnível de string)",
			path: "input.source.sublevel",
			data: vars,
			want: nil,
		},
		{
			name: "Caminho para variável de ambiente simulada",
			path: "env.TABLE_NAME",
			data: vars,
			want: "MyDynamoTable",
		},
		{
			name: "Caminho para variável 'this' (simulando)",
			path: "this.response",
			data: vars,
			want: "placeholder for response",
		},
		{
			name: "Caminho vazio",
			path: "",
			data: vars,
			want: nil,
		},
		{
			name: "Data é nil",
			path: "input.body.userId",
			data: nil,
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getDeepValue(tt.data, tt.path)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getDeepValue() got = %v (%T), want %v (%T)", got, got, tt.want, tt.want)
			}
		})
	}
}

// --- Testes para resolveStringTemplates ---

func TestResolveStringTemplates(t *testing.T) {
	// Contexto de variáveis similar ao anterior
	vars := map[string]interface{}{
		"input": map[string]interface{}{
			"body": map[string]interface{}{
				"userId":     "user-123",
				"productSku": "sku-abc",
				"quantity":   5, // Adiciona um número
			},
		},
		"step": map[string]interface{}{
			"requestAccessToken": map[string]interface{}{
				"outputVar": "xyz789token",
			},
			"getUserDetails": map[string]interface{}{
				"outputVar": map[string]interface{}{"id": "user-123", "name": "John Doe"},
			},
		},
		"env": map[string]interface{}{
			"API_KEY": "secret-key",
		},
		"constants": map[string]interface{}{
			"timeout": 10,
		},
	}

	tests := []struct {
		name string
		s    string
		vars map[string]interface{}
		want interface{} // Pode retornar o valor original se for template exato
	}{
		{
			name: "String com template exato (string)",
			s:    "${step.requestAccessToken.outputVar}",
			vars: vars,
			want: "xyz789token",
		},
		{
			name: "String com template exato (int)",
			s:    "${input.body.quantity}",
			vars: vars,
			want: 5, // Deve retornar o tipo original
		},
		{
			name: "String com template exato (map)",
			s:    "${step.getUserDetails.outputVar}",
			vars: vars,
			want: map[string]interface{}{"id": "user-123", "name": "John Doe"}, // Deve retornar o mapa
		},
		{
			name: "String com template embutido",
			s:    "Authorization: Bearer ${step.requestAccessToken.outputVar}",
			vars: vars,
			want: "Authorization: Bearer xyz789token",
		},
		{
			name: "String com múltiplos templates embutidos",
			s:    "User ID: ${input.body.userId}, SKU: ${input.body.productSku}",
			vars: vars,
			want: "User ID: user-123, SKU: sku-abc",
		},
		{
			name: "String com template numérico embutido",
			s:    "Order quantity: ${input.body.quantity}",
			vars: vars,
			want: "Order quantity: 5",
		},
		{
			name: "String sem templates",
			s:    "Fixed string value",
			vars: vars,
			want: "Fixed string value",
		},
		{
			name: "String com template para valor inexistente",
			s:    "Value: ${step.nonExistent.value}",
			vars: vars,
			want: "Value: <nil>", // Ou como fmt.Sprintf("%v", nil) renderiza
		},
		{
			name: "String vazia",
			s:    "",
			vars: vars,
			want: "",
		},
		{
			name: "Template referenciando env var",
			s:    "X-Api-Key: ${env.API_KEY}",
			vars: vars,
			want: "X-Api-Key: secret-key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := resolveStringTemplates(tt.s, tt.vars)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("resolveStringTemplates() got = %v (%T), want %v (%T)", got, got, tt.want, tt.want)
			}
		})
	}
}

// --- Testes para resolveMapTemplates ---
// (resolveMapTemplates chama ResolveTemplates recursivamente, que por sua vez chama resolveStringTemplates, etc.)
func TestResolveMapTemplates(t *testing.T) {
	// Contexto de variáveis
	vars := map[string]interface{}{
		"input": map[string]interface{}{
			"body": map[string]interface{}{
				"userId": "user-123",
				"data":   map[string]interface{}{"value": 10},
			},
			"pathParameters": map[string]interface{}{
				"orderId": "ord-999",
			},
		},
		"step": map[string]interface{}{
			"getUser": map[string]interface{}{"name": "Alice"},
		},
		"env": map[string]interface{}{
			"ENDPOINT": "https://service.example.com",
		},
	}

	tests := []struct {
		name string
		m    map[string]interface{}
		vars map[string]interface{}
		want map[string]interface{}
	}{
		{
			name: "Mapa com strings e templates",
			m: map[string]interface{}{
				"url":      "https://api.userservice.com/v1/users/${input.body.userId}",
				"static":   "value",
				"userName": "${step.getUser.name}",
				"complex":  "Order ${input.pathParameters.orderId} for ${step.getUser.name}",
			},
			vars: vars,
			want: map[string]interface{}{
				"url":      "https://api.userservice.com/v1/users/user-123",
				"static":   "value",
				"userName": "Alice",
				"complex":  "Order ord-999 for Alice",
			},
		},
		{
			name: "Mapa com valor sendo um template exato para um mapa",
			m: map[string]interface{}{
				"userData": "${step.getUser}", // Deve resolver para o mapa em si
			},
			vars: vars,
			want: map[string]interface{}{
				"userData": map[string]interface{}{"name": "Alice"},
			},
		},
		{
			name: "Mapa aninhado com templates",
			m: map[string]interface{}{
				"config": map[string]interface{}{
					"baseURL": "${env.ENDPOINT}",
					"path":    "/orders/${input.pathParameters.orderId}",
					"nested": map[string]interface{}{
						"user": "${step.getUser.name}",
					},
				},
				"timeout": 30,
			},
			vars: vars,
			want: map[string]interface{}{
				"config": map[string]interface{}{
					"baseURL": "https://service.example.com",
					"path":    "/orders/ord-999",
					"nested": map[string]interface{}{
						"user": "Alice",
					},
				},
				"timeout": 30, // Número não é template, permanece
			},
		},
		{
			name: "Mapa vazio",
			m:    map[string]interface{}{},
			vars: vars,
			want: map[string]interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// resolveMapTemplates chama ResolveTemplates internamente
			// mas para testar especificamente a lógica de mapa, chamamos diretamente
			got := resolveMapTemplates(tt.m, tt.vars)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("resolveMapTemplates() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// --- Testes para resolveSliceTemplates ---
// (resolveSliceTemplates chama ResolveTemplates recursivamente)
func TestResolveSliceTemplates(t *testing.T) {
	// Contexto de variáveis
	vars := map[string]interface{}{
		"input": map[string]interface{}{"userId": "u-456"},
		"step": map[string]interface{}{
			"processItems": map[string]interface{}{"result": []string{"item1", "item2"}},
			"config":       map[string]interface{}{"value": "configValue"},
		},
	}

	tests := []struct {
		name string
		s    []interface{}
		vars map[string]interface{}
		want []interface{}
	}{
		{
			name: "Slice com strings contendo templates",
			s: []interface{}{
				"User: ${input.userId}",
				"Fixed",
				"${step.config.value}",
			},
			vars: vars,
			want: []interface{}{
				"User: u-456",
				"Fixed",
				"configValue",
			},
		},
		{
			name: "Slice com template exato para slice",
			s: []interface{}{
				"Header",
				"${step.processItems.result}", // Espera-se que resolva para o slice
				"Footer",
			},
			vars: vars,
			want: []interface{}{
				"Header",
				[]string{"item1", "item2"}, // O valor resolvido é o slice
				"Footer",
			},
		},
		{
			name: "Slice com mapas contendo templates",
			s: []interface{}{
				map[string]interface{}{"id": 1, "user": "${input.userId}"},
				map[string]interface{}{"id": 2, "config": "${step.config.value}"},
			},
			vars: vars,
			want: []interface{}{
				map[string]interface{}{"id": 1, "user": "u-456"},
				map[string]interface{}{"id": 2, "config": "configValue"},
			},
		},
		{
			name: "Slice vazio",
			s:    []interface{}{},
			vars: vars,
			want: []interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := resolveSliceTemplates(tt.s, tt.vars)
			// DeepEqual lida bem com slices e tipos mistos
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("resolveSliceTemplates() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// --- Testes para getSafeValueFromPath ---

func TestGetSafeValueFromPath(t *testing.T) {
	// Contexto de variáveis
	vars := map[string]interface{}{
		"input": map[string]interface{}{
			"body": map[string]interface{}{
				"userId":     "user-123",
				"productSku": "sku-abc",
				"quantity":   5,
				"enabled":    true,
				"tags":       []string{"A", "B"},
			},
			"pathParameters": map[string]interface{}{ // Adicionado para teste
				"orderId": "ord-555",
			},
			"source": "apiGatewayProxy",
		},
		"step": map[string]interface{}{
			"getUserDetails": map[string]interface{}{
				"outputVar": map[string]interface{}{
					"id":      "user-123",
					"name":    "John Doe",
					"address": map[string]interface{}{"city": "Anytown"},
				},
			},
		},
	}

	tests := []struct {
		name string
		path string
		vars map[string]interface{}
		want string
	}{
		{
			name: "Caminho para string",
			path: "input.body.userId",
			vars: vars,
			want: "user-123",
		},
		{
			name: "Caminho para número",
			path: "input.body.quantity",
			vars: vars,
			want: "5",
		},
		{
			name: "Caminho para booleano",
			path: "input.body.enabled",
			vars: vars,
			want: "true",
		},
		// {
		// 	name: "Caminho para slice (JSON)",
		// 	path: "input.body.tags",
		// 	vars: vars,
		// 	want: `["A","B"]`,
		// },
		{
			name: "Caminho para mapa (JSON)",
			path: "step.getUserDetails.outputVar.address",
			vars: vars,
			want: `{"city":"Anytown"}`,
		},
		{
			name: "Caminho inexistente no final",
			path: "input.body.nonexistent",
			vars: vars,
			want: "",
		},
		{
			name: "Caminho inexistente no meio",
			path: "step.nonexistent.outputVar",
			vars: vars,
			want: "",
		},
		{
			name: "Caminho sobre tipo inválido (string)",
			path: "input.source.sublevel", // source é string
			vars: vars,
			want: "",
		},
		{
			name: "Caso especial pathParameters existente",
			path: "input.pathParameters.orderId",
			vars: vars,
			want: "ord-555",
		},
		{
			name: "Caso especial pathParameters inexistente (chave)",
			path: "input.pathParameters.customerId",
			vars: vars,
			want: "",
		},
		{
			name: "Caso especial pathParameters inexistente (mapa)",
			path: "input.pathParameters.orderId",
			vars: map[string]interface{}{"input": map[string]interface{}{"body": "test"}}, // input existe, mas sem pathParameters
			want: "",
		},
		{
			name: "Vars é nil",
			path: "input.body.userId",
			vars: nil,
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getSafeValueFromPath(tt.path, tt.vars)
			if got != tt.want {
				t.Errorf("getSafeValueFromPath() got = %q, want %q", got, tt.want)
			}
		})
	}
}

// --- Testes para findFieldByJSONTag ---

// Estrutura de exemplo para testar findFieldByJSONTag
type SampleStruct struct {
	UserName  string `json:"name"`
	UserAge   int    `json:"age,omitempty"`
	IsEnabled bool   `json:"enabled"`
	NoTag     string
}

func TestFindFieldByJSONTag(t *testing.T) {
	sample := SampleStruct{UserName: "Test", UserAge: 30, IsEnabled: true, NoTag: "Untagged"}
	val := reflect.ValueOf(sample) // Obter reflect.Value da struct

	tests := []struct {
		name        string
		structValue reflect.Value
		jsonTag     string
		wantValid   bool         // Queremos um reflect.Value válido?
		wantKind    reflect.Kind // Se válido, qual o Kind esperado?
	}{
		{
			name:        "Encontrar tag JSON simples ('name')",
			structValue: val,
			jsonTag:     "name",
			wantValid:   true,
			wantKind:    reflect.String,
		},
		{
			name:        "Encontrar tag JSON com opção ('age,omitempty')",
			structValue: val,
			jsonTag:     "age", // Deve encontrar pela parte antes da vírgula
			wantValid:   true,
			wantKind:    reflect.Int,
		},
		{
			name:        "Encontrar tag JSON booleana ('enabled')",
			structValue: val,
			jsonTag:     "enabled",
			wantValid:   true,
			wantKind:    reflect.Bool,
		},
		{
			name:        "Tag JSON não encontrada ('nonexistent')",
			structValue: val,
			jsonTag:     "nonexistent",
			wantValid:   false,
			wantKind:    reflect.Invalid,
		},
		{
			name:        "Tentar encontrar campo sem tag JSON pelo nome do campo ('NoTag')",
			structValue: val,
			jsonTag:     "NoTag", // Não é uma tag JSON
			wantValid:   false,
			wantKind:    reflect.Invalid,
		},
		// {
		// 	name:        "Tag JSON vazia",
		// 	structValue: val,
		// 	jsonTag:     "",
		// 	wantValid:   false,
		// 	wantKind:    reflect.Invalid,
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Garantir que estamos testando com uma struct
			if tt.structValue.Kind() != reflect.Struct {
				t.Skip("Skipping test case that requires a struct value.")
				return
			}
			got := findFieldByJSONTag(tt.structValue, tt.jsonTag)

			isValid := got.IsValid()
			if isValid != tt.wantValid {
				t.Errorf("findFieldByJSONTag() validity got = %v, want %v", isValid, tt.wantValid)
				return // Não checar Kind se a validade estiver errada
			}

			// Só checa o Kind se esperamos um campo válido
			if tt.wantValid && got.Kind() != tt.wantKind {
				t.Errorf("findFieldByJSONTag() kind got = %v, want %v", got.Kind(), tt.wantKind)
			}
		})
	}
}
