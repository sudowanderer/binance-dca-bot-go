package env

import (
	"os"
	"testing"
)

func TestIsLambdaEnvironment(t *testing.T) {
	// 保存原始环境变量
	originalEnv := os.Environ()

	// 清理环境变量，模拟本地环境
	os.Clearenv()
	if IsLambdaEnvironment() {
		t.Errorf("Expected false in local environment")
	}

	// 模拟 Lambda 环境
	os.Setenv("AWS_LAMBDA_FUNCTION_NAME", "my-lambda")
	if !IsLambdaEnvironment() {
		t.Errorf("Expected true in Lambda environment")
	}

	// 恢复环境变量
	restoreEnv(originalEnv)
}

// restoreEnv 重新设置原始环境变量
func restoreEnv(env []string) {
	os.Clearenv()
	for _, kv := range env {
		parts := splitOnce(kv, '=')
		if len(parts) == 2 {
			os.Setenv(parts[0], parts[1])
		}
	}
}

func splitOnce(s string, sep byte) []string {
	for i := 0; i < len(s); i++ {
		if s[i] == sep {
			return []string{s[:i], s[i+1:]}
		}
	}
	return []string{s}
}
