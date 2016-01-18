package csrf

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"io"
)

// functions derived from nosurf https://github.com/justinas/nosurf

const (
	tokenLength = 32
)

func init() {
	checkForPRNG()
}

func checkForPRNG() {
	// Check that cryptographically secure PRNG is available
	// In case it's not, panic.
	buf := make([]byte, 1)
	_, err := io.ReadFull(rand.Reader, buf)

	if err != nil {
		panic(fmt.Sprintf("crypto/rand is unavailable: Read() failed with %#v", err))
	}
}

/*
There are two types of tokens.

* The unmasked "real" token consists of 32 random bytes.
  It is stored in a cookie (base64-encoded) and it's the
  "reference" value that sent tokens get compared to.

* The masked "sent" token consists of 64 bytes:
  32 byte key used for one-time pad masking and
  32 byte "real" token masked with the said key.
  It is used as a value (base64-encoded as well)
  in forms and/or headers.

Upon processing, both tokens are base64-decoded
and then treated as 32/64 byte slices.
*/

// A token is generated by returning tokenLength bytes
// from crypto/rand
func generateToken() []byte {
	bytes := make([]byte, tokenLength)

	if _, err := io.ReadFull(rand.Reader, bytes); err != nil {
		panic(err)
	}

	return bytes
}

func b64encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func b64decode(data string) []byte {
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil
	}
	return decoded
}

func verifyToken(realToken, sentToken []byte) bool {
	realN := len(realToken)
	sentN := len(sentToken)

	// sentN == tokenLength means the token is unmasked
	// sentN == 2*tokenLength means the token is masked.

	if realN == tokenLength && sentN == 2*tokenLength {
		return verifyMasked(realToken, sentToken)
	} else {
		return false
	}
}

// Verifies the masked token
func verifyMasked(realToken, sentToken []byte) bool {
	sentPlain := unmaskToken(sentToken)
	return subtle.ConstantTimeCompare(realToken, sentPlain) == 1
}

// Masks/unmasks the given data *in place*
// with the given key
// Slices must be of the same length, or oneTimePad will panic
func oneTimePad(data, key []byte) {
	n := len(data)
	if n != len(key) {
		panic("Lengths of slices are not equal")
	}

	for i := 0; i < n; i++ {
		data[i] ^= key[i]
	}
}

func maskToken(data []byte) []byte {
	if len(data) != tokenLength {
		return nil
	}

	// tokenLength*2 == len(enckey + token)
	result := make([]byte, 2*tokenLength)
	// the first half of the result is the OTP
	// the second half is the masked token itself
	key := result[:tokenLength]
	token := result[tokenLength:]
	copy(token, data)

	// generate the random token
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		panic(err)
	}

	oneTimePad(token, key)
	return result
}

func unmaskToken(data []byte) []byte {
	if len(data) != tokenLength*2 {
		return nil
	}

	key := data[:tokenLength]
	token := data[tokenLength:]
	oneTimePad(token, key)

	return token
}