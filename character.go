package main

type CharacterJSON struct {
}

func NewCharacterJSON(filepath string) (*CharacterJSON, error) {
	return &CharacterJSON{}, nil
}
