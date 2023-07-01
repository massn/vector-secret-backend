package main

import (
	"os"
	"testing"
)

func TestToSecretsInfo(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		wantSecretInfo SecretsInfo
		wantErr        string
	}{
		{
			name:  "valid input",
			input: `{"version":"1.0","secrets":["secret1","secret2"]}`,
			wantSecretInfo: SecretsInfo{
				Version: "1.0",
				Secrets: []string{
					"secret1",
					"secret2",
				},
			},
			wantErr: "",
		},
		{
			name:  "bad json input",
			input: `{"version":"1.0","secrets":["secret1","secret2"]`,
			wantSecretInfo: SecretsInfo{
				Version: "",
				Secrets: []string{},
			},
			wantErr: "unexpected end of JSON input",
		},
		{
			name:  "empty input",
			input: ``,
			wantSecretInfo: SecretsInfo{
				Version: "",
				Secrets: []string{},
			},
			wantErr: "unexpected end of JSON input",
		},
		{
			name:  "no secretes in input",
			input: `{"version":"1.0","secrets":[]}`,
			wantSecretInfo: SecretsInfo{
				Version: "1.0",
				Secrets: []string{},
			},
			wantErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSecretInfo, err := toSecretsInfo(tt.input)
			if err != nil && err.Error() != tt.wantErr {
				t.Errorf("toSecretsInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotSecretInfo.Version != tt.wantSecretInfo.Version {
				t.Errorf("toSecretsInfo() gotSecretInfo.Version = %v, want %v", gotSecretInfo.Version, tt.wantSecretInfo.Version)
			}
			if len(gotSecretInfo.Secrets) != len(tt.wantSecretInfo.Secrets) {
				t.Errorf("toSecretsInfo() gotSecretInfo.Secrets = %v, want %v", gotSecretInfo.Secrets, tt.wantSecretInfo.Secrets)
			}
		})
	}
}

func TestGetResponse(t *testing.T) {
	tests := []struct {
		name       string
		envVals    map[string]string
		secretKeys []string
		wantRes    string
	}{
		{
			name: "valid input",
			envVals: map[string]string{
				"VECTOR_SECRET_SECRET1": "value1",
				"VECTOR_SECRET_SECRET2": "value2",
			},
			secretKeys: []string{
				"secret1",
				"secret2",
			},
			wantRes: `{"secret1":{"value":"value1","error":null},"secret2":{"value":"value2","error":null}}`,
		},
		{
			name:    "invalid input with no env variables",
			envVals: map[string]string{},
			secretKeys: []string{
				"secret1",
				"secret2",
			},
			wantRes: `{"secret1":{"value":"","error":"secret secret1 not found"},"secret2":{"value":"","error":"secret secret2 not found"}}`,
		},
		{
			name: "invalid input with one env variable",
			envVals: map[string]string{
				"VECTOR_SECRET_SECRET1": "value1",
			},
			secretKeys: []string{
				"secret1",
				"secret2",
			},
			wantRes: `{"secret1":{"value":"value1","error":null},"secret2":{"value":"","error":"secret secret2 not found"}}`,
		},
	}

	for _, tt := range tests {
		for k, v := range tt.envVals {
			err := os.Setenv(k, v)
			if err != nil {
				panic(err)
			}
		}
		t.Run(tt.name, func(t *testing.T) {
			gotRes := getResponse(tt.secretKeys)
			if gotRes != tt.wantRes {
				t.Errorf("getResponse() gotRes = %v, want %v", gotRes, tt.wantRes)
			}
		})
		for k := range tt.envVals {
			err := os.Unsetenv(k)
			if err != nil {
				panic(err)
			}
		}
	}
}
