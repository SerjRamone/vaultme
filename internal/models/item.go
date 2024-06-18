package models

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// ItemStorage is a interface for working with items in database
type ItemStorage interface {
	// GetItem gets user's item by a given ID
	GetItem(ctx context.Context, userID string, itemID string) (*Item, error)
	// CreateItem creates new item
	CreateItem(ctx context.Context, userID string, item *ItemDTO) (*Item, error)
	// UpdateItem updates item
	UpdateItem(ctx context.Context, userID string, item *Item) (*Item, error)
	// ListItems return list of user's items with pagination
	ListItems(ctx context.Context, userID string, limit int, offset int) ([]*Item, error)
}

// ItemType item type (credential, text, raw, card)
// - CREDENTAIL - item contents credentals (login + password)
// - TEXT - item contents simple text data
// - RAW - item contents binary data (file bytes)
// - CARD - item contains card (bank plastic card) data: number, owner and validity date
type ItemType string

const (
	// Credentials item type
	CredentialType ItemType = "CREDENTIAL"
	// Text item type
	TextType ItemType = "TEXT"
	// Raw item type (binary data)
	RawType ItemType = "RAW"
	// Card item type
	CardType ItemType = "CARD"
)

// Item contents all user's item data
type Item struct {
	// ID - item UUID
	ID string `json:"id"`
	// UserID - user UUID
	UserID string `json:"user_id"`
	// Name - user-given item name or file name
	Name string `json:"name"`
	// Type - specifies type of item
	Type string `json:"type"`
	// Version - item version
	Version int64 `json:"version"`
	// CreatedAt - item creation date
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt - item update date
	UpdatedAt time.Time `json:"updated_at"`
	// Data - item type depend data. See `ItemType`
	Data []byte `json:"data"`
	// Meta - item metadata (key-value pairs)
	Meta []*Meta `json:"meta"`
}

// NewItem creates new item
func NewItem(id string,
	name string,
	itemType ItemType,
	version int64,
	createdAt time.Time,
	updatedAt time.Time,
	data ItemDataType,
	meta []*Meta,
) (*Item, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("marshaling data error: %w", err)
	}

	return &Item{
		ID:        id,
		Name:      name,
		Type:      string(itemType),
		Version:   version,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Data:      b,
		Meta:      meta,
	}, nil
}

// ItemDTO is a data transfer object for item
type ItemDTO struct {
	// Name - user-given item name
	// it's file name if item type is RAW
	Name string `json:"name"`
	// Type - specifies type of item
	Type string `json:"type"`
	// Data - item type depend data. See `ItemType`
	Data []byte `json:"data"`
	// Meta - key-value pairs to store additional data
	Meta []*Meta `json:"meta"`
}

// NewItemDTO return new ItemDTO object
func NewItemDTO(name string, itemType ItemType, data ItemDataType, meta []*Meta) (*ItemDTO, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("marshaling data error: %w", err)
	}

	return &ItemDTO{
		Name: name,
		Type: string(itemType),
		Data: b,
		Meta: meta,
	}, nil
}
