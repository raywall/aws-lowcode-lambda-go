package loader

// import (
// 	"testing"

// 	"github.com/go-playground/assert/v2"
// )

// func TestLambdaConfigTemplate(t *testing.T) {
// 	loader := LocalLoader{Path: local}
// 	data, _ := loader.Load()

// 	t.Run("Deve retornar a versão de template 2025-04-19", func(t *testing.T) {
// 		assert.Equal(t, data.TemplateFormatVersion, "2025-04-19")
// 	})

// 	t.Run("Deve retornar o transform do lowcode lambda: LC::Lambda", func(t *testing.T) {
// 		assert.Equal(t, data.Transform, "LC::Lambda")
// 	})
// }

// func TestLambdaConfigMetadata(t *testing.T) {
// 	loader := LocalLoader{Path: local}
// 	data, _ := loader.Load()

// 	t.Run("Deve retornar o nome do serviço", func(t *testing.T) {
// 		assert.Equal(t, data.Metadata.Name, "user-service")
// 	})
// }
