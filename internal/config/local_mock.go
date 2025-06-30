package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
)

type LocalMockProvider struct {
	values map[string]string
}

func NewLocalMockProviderFromFile(filePath string) (*LocalMockProvider, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open local config file: %v", err)
	}
	defer f.Close()

	values := make(map[string]string)
	if err := json.NewDecoder(f).Decode(&values); err != nil {
		return nil, fmt.Errorf("failed to parse local config JSON: %v", err)
	}

	return &LocalMockProvider{values: values}, nil
}

// NewConfigEventFromLocalMock 将 mockProvider.values 直接映射成 ConfigEvent
func NewConfigEventFromLocalMock(provider *LocalMockProvider, cfgEvent *ConfigEvent) error {
	raw, err := json.Marshal(provider.values)
	if err != nil {
		return fmt.Errorf("failed to marshal mock values: %v", err)
	}

	if err := json.Unmarshal(raw, cfgEvent); err != nil {
		return fmt.Errorf("failed to unmarshal mock values into ConfigEvent: %v", err)
	}

	return nil
}

func (p *LocalMockProvider) GetParameter(ctx context.Context, paramName string) (string, error) {
	if v, ok := p.values[paramName]; ok {
		return v, nil
	}
	return "", fmt.Errorf("mock value for %s not found", paramName)
}
