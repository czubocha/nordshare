package hash

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"nordshare/pkg/note"
)

func HashNote(note *note.Note) error {
	if len(note.ReadPassword) > 0 {
		if err := hash(&note.ReadPassword); err != nil {
			return fmt.Errorf("hash: %w", err)
		}
	}
	if len(note.WritePassword) > 0 {
		if err := hash(&note.WritePassword); err != nil {
			return fmt.Errorf("hash: %w", err)
		}
	}
	return nil
}

func hash(content *[]byte) error {
	password, err := bcrypt.GenerateFromPassword(*content, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	*content = password
	return nil
}
