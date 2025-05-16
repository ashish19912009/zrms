package jwk

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lestrrat-go/jwx/v2/jwk"
)

var jwkSet jwk.Set

// InitializeJWK initializes the JWK Set from a given RSA public key
func InitializeJWK(pubKey *rsa.PublicKey, keyID string) error {
	if pubKey == nil {
		return fmt.Errorf("nil public key provided")
	}

	key, err := jwk.FromRaw(pubKey)
	if err != nil {
		return err
	}

	if err := key.Set(jwk.KeyIDKey, keyID); err != nil {
		return err
	}

	jwkSet = jwk.NewSet()
	jwkSet.AddKey(key)

	return nil
}

// GetJWKSet returns the initialized JWK Set
func GetJWKSet() jwk.Set {
	return jwkSet
}

func Handler(w http.ResponseWriter, r *http.Request) {
	if jwkSet == nil {
		http.Error(w, "JWK set not initialized", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(jwkSet); err != nil {
		http.Error(w, "Failed to encode JWK set", http.StatusInternalServerError)
	}
}
