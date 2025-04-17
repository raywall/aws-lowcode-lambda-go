package customPlugins

// import (
// 	"github.com/aws/aws-lambda-go/lambda"
// 	"github.com/gorilla/mux"

// 	"github.com/raywall/aws-lowcode-lambda-go/pkg/engine"
// )

// func init() {
// 	// Carrega plugins
// 	registry := engine.LoadPlugins([]string{
// 		"/plugins/email.go",
// 		"/plugins/database.go",
// 	})

// 	// Inicializa framework com plugins
// 	engine.Init(registry)

// 	// Configura roteador
// 	router := mux.NewRouter()
// 	router.Use(engine.LambdaMiddleware)

// 	// Carrega rotas do YAML
// 	engine.LoadRoutesFromConfig("config.yaml")
// }

// func main() {
// 	lambda.Start(engine.Handler)
// }
