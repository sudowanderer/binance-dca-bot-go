package config

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

// SSMParameterProvider 用于生产环境，从 Parameter Store 拉取参数
type SSMParameterProvider struct {
	Client *ssm.Client
}

func (p *SSMParameterProvider) GetParameter(ctx context.Context, paramName string) (string, error) {
	resp, err := p.Client.GetParameter(ctx, &ssm.GetParameterInput{
		Name:           aws.String(paramName),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return "", err
	}
	return *resp.Parameter.Value, nil
}
