package note

func FullFlowScenario(h *Handler) error {
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
	// read note with read password
	expectedReadNote := ReadModifyNote{
		Content: content,
		TTL:     ttl,
	}
	if err = get(h.Lambda, h.Reader, id, readPassword, expectedReadNote); err != nil {
		return err
	}
	// read note with write password
	if err = get(h.Lambda, h.Reader, id, writePassword, expectedReadNote); err != nil {
		return err
	}
	// modify note
	modifiedNote := ReadModifyNote{
		Content: modifiedContent,
		TTL:     modifiedTtl,
	}
	if err = modify(h.Lambda, h.Modifier, id, writePassword, modifiedNote); err != nil {
		return err
	}
	// read note to check if note is modified
	if err = get(h.Lambda, h.Reader, id, readPassword, modifiedNote); err != nil {
		return err
	}
	// delete note
	if err = remove(h.Lambda, h.Remover, writePassword, id); err != nil {
		return err
	}
	return nil
}
