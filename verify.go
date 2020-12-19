package interactions

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"io"
	"net/http"

	"github.com/kevinburke/nacl/sign"
)

func Verify(ctx context.Context, r *http.Request, key ed25519.PublicKey) bool {
	r = r.Clone(ctx)
	defer r.Body.Close()

	var buf bytes.Buffer

	timestamp := r.Header.Get("X-Signature-Timestamp")
	if timestamp == "" {
		return false
	}

	signature := r.Header.Get("X-Signature-Ed25519")
	if signature == "" {
		return false
	}

	sig, err := hex.DecodeString(signature)
	if err != nil {
		return false
	}

	buf.Write(sig)
	buf.WriteString(timestamp)

	_, err = io.Copy(&buf, r.Body)
	if err != nil {
		return false
	}

	return sign.PublicKey(key).Verify(buf.Bytes())
}
