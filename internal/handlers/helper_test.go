package handlers

import (
	"fmt"
	"testing"
)

func TestGenerateTokens(t *testing.T) {
	token, _ := GenerateTokens("alex@gmail.com", false, 1)
	fmt.Println(token)
}
