package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type SecretsInfo struct {
	Version string   `json:"version"`
	Secrets []string `json:"secrets"`
}

type ValErr struct {
	Value string  `json:"value"`
	Error *string `json:"error"`
}

func main() {
	if len(os.Args) < 2 {
		panic("No secrets provided")
	}
	secretsInfo, err := toSecretsInfo(os.Args[1])
	if err != nil {
		panic(err)
	}
	if secretsInfo.Version != "1.0" {
		panic("Invalid version")
	}
	res := getResponse(secretsInfo.Secrets)
	fmt.Println(res)
}

func toSecretsInfo(input string) (*SecretsInfo, error) {
	var si SecretsInfo
	err := json.Unmarshal([]byte(input), &si)
	if err != nil {
		return &SecretsInfo{}, err
	}
	return &si, nil
}

func getResponse(secretKeys []string) string {
	res := make(map[string]ValErr)
	for _, key := range secretKeys {
		val, err := getSecret(key)
		valErr := ValErr{Value: val}
		if err != nil {
			errString := err.Error() // helper variable to get pointer to string
			valErr.Error = &errString
		} else {
			valErr.Error = nil
		}
		res[key] = valErr
	}
	jsonRes, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}
	return string(jsonRes)
}

func getSecret(key string) (string, error) {
	envKey := "VECTOR_SECRET_" + strings.ToUpper(key)
	secret, ok := os.LookupEnv(envKey)
	if !ok {
		return "", fmt.Errorf("secret %s not found", key)
	}
	return secret, nil
}
