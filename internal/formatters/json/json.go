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
	payload := convertChangesToDTOMap(changes)

	encoded, err := json.MarshalIndent(payload, "", common.IndentBase)
	if err != nil {
		return "", err
	}

	output := string(encoded) + common.NewLine

	return output, nil
}

func convertChangesToDTOMap(changes []diff.Change) map[string]changeDTO {
	result := make(map[string]changeDTO, len(changes))
	for _, change := range changes {
		result[change.Key] = convertChangeToDTO(change)
	}

	return result
}

func convertChangeToDTO(change diff.Change) changeDTO {
	dto := changeDTO{Type: changeTypeToString(change)}

	switch change.Type {
	case diff.Nested:
		dto.Children = convertChangesToDTOMap(change.Children)

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
	switch change.Type {
	case diff.Nested:
		return changeTypeNested

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
