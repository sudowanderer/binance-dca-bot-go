package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

// ParameterProvider defines how to fetch a parameter value by name.
type ParameterProvider interface {
	// GetParameter fetches the *decrypted* parameter value.
	GetParameter(ctx context.Context, name string) (string, error)
}

// SSMParameterProvider fetches from AWS SSM (with decryption).
type SSMParameterProvider struct {
	Client *ssm.Client
}

func (p *SSMParameterProvider) GetParameter(ctx context.Context, name string) (string, error) {
	out, err := p.Client.GetParameter(ctx, &ssm.GetParameterInput{
		Name:           &name,
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return "", fmt.Errorf("ssm GetParameter %q: %w", name, err)
	}
	if out.Parameter == nil || out.Parameter.Value == nil {
		return "", fmt.Errorf("ssm GetParameter %q: empty parameter value", name)
	}
	return *out.Parameter.Value, nil
}

// LocalMockProvider loads all parameters from a local JSON file.
type LocalMockProvider struct {
	values map[string]string
}

// NewLocalMockProviderFromFile reads a JSON file containing a map[string]string.
func NewLocalMockProviderFromFile(path string) (*LocalMockProvider, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading local config file %q: %w", path, err)
	}
	var data map[string]string
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, fmt.Errorf("unmarshaling local config JSON: %w", err)
	}
	return &LocalMockProvider{values: data}, nil
}

func (m *LocalMockProvider) GetParameter(ctx context.Context, name string) (string, error) {
	v, ok := m.values[name]
	if !ok {
		return "", fmt.Errorf("local mock: parameter %q not found", name)
	}
	return v, nil
}
