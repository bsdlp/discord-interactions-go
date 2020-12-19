package interactions

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

const (
	publicKeySize = 32
	signatureSize = 64
)

func verify(sig []byte, key ed25519.PublicKey) bool {
	if l := len(key); l != publicKeySize {
		panic("sign: bad public key length: " + strconv.Itoa(l))
	}

	if len(sig) < signatureSize || sig[63]&224 != 0 {
		return false
	}
	msg := sig[signatureSize:]
	sig = sig[:signatureSize]

	return ed25519.Verify(ed25519.PublicKey(key), msg, sig)
}

// Verify implements the verification side of the discord interactions api
// signing algorithm, as documented here:
// https://discord.com/developers/docs/interactions/slash-commands#security-and-authorization
func Verify(ctx context.Context, r *http.Request, key ed25519.PublicKey) bool {
	var payloadBuffer bytes.Buffer

	signature := r.Header.Get("X-Signature-Ed25519")
	if signature == "" {
		return false
	}

	sig, err := hex.DecodeString(signature)
	if err != nil {
		return false
	}
	payloadBuffer.Write(sig)

	timestamp := r.Header.Get("X-Signature-Timestamp")
	if timestamp == "" {
		return false
	}

	payloadBuffer.WriteString(timestamp)

	defer r.Body.Close()
	var body bytes.Buffer

	// at the end of the function, copy the original body back into the request
	defer func() {
		r.Body = ioutil.NopCloser(&body)
	}()

	// copy body into buffers
	_, err = io.Copy(&payloadBuffer, io.TeeReader(r.Body, &body))
	if err != nil {
		return false
	}

	return verify(payloadBuffer.Bytes(), key)
}
