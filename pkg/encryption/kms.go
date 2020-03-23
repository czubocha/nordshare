package encryption

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/kms"
	"nordshare/pkg/note"
)

type service struct {
	keyID  string
	client *kms.KMS
}

func NewService(keyID string, client *kms.KMS) *service {
	return &service{keyID: keyID, client: client}
}

func (s service) Encrypt(content *[]byte) error {
	if err := s.encrypt(content, s.keyID); err != nil {
		return fmt.Errorf("encrypt: %w", err)
	}
	return nil
}

func (s service) DecryptContent(note *note.Note) error {
	if err := s.decrypt(&note.Content); err != nil {
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
