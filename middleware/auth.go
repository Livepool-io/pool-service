package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"os"

	"github.com/ethereum/go-ethereum/crypto"
)

func IsHMACAuthorized(authHeader string, body []byte) bool {
	hash := hmac.New(sha256.New, []byte(os.Getenv("SECRET")))
	hash.Write(body)
	return hex.EncodeToString(hash.Sum(nil)) == authHeader
}

// TODO: Use Sign-in with Ethereum https://docs.login.xyz/libraries/go
func IsECDSAAuthorized(address, signature string, data []byte) bool {
	hash := crypto.Keccak256Hash(data)

	sigPublicKeyECDSA, err := crypto.SigToPub(hash.Bytes(), []byte(signature))
	if err != nil {
		return false
	}
	addr := crypto.PubkeyToAddress(*sigPublicKeyECDSA)

	return addr.Hex() == address
}
