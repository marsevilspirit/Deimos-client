package testutil

import (
	"fmt"
	"log"
)

// MockT 实现 assert.TestingT 接口，用于在 main 函数中使用断言
type MockT struct {
	FailOnError bool // 是否在断言失败时退出程序
}

// NewMockT 创建一个新的 MockT 实例
func NewMockT(failOnError bool) *MockT {
	return &MockT{FailOnError: failOnError}
}

func (t *MockT) Errorf(format string, args ...interface{}) {
	fmt.Printf("❌ 断言失败: "+format+"\n", args...)
}

func (t *MockT) FailNow() {
	if t.FailOnError {
		log.Fatal("断言失败，程序退出")
	} else {
		fmt.Println("⚠️  断言失败，但继续执行")
	}
}

// Helper 返回 true，用于兼容 testify 接口
func (t *MockT) Helper() {}
