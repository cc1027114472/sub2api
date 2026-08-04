package payment

import "testing"

func TestLooksLikeCiphertext(t *testing.T) {
	t.Parallel()

	key := make([]byte, AES256KeySize)
	for i := range key {
		key[i] = byte(i + 1)
	}
	enc, err := Encrypt(`{"a":"b"}`, key)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	if !LooksLikeCiphertext(enc) {
		t.Fatalf("expected ciphertext to match LooksLikeCiphertext: %q", enc)
	}
	if LooksLikeCiphertext(`{"a":"b"}`) {
		t.Fatal("plaintext JSON must not look like ciphertext")
	}
	if LooksLikeCiphertext("not-cipher") {
		t.Fatal("garbage must not look like ciphertext")
	}
}
