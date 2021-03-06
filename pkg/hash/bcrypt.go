package hash

import (
	"fmt"
	"github/czubocha/nordshare"
	"golang.org/x/crypto/bcrypt"
)

func Passwords(note *nordshare.Note) error {
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

func HasReadAccess(note nordshare.Note, password []byte) bool {
	return notExistsOrMatches(note.ReadPassword, password) || existsAndMatches(note.WritePassword, password)
}

func HasWriteAccess(note nordshare.Note, password []byte) bool {
	return existsAndMatches(note.WritePassword, password)
}

func notExistsOrMatches(passwordHash, password []byte) bool {
	return len(passwordHash) == 0 || bcrypt.CompareHashAndPassword(passwordHash, password) == nil
}

func existsAndMatches(passwordHash, password []byte) bool {
	return len(passwordHash) > 0 && bcrypt.CompareHashAndPassword(passwordHash, password) == nil
}

func hash(content *[]byte) error {
	password, err := bcrypt.GenerateFromPassword(*content, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	*content = password
	return nil
}
