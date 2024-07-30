package encryptoha

import (
	"encoding/base64"
	"golang.org/x/crypto/argon2"
)

type IHashEncryptor interface {
	Hash(string) (string, error)
	Compare(string, string) (bool, error)
}

type Argon2Config struct {
	Salt      []byte
	Time      uint32
	Memory    uint32
	Threads   uint8
	KeyLength uint32
}

type Argon2Encryptor struct {
	config Argon2Config
}

func NewArgon2Encryptor(config Argon2Config) *Argon2Encryptor {
	return &Argon2Encryptor{config: config}
}

// Hash generates a hash for the given input string
func (e *Argon2Encryptor) Hash(input string) (string, error) {
	hash := argon2.IDKey([]byte(input), e.config.Salt, e.config.Time, e.config.Memory, e.config.Threads, e.config.KeyLength)
	encodedHash := base64.RawStdEncoding.EncodeToString(hash)
	return encodedHash, nil
}

// Compare checks if the given hash matches the input string
func (e *Argon2Encryptor) Compare(hash, input string) (bool, error) {
	newHash, err := e.Hash(input)
	if err != nil {
		return false, err
	}
	return hash == newHash, nil
}
