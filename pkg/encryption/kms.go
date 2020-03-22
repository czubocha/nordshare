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

func (s service) EncryptContent(note *note.Note) error {
	if err := s.encrypt(&note.Content, s.keyID); err != nil {
		return fmt.Errorf("encrypt: %w", err)
	}
	return nil
}

func (s service) encrypt(content *[]byte, keyID string) error {
	input := kms.EncryptInput{
		KeyId:     &keyID,
		Plaintext: *content,
	}
	output, err := s.client.Encrypt(&input)
	if err != nil {
		return err
	}
	*content = output.CiphertextBlob
	return nil
}
