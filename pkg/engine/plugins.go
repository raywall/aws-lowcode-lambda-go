package engine

// import (
// 	"plugin"

// 	"github.com/raywall/aws-lowcode-lambda-go/pkg/plugin"
// )

// func LoadPlugins(pluginPaths []string) (*plugin.Registry, error) {
// 	registry := plugin.NewRegistry()

// 	for _, path := range pluginPaths {
// 		p, err := plugin.Open(path)
// 		if err != nil {
// 			return nil, err
// 		}

// 		sym, err := p.Lookup("Plugin")
// 		if err != nil {
// 			continue
// 		}

// 		if plugin, ok := sym.(plugin.Plugin); ok {
// 			plugin.Register(registry)
// 		}
// 	}

// 	return registry, nil
// }
