package skeleton

import (
	"testing"
)

func TestExtractSymbolsGo(t *testing.T) {
	content := `
package main
type User struct { Name string }
func (u User) SayHello() { println("hi") }
func main() { SayHello() }
`
	symbols, err := ExtractSymbols("test.go", []byte(content))
	if err != nil {
		t.Fatal(err)
	}
	
	if len(symbols) != 3 {
		t.Errorf("expected 3 symbols, got %d", len(symbols))
	}
	
	// Struct
	if symbols[0].Name != "User" || symbols[0].Type != "struct" {
		t.Errorf("expected User struct, got %+v", symbols[0])
	}
	
	// Method
	if symbols[1].Name != "SayHello" || symbols[1].Type != "method" || symbols[1].Parent != "User" {
		t.Errorf("expected SayHello method of User, got %+v", symbols[1])
	}
	
	// Function
	if symbols[2].Name != "main" || symbols[2].Type != "function" {
		t.Errorf("expected main function, got %+v", symbols[2])
	}
}

func TestExtractSymbolsPython(t *testing.T) {
	content := `
class Test:
    def method(self):
        pass
def top_level():
    pass
`
	symbols, err := ExtractSymbols("test.py", []byte(content))
	if err != nil {
		t.Fatal(err)
	}
	
	if len(symbols) != 3 {
		t.Errorf("expected 3 symbols, got %d", len(symbols))
	}
	
	if symbols[0].Name != "Test" || symbols[0].Type != "class" {
		t.Errorf("expected Test class, got %+v", symbols[0])
	}
	if symbols[1].Name != "method" || symbols[1].Type != "function" || symbols[1].Parent != "Test" {
		t.Errorf("expected method of Test, got %+v", symbols[1])
	}
	if symbols[2].Name != "top_level" || symbols[2].Type != "function" {
		t.Errorf("expected top_level function, got %+v", symbols[2])
	}
}

func TestExtractSymbolsBraces(t *testing.T) {
	content := `
class MyClass {
    myMethod() {
    }
}
function topLevel() {
}
`
	symbols, err := ExtractSymbols("test.ts", []byte(content))
	if err != nil {
		t.Fatal(err)
	}
	
	if len(symbols) != 3 {
		t.Errorf("expected 3 symbols, got %d", len(symbols))
	}
	
	if symbols[0].Name != "MyClass" || symbols[0].Type != "class" {
		t.Errorf("expected MyClass class, got %+v", symbols[0])
	}
	if symbols[1].Name != "myMethod" || symbols[1].Type != "method" || symbols[1].Parent != "MyClass" {
		t.Errorf("expected myMethod of MyClass, got %+v", symbols[1])
	}
	if symbols[2].Name != "topLevel" || symbols[2].Type != "function" {
		t.Errorf("expected topLevel function, got %+v", symbols[2])
	}
}
