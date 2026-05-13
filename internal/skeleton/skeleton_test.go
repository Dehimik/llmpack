package skeleton

import (
	"bytes"
	"testing"
)

func TestReduceBraces(t *testing.T) {
	input := `
function test() {
    console.log("hello");
}

class MyClass {
    method() {
        return 1;
    }
}
`
	output := reduceBraces([]byte(input))
	if !bytes.Contains(output, []byte("`... implementation hidden ...`")) {
		t.Errorf("Expected placeholder in output, got:\n%s", string(output))
	}
	if bytes.Contains(output, []byte("console.log")) {
		t.Errorf("Expected function body to be stripped, got:\n%s", string(output))
	}
}

func TestReduceIndentation(t *testing.T) {
	input := `
def my_func():
    print("hello")
    return True

class MyClass:
    def method(self):
        pass
`
	output := reduceIndentation([]byte(input))
	if !bytes.Contains(output, []byte("/* ... implementation hidden ... */")) {
		t.Errorf("Expected placeholder in output, got:\n%s", string(output))
	}
	if bytes.Contains(output, []byte("print(\"hello\")")) {
		t.Errorf("Expected function body to be stripped, got:\n%s", string(output))
	}
}

func TestProcess(t *testing.T) {
	goCode := "package main\nfunc main() { println(1) }"
	out, err := Process("test.go", []byte(goCode))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(out, []byte("`... implementation hidden ...`")) {
		t.Errorf("Go skeletonization failed")
	}

	pyCode := "def test():\n    pass"
	out, _ = Process("test.py", []byte(pyCode))
	if !bytes.Contains(out, []byte("/* ... implementation hidden ... */")) {
		t.Errorf("Python skeletonization failed")
	}

	tsCode := "function test() { return 1; }"
	out, _ = Process("test.ts", []byte(tsCode))
	if !bytes.Contains(out, []byte("`... implementation hidden ...`")) {
		t.Errorf("TypeScript skeletonization failed")
	}
}
