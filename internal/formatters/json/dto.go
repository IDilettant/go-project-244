package json

type changeDTO struct {
	Key      string      `json:"key"`
	Type     string      `json:"type"`
	OldValue any         `json:"oldValue,omitempty"`
	NewValue any         `json:"newValue,omitempty"`
	Children []changeDTO `json:"children,omitempty"`
}
