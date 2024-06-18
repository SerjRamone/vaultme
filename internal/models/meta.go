package models

// Item metadata (key-value pairs)
// some described item data
// f.e. site or bank name
type Meta struct {
	Tag  string `json:"tag"`
	Text string `json:"text"`
}
