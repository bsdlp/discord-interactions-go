package interactions_test

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/bsdlp/discord-interactions-go/interactions"
	"github.com/kevinburke/nacl/sign"
	"github.com/stretchr/testify/suite"
)

type InteractionsTestSuite struct {
	suite.Suite
	priv      ed25519.PrivateKey
	pub       ed25519.PublicKey
	timestamp string
}

func (suite *InteractionsTestSuite) SetupTest() {
	var err error
	suite.pub, suite.priv, err = ed25519.GenerateKey(nil)
	suite.Require().NoError(err)
	suite.timestamp = strconv.FormatInt(time.Now().Unix(), 10)
}

func (suite *InteractionsTestSuite) TestVerifySuccess() {
	body := "body"
	request := httptest.NewRequest("POST", "http://localhost/interaction", strings.NewReader(body))
	request.Header.Set("X-Signature-Timestamp", suite.timestamp)

	var msg bytes.Buffer
	msg.WriteString(suite.timestamp)
	msg.WriteString(body)
	signature := sign.Sign(msg.Bytes(), sign.PrivateKey(suite.priv))
	request.Header.Set("X-Signature-Ed25519", hex.EncodeToString(signature[:sign.SignatureSize]))

	suite.Assert().True(interactions.Verify(request, suite.pub))
}

func (suite *InteractionsTestSuite) TestVerifyModifiedBody() {
	body := "body"
	request := httptest.NewRequest("POST", "http://localhost/interaction", strings.NewReader("WRONG"))
	request.Header.Set("X-Signature-Timestamp", suite.timestamp)

	var msg bytes.Buffer
	msg.WriteString(suite.timestamp)
	msg.WriteString(body)
	signature := sign.Sign(msg.Bytes(), sign.PrivateKey(suite.priv))
	request.Header.Set("X-Signature-Ed25519", hex.EncodeToString(signature[:sign.SignatureSize]))

	suite.Assert().False(interactions.Verify(request, suite.pub))
}
func (suite *InteractionsTestSuite) TestVerifyModifiedTimestamp() {
	body := "body"
	request := httptest.NewRequest("POST", "http://localhost/interaction", strings.NewReader("WRONG"))
	request.Header.Set("X-Signature-Timestamp", strconv.FormatInt(time.Now().Add(time.Minute).Unix(), 10))

	var msg bytes.Buffer
	msg.WriteString(suite.timestamp)
	msg.WriteString(body)
	signature := sign.Sign(msg.Bytes(), sign.PrivateKey(suite.priv))
	request.Header.Set("X-Signature-Ed25519", hex.EncodeToString(signature[:sign.SignatureSize]))

	suite.Assert().False(interactions.Verify(request, suite.pub))
}
func TestInteractions(t *testing.T) {
	suite.Run(t, &InteractionsTestSuite{})
}
