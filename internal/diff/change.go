package diff

type ChangeType int

const (
	Unchanged ChangeType = iota
	Removed
	Added
	Updated
)

type Change struct {
	Key      string
	Type     ChangeType
	OldValue any
	NewValue any
}
