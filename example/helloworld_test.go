package example

import "testing"

func TestHelloWorld(t *testing.T) {
	if HelloWorld() != "Hello, World!" {
		t.Error("Test failed")
	}
}
