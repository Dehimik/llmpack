package skeleton

import (
	"strings"
	"testing"
)

func TestScannerStringEscaping(t *testing.T) {
	content := []byte(`
func foo() {
    s := "\\\\"
}
func bar() {}
`)
	blocks := scanBraces(content)
	t.Logf("Blocks found: %+v", blocks)
	
	foundBar := false
	for _, b := range blocks {
		if b.StartLine == 5 {
			foundBar = true
		}
	}
	if !foundBar {
		t.Errorf("Could not find block starting at line 5 (bar), string escaping might be broken")
	}
}

func TestPythonCommentIndentation(t *testing.T) {
	content := []byte(`
class MyClass:
    def method1(self):
        pass
    # This comment should not close the class block if it has less indent
# But this line should
def other():
    pass
`)
	blocks := scanIndentation(content)
	
	var myClassBlock *Block
	for _, b := range blocks {
		if b.StartLine == 2 {
			myClassBlock = &b
		}
	}
	
	if myClassBlock == nil {
		t.Fatal("MyClass block not found")
	}
	
	// If comment at line 5 closes the block, EndLine will be 4 or 5.
	if myClassBlock.EndLine < 6 {
		t.Errorf("MyClass block closed too early at line %d", myClassBlock.EndLine)
	}
}

func TestExtractBraceSymbolsControlFlow(t *testing.T) {
	content := []byte(`
class MyClass {
    myMethod() {
        if (true) {
            doSomething();
        }
    }
}
`)
	symbols, err := extractBraceSymbols(content)
	if err != nil {
		t.Fatal(err)
	}
	
	for _, s := range symbols {
		if s.Name == "if" {
			t.Errorf("Should not extract 'if' as a symbol")
		}
	}
}

func TestGenericReducerClassSkeleton(t *testing.T) {
	content := []byte(`
class MyClass {
    method1() {
        console.log("hello");
    }
    method2() {
        return 42;
    }
}
`)
	reduced := string(reduceBraces(content, ""))
	if !strings.Contains(reduced, "method1()") {
		t.Errorf("Class skeleton should contain method1(), but got:\n%s", reduced)
	}
	if !strings.Contains(reduced, "method2()") {
		t.Errorf("Class skeleton should contain method2(), but got:\n%s", reduced)
	}
}

func TestSorting(t *testing.T) {
	content := []byte(`
def second():
    pass
def first():
    pass
`)
	symbols, err := extractIndentationSymbols(content)
	if err != nil {
		t.Fatal(err)
	}
	if len(symbols) < 2 {
		t.Fatalf("Expected 2 symbols, got %d", len(symbols))
	}
	// They should be sorted by line number: second (line 2) then first (line 4)
	if symbols[0].Name != "second" || symbols[1].Name != "first" {
		t.Errorf("Symbols should be sorted by line number, but got %s then %s", symbols[0].Name, symbols[1].Name)
	}
}
