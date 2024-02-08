package config

import "errors"

func (res *ResourceItem) GetToken() (*string, error) {
	if res.ResourceType != "RESTfulApi" {
		return nil, errors.New("resource is not a database")
	}

	// logica para autenticacao aqui

	return nil, nil
}
