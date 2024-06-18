package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type ItemDataType interface {
	Raw() ([]byte, error)
}

// Card - bank card data
type Card struct {
	Number     string    `json:"number"`
	Owner      string    `json:"owner"`
	ValidityTo time.Time `json:"validity_to"`
}

// Raw - encode Card to byte slice
func (c *Card) Raw() ([]byte, error) {
	b, err := json.Marshal(c)
	if err != nil {
		return nil, fmt.Errorf("marshaling card data error: %w", err)
	}
	return b, nil
}

// Text - is a Text data type for encoding text
type Text struct {
	Data string `json:"data"`
}

// Raw - encode Text to byte slice
func (t *Text) Raw() ([]byte, error) {
	b, err := json.Marshal(t)
	if err != nil {
		return nil, fmt.Errorf("marshaling text data error: %w", err)
	}
	return b, nil
}

// Credential - is a Credential data type for encoding login and password
type Credential struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// Raw - encode Credential to byte slice
func (c *Credential) Raw() ([]byte, error) {
	b, err := json.Marshal(c)
	if err != nil {
		return nil, fmt.Errorf("marshaling credential data error: %w", err)
	}
	return b, nil
}

// File - file data
type File struct {
	Name      string `json:"name"`
	Extension string `json:"ext"`
	Data      []byte `json:"data"`
}

// Raw - encode File to byte slice
func (f *File) Raw() ([]byte, error) {
	b, err := json.Marshal(f)
	if err != nil {
		return nil, fmt.Errorf("marshaling file data error: %w", err)
	}
	return b, nil
}
