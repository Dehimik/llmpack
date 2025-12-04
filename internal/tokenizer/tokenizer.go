package tokenizer

import (
	"github.com/pkoukk/tiktoken-go"
)

type TikToken struct {
	tk *tiktoken.Tiktoken
}

func New() *TikToken {
	// encoding GPT-4o / GPT-4
	tk, _ := tiktoken.GetEncoding("cl100k_base")
	return &TikToken{tk: tk}
}

func (t *TikToken) Count(text string) int {
	if t.tk == nil {
		return 0
	}
	return len(t.tk.Encode(text, nil, nil))
}
