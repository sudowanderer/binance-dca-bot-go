package env

import "os"

// IsLambdaEnvironment 检测当前环境是否在 AWS Lambda 中
func IsLambdaEnvironment() bool {
	for _, env := range os.Environ() {
		if len(env) > 11 && env[:11] == "AWS_LAMBDA_" {
			return true
		}
	}
	return false
}
