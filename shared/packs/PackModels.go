package packs

type Pack struct {
	Settings PackSettings `json:"settings"`
	Prompts  []Prompt     `json:"prompts"`
}

type PackSettings struct {
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Disabled    bool   `json:"disabled"`
}

type Prompt struct {
	UUID    string         `json:"uuid"`
	Details map[string]any `json:"details"`
}
