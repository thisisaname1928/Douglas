package app

import (
	"fmt"
	"testing"
)

func TestApp(t *testing.T) {
	fmt.Println(call4AI("Hello", GEMINI_API_KEY))
}
