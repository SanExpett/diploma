package requestId

import (
	"fmt"
	"testing"
)

func TestGenerateRequestID(t *testing.T) {
	fmt.Println(GenerateRequestID())
}
