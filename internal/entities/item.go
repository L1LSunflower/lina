package entities

import "github.com/google/uuid"

const (
	ReadyStatus = "ready"
	DoneStatus  = "done"
)

type Item struct {
	Id            string   `json:"id,omitempty"`
	Url           string   `json:"url"`
	Name          string   `json:"name"`
	Article       string   `json:"article"`
	ExpectedPrice int      `json:"expected_price"`
	ActualPrice   int      `json:"actual_price"`
	Currency      string   `json:"currency"`
	Colors        []string `json:"colors,omitempty"`
	Sizes         []string `json:"sizes,omitempty"`
	ImageLinks    []string `json:"image_links,omitempty"`
	Hash          string   `json:"hash,omitempty"`
	Status        string   `json:"status,omitempty"`
}

func (i *Item) SanitizedItem() Item {
	itemCopy := *i
	itemCopy.Id = ""
	itemCopy.Hash = ""
	itemCopy.Status = ""
	itemCopy.Colors = nil
	itemCopy.Sizes = nil
	itemCopy.ImageLinks = nil
	return itemCopy
}

func (i *Item) PrepareToSave(hash string) {
	uid, err := uuid.NewV7()
	if err != nil {
		return
	}
	i.Id = uid.String()
	i.Hash = hash
	i.Status = ReadyStatus
}
