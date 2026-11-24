package voice

import (
	"bytes"
	"testing"
)

func TestNoopEncrypter_Decrypt(t *testing.T) {
	enc := NewNoopEncrypter()

	// Short packet (< RTPHeaderSize) returns an error instead of panicking
	if _, err := enc.Decrypt(RTPHeaderSize, []byte{1, 2, 3}); err == nil {
		t.Fatal("expected error for short packet, got nil")
	}

	// Valid packet
	expected := []byte("1234567890")
	packet := make([]byte, RTPHeaderSize+len(expected))
	copy(packet[RTPHeaderSize:], expected)

	decrypted, err := enc.Decrypt(RTPHeaderSize, packet)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Equal(expected, decrypted) {
		t.Fatalf("expected %v, got %v", expected, decrypted)
	}
}
