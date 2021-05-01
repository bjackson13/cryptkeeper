package cmd

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"io"
)

// generate sha256 hash from users passphrase
func generateHash(passPhrase string) []byte {
	h := sha256.New()
	h.Write([]byte(passPhrase))
	return h.Sum(nil)
}

/**
* Encrypts the journal contents. Uses the usrs passphrase to generate hash.
 */
func EncryptJournal(journalBytes []byte, passPhrase string) ([]byte, error) {
	// create new cipher Block
	cphr, err := aes.NewCipher(generateHash(passPhrase))
	if err != nil {
		return []byte{}, err
	}

	// create new Gaois/Counter Mode
	gcm, err := cipher.NewGCM(cphr)
	if err != nil {
		return []byte{}, err
	}

	// generate nonce by writing random bits to nonce-sized length buffer
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return []byte{}, err
	}

	// return encrypted byte array
	return gcm.Seal(nonce, nonce, journalBytes, nil), nil
}
