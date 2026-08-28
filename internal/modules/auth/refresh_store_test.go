package auth

import "testing"

func TestRandomTokenUniqueAndSized(t *testing.T) {
	seen := map[string]bool{}
	for i := 0; i < 100; i++ {
		tok, err := randomToken(32)
		if err != nil {
			t.Fatalf("randomToken: %v", err)
		}
		if len(tok) != 64 { // 32 bytes -> 64 hex chars
			t.Fatalf("len = %d, want 64", len(tok))
		}
		if seen[tok] {
			t.Fatalf("duplicate token generated: %s", tok)
		}
		seen[tok] = true
	}
}
