package skeleton

import (
	"testing"
)

func TestBraceScanner(t *testing.T) {
	content := `
function test() {
    console.log("hello {");
    if (true) {
        // block
    }
}
`
	blocks := scanBraces([]byte(content))
	if len(blocks) != 2 {
		t.Errorf("expected 2 blocks, got %d", len(blocks))
	}
	
	// block 1: if (true) { ... }
	if blocks[0].StartLine != 4 || blocks[0].EndLine != 6 {
		t.Errorf("block 0 unexpected: %+v", blocks[0])
	}
	
	// block 2: function test() { ... }
	if blocks[1].StartLine != 2 || blocks[1].EndLine != 7 {
		t.Errorf("block 1 unexpected: %+v", blocks[1])
	}
}

func TestIndentationScanner(t *testing.T) {
	content := `
class Test:
    def method(self):
        print("hello")
    def other(self):
        pass
`
	blocks := scanIndentation([]byte(content))
	if len(blocks) != 3 {
		t.Errorf("expected 3 blocks, got %d", len(blocks))
	}
	
	// Inner blocks should come first (depth-first or order-of-completion)
	// method(self)
	if blocks[0].StartLine != 3 || blocks[0].EndLine != 4 {
		t.Errorf("block 0 unexpected: %+v", blocks[0])
	}
	// other(self)
	if blocks[1].StartLine != 5 || blocks[1].EndLine != 6 {
		t.Errorf("block 1 unexpected: %+v", blocks[1])
	}
	// class Test
	if blocks[2].StartLine != 2 || blocks[2].EndLine != 6 {
		t.Errorf("block 2 unexpected: %+v", blocks[2])
	}
}
