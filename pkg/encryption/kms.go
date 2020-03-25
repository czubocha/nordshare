package encryption

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/kms"
)

type service struct {
	keyID  string
	client crypter
}
type crypter interface {
	Encrypt(input *kms.EncryptInput) (*kms.EncryptOutput, error)
	Decrypt(input *kms.DecryptInput) (*kms.DecryptOutput, error)
}

func NewService(keyID string, client crypter) *service {
	return &service{keyID: keyID, client: client}
}

func (s service) Encrypt(content *[]byte) error {
	if err := s.encrypt(content, s.keyID); err != nil {
		return fmt.Errorf("encrypt: %w", err)
	}
	return nil
}

func (s service) Decrypt(content *[]byte) error {
	if err := s.decrypt(content); err != nil {
		return fmt.Errorf("decrypt: %w", err)
	}
	return nil
}

func (s service) encrypt(plaintext *[]byte, keyID string) error {
	input := kms.EncryptInput{
		KeyId:     &keyID,
		Plaintext: *plaintext,
	}
	output, err := s.client.Encrypt(&input)
	if err != nil {
		return err
	}
	*plaintext = output.CiphertextBlob
	return nil
}

func (s service) decrypt(ciphertext *[]byte) error {
	input := kms.DecryptInput{
		CiphertextBlob: *ciphertext,
	}
	output, err := s.client.Decrypt(&input)
	if err != nil {
		return err
	}
	*ciphertext = output.Plaintext
	return nil
}
