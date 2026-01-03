package json

import (
	"encoding/json"

	"code/internal/diff"
	"code/internal/formatters/common"
)

const (
	changeTypeNested    = "nested"
	changeTypeUnchanged = "unchanged"
	changeTypeRemoved   = "removed"
	changeTypeAdded     = "added"
	changeTypeUpdated   = "updated"
	changeTypeUnknown   = "unknown"
)

type Formatter struct{}

func New() *Formatter { return &Formatter{} }

func (f *Formatter) Format(changes []diff.Change) (string, error) {
	payload := convertChangesToDTO(changes)

	encoded, err := json.MarshalIndent(payload, "", common.IndentBase)
	if err != nil {
		return "", err
	}

	return string(encoded) + common.NewLine, nil
}

func convertChangesToDTO(changes []diff.Change) []changeDTO {
	result := make([]changeDTO, 0, len(changes))
	for _, change := range changes {
		result = append(result, convertChangeToDTO(change))
	}
	return result
}

func convertChangeToDTO(change diff.Change) changeDTO {
	dto := changeDTO{
		Key:  change.Key,
		Type: changeTypeToString(change),
	}

	if change.IsContainer() {
		dto.Children = convertChangesToDTO(change.Children)

		return dto
	}

	switch change.Type {
	case diff.Unchanged, diff.Removed:
		dto.OldValue = change.OldValue

	case diff.Added:
		dto.NewValue = change.NewValue

	case diff.Updated:
		dto.OldValue = change.OldValue
		dto.NewValue = change.NewValue
	}

	return dto
}

func changeTypeToString(change diff.Change) string {
	if change.IsContainer() {
		return changeTypeNested
	}

	switch change.Type {
	case diff.Unchanged:
		return changeTypeUnchanged
	case diff.Removed:
		return changeTypeRemoved
	case diff.Added:
		return changeTypeAdded
	case diff.Updated:
		return changeTypeUpdated
	default:
		return changeTypeUnknown
	}
}

