package security

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

// regexes list
var (
	// AWS Access Key ID
	reAWS = regexp.MustCompile(`\bAKIA[0-9A-Z]{16}\b`)
	// OpenAI (sk-...)
	reOpenAI = regexp.MustCompile(`\bsk-[a-zA-Z0-9]{20,}\b`)
	// Private Key Headers
	rePrivateKey = regexp.MustCompile(`-----BEGIN [A-Z]+ PRIVATE KEY-----`)
)

type Scanner struct {
	Disabled bool
}

func New(disabled bool) *Scanner {
	return &Scanner{Disabled: disabled}
}

func (s *Scanner) Scan(path string, content []byte) error {
	if s.Disabled {
		return nil
	}

	// 1. Filename Check
	base := filepath.Base(path)
	if isSensitiveFilename(base) {
		return fmt.Errorf("sensitive filename detected: %s", base)
	}

	// 2. Content Check
	if reAWS.Match(content) {
		return fmt.Errorf("potential AWS Access Key detected")
	}
	if reOpenAI.Match(content) {
		return fmt.Errorf("potential OpenAI Key detected")
	}
	if rePrivateKey.Match(content) {
		return fmt.Errorf("private key header detected")
	}

	return nil
}

func isSensitiveFilename(name string) bool {
	switch name {
	case ".env", ".env.local", ".env.production", "id_rsa", "id_dsa", "id_ed25519":
		return true
	}

	ext := filepath.Ext(name)
	switch ext {
	case ".pem", ".key", ".p12", ".pfx", ".kdbx":
		return true
	}

	if strings.Contains(name, "_secret") || strings.Contains(name, "_token") {
		// now empty, bcs dangerous
	}

	return false
}
