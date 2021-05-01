package cmd

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"
)

// Generates Gaois/Counter Mode using aes256 and hashed user passphrase
func getGCM(passPhrase string) (cipher.AEAD, error) {
	// create new aes cipher block
	cphr, err := aes.NewCipher(generateHash(passPhrase))
	if err != nil {
		return nil, err
	}

	// create new Gaois/Counter Mode with our aes cipher block
	gcm, err := cipher.NewGCM(cphr)
	if err != nil {
		return nil, err
	}

	return gcm, nil
}

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
	gcm, err := getGCM(passPhrase)
	if err != nil {
		return []byte{}, nil
	}

	// generate nonce by writing random bits to nonce-sized length buffer
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return []byte{}, err
	}

	// return encrypted byte array. Nonce is stored at the beginning of the text for use in decryption
	return gcm.Seal(nonce, nonce, journalBytes, nil), nil
}

/**
* Decrypt journal entry
* TODO: Refactor for use beyond testing
 */
func DecryptJournal(encryptedJournal []byte, passPhrase string) ([]byte, error) {
	gcm, err := getGCM(passPhrase)
	if err != nil {
		return []byte{}, nil
	}

	nonceSize := gcm.NonceSize()
	if len(encryptedJournal) < nonceSize {
		return []byte{}, errors.New("byte array read from file is invalid")
	}

	nonce := encryptedJournal[:nonceSize]
	encryptedJournalText := encryptedJournal[nonceSize:]
	return gcm.Open(nil, nonce, encryptedJournalText, nil)
}
