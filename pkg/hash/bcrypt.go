package hash

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"nordshare/pkg/note"
)

func Passwords(note *note.Note) error {
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

func HasReadAccess(note note.Note, password []byte) bool {
	if len(note.ReadPassword) == 0 {
		return true
	}
	if err := bcrypt.CompareHashAndPassword(note.ReadPassword, password); err == nil {
		return true
	}
	if len(note.WritePassword) > 0 && bcrypt.CompareHashAndPassword(note.WritePassword, password) == nil {
		return true
	}
	return false
}

func hash(content *[]byte) error {
	password, err := bcrypt.GenerateFromPassword(*content, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	*content = password
	return nil
}
