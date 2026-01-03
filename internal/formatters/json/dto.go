package json

type changeDTO struct {
	Type     string               `json:"type"`
	OldValue any                  `json:"oldValue,omitempty"`
	NewValue any                  `json:"newValue,omitempty"`
	Children map[string]changeDTO `json:"children,omitempty"`
}
