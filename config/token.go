package config

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type (
	Request struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		GrantType    string `json:"grant_type"`
	}

	Token struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		Scope        string `json:"scope"`
		Active       bool   `json:"active"`
	}

	Credential struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
	}
)

var (
	EndpointSTS   string
	CredentialSTS Credential
)

type STS interface {
	GetToken()
}

func (res *ResourceItem) GetToken() (*string, error) {
	if res.ResourceType != "RESTfulApi" {
		return nil, errors.New("resource is not a database")
	}

	client := &http.Client{}
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	// url sts
	stsUrl := EndpointSTS + "api/oauth/token"

	data := url.Values{}
	data.Set("client_id", CredentialSTS.ClientID)
	data.Set("client_secret", CredentialSTS.ClientSecret)
	data.Set("grant_type", "client_credentials")

	form := data.Encode()
	request, err := http.NewRequest("GET", stsUrl, strings.NewReader(form))
	if err != nil {
		log.Fatalf("falha ao gerar uma nova requisição: %v", err)
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response, err := client.Do(request)
	if err != nil {
		log.Fatalf("falha ao resgatar o token STS: %v", err)
	}

	var token Token
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("falha ao ler a resposta da requisição do token: %v", err)
	}

	err = json.Unmarshal(responseBody, &token)
	if err != nil {
		log.Fatalf("falha ao converter resposta da requisição do token: %v", err)
	}

	response.Body.Close()

	return &token.AccessToken, nil
}
