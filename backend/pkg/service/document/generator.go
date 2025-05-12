package document

import (
	"context"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/ONLYOFFICE/onlyoffice-miro/backend/pkg/common"
)

type KeyGenerator interface {
	Generate(context.Context, DocumentConfigurer) (string, error)
}

type modificationKeyGenerator struct{}

func NewModificationKeyGenerator() KeyGenerator {
	return &modificationKeyGenerator{}
}

func (g *modificationKeyGenerator) Generate(ctx context.Context, configurer DocumentConfigurer) (string, error) {
	hasher := md5.New()
	hasher.Write([]byte(common.Concat(configurer.ID(), configurer.ModifiedAt())))
	hash := hasher.Sum(nil)

	return base64.URLEncoding.EncodeToString(hash), nil
}

type SignatureGenerator interface {
	Sign(key []byte, payload []byte) (string, error)
}

type jwtSignatureGenerator struct{}

func NewJwtSignatureGenerator() SignatureGenerator {
	return &jwtSignatureGenerator{}
}

func (g *jwtSignatureGenerator) Sign(key []byte, payload []byte) (string, error) {
	header := `{"alg":"HS256","typ":"JWT"}`
	hencoded := base64.RawURLEncoding.EncodeToString([]byte(header))

	pencoded, err := encodeClaims(payload)
	if err != nil {
		return "", err
	}

	token := common.Concat(hencoded, ".", pencoded)
	signature := computeHMAC(token, key)

	return common.Concat(token, ".", signature), nil
}

func encodeClaims(payload []byte) (string, error) {
	var claims map[string]any
	if err := json.Unmarshal(payload, &claims); err != nil {
		claims = make(map[string]any)
	}

	now := time.Now()
	claims["iat"] = now.Unix()
	claims["exp"] = now.Add(5 * time.Minute).Unix()

	npayload, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(npayload), nil
}

func computeHMAC(message string, key []byte) string {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(message))

	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
