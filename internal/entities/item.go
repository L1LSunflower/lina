package entities

type Item struct {
	Id            string   `json:"id"`
	Url           string   `json:"url"`
	Name          string   `json:"name"`
	Article       string   `json:"article"`
	ExpectedPrice int      `json:"expected_price"`
	ActualPrice   int      `json:"actual_price"`
	Currency      string   `json:"currency"`
	Colors        []string `json:"colors"`
	Sizes         []string `json:"sizes"`
	ImageLinks    []string `json:"image_links"`
	Hash          string   `json:"hash,omitempty"`
	Status        string   `json:"status,omitempty"`
}
