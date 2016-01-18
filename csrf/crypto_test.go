package csrf

import (
	"bytes"
	"crypto/rand"
	"errors"
	"testing"
)

// A reader that always fails on Read()
// Suitable for testing the case of crypto/rand unavailability
type failReader struct{}

func (f failReader) Read(p []byte) (n int, err error) {
	err = errors.New("dummy error")
	return
}

func TestChecksForPRNG(t *testing.T) {
	// Monkeypatch crypto/rand with an always-failing reader
	oldReader := rand.Reader
	rand.Reader = failReader{}
	// Restore it later for other tests
	defer func() {
		rand.Reader = oldReader
	}()

	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("Expected checkForPRNG() to panic")
		}
	}()

	checkForPRNG()
}

func TestGeneratesAValidToken(t *testing.T) {
	// We can't test much with any certainity here,
	// since we generate tokens randomly
	// Basically we check that the length of the
	// token is what it should be

	token := generateToken()
	l := len(token)

	if l != tokenLength {
		t.Errorf("Bad decoded token length: expected %d, got %d", tokenLength, l)
	}
}

func TestVerifyTokenChecksLengthCorrectly(t *testing.T) {
	for i := 0; i < 64; i++ {
		slice := make([]byte, i)
		result := verifyToken(slice, slice)
		if result != false {
			t.Errorf("VerifyToken should've returned false with slices of length %d", i)
		}
	}

	slice := make([]byte, 64)
	result := verifyToken(slice[:32], slice)
	if result != true {
		t.Errorf("VerifyToken should've returned true on a zeroed slice of length 64")
	}
}

func TestVerifiesMaskedTokenCorrectly(t *testing.T) {
	realToken := []byte("qwertyuiopasdfghjklzxcvbnm123456")
	sentToken := []byte("qwertyuiopasdfghjklzxcvbnm123456" +
		"\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00" +
		"\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00")

	if !verifyToken(realToken, sentToken) {
		t.Errorf("VerifyToken returned a false negative")
	}

	realToken[0] = 'x'

	if verifyToken(realToken, sentToken) {
		t.Errorf("VerifyToken returned a false positive")
	}
}

func TestOtpPanicsOnLengthMismatch(t *testing.T) {
	data := make([]byte, 1)
	key := make([]byte, 2)

	defer func() {
		if r := recover(); r == nil {
			t.Error("One time pad should've panicked on receiving slices" +
				"of different length, but it didn't")
		}
	}()
	oneTimePad(data, key)
}
func TestOtpMasksCorrectly(t *testing.T) {
	data := []byte("Inventors of the shish-kebab")
	key := []byte("They stop Cthulhu eating ye.")
	// precalculated
	expected := []byte("\x1d\x06\x13\x1cN\x07\x1b\x1d\x03\x00,\x12H\x01\x04" +
		"\rUS\r\x08\x07\x01C\x0cE\x1b\x04L")

	oneTimePad(data, key)

	if !bytes.Equal(data, expected) {
		t.Errorf("oneTimePad masked the data incorrectly: expected %#v, got %#v",
			expected, data)
	}
}

func TestOtpUnmasksCorrectly(t *testing.T) {
	orig := []byte("a very secret message")
	data := make([]byte, len(orig))
	copy(data, orig)
	if !bytes.Equal(orig, data) {
		t.Fatal("copy failed")
	}

	key := []byte("even more secret key!")

	oneTimePad(data, key)
	oneTimePad(data, key)

	if !bytes.Equal(orig, data) {
		t.Errorf("2x oneTimePad didn't return the original data:"+
			" expected %#v, got %#v", orig, data)
	}
}

func TestMasksTokenCorrectly(t *testing.T) {
	// needs to be of tokenLength
	token := []byte("12345678901234567890123456789012")
	fullToken := maskToken(token)

	if len(fullToken) != 2*tokenLength {
		t.Errorf("len(fullToken) is not %d, but %d", 2*tokenLength, len(fullToken))
	}

	key := fullToken[:tokenLength]
	encToken := fullToken[tokenLength:]

	// perform unmasking
	oneTimePad(encToken, key)

	if !bytes.Equal(encToken, token) {
		t.Errorf("Unmasked token is invalid: expected %v, got %v", token, encToken)
	}
}

func TestUnmasksTokenCorrectly(t *testing.T) {
	token := []byte("12345678901234567890123456789012")
	fullToken := maskToken(token)

	decToken := unmaskToken(fullToken)

	if !bytes.Equal(decToken, token) {
		t.Errorf("Unmasked token is invalid: expected %v, got %v", token, decToken)
	}
}
