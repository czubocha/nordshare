package note

import (
	"fmt"
	"strings"
)

func ExpirationScenario(h *Handler) error {
	// create note
	noteToCreate := CreateNote{
		Content:       content,
		ReadPassword:  readPassword,
		WritePassword: writePassword,
		TTL:           ttl,
	}
	id, err := create(h.Lambda, h.Saver, noteToCreate)
	if err != nil {
		return err
	}
	if err = expire(id, h.TableName, h.DynamoDB); err != nil {
		return err
	}
	err = get(h.Lambda, h.Reader, id, readPassword, ReadModifyNote{})
	if err == nil || strings.HasPrefix(err.Error(), "read note: expected status") {
		return fmt.Errorf("expiration scenario: expired note read")
	}
	return nil
}
