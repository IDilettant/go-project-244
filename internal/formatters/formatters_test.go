package formatters

import (
	"testing"

	"github.com/stretchr/testify/require"

	"code/internal/formatters/common"
)

func TestSelectFormatter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		format  string
		wantErr bool
	}{
		{name: "ok/default selects stylish", format: "", wantErr: false},
		{name: "ok/stylish selects stylish", format: common.FormatStylish, wantErr: false},
		{name: "error/unknown format", format: "unknown", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			f, err := SelectFormatter(tt.format)

			if tt.wantErr {
				require.Error(t, err)
				require.ErrorIs(t, err, common.ErrUnknownFormat)
				require.Nil(t, f)

				return
			}

			require.NoError(t, err)
			require.NotNil(t, f)
		})
	}
}
