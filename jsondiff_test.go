package jsondiff

import (
	"testing"
)

func TestCompareJSON(t *testing.T) {
	tests := []struct {
		name        string
		a           []byte
		b           []byte
		opts        *CompareJSONOptions
		expectedRes string
		expectedErr error
	}{
		{
			name:        "empty json",
			a:           []byte(`{}`),
			b:           []byte(`{}`),
			opts:        nil,
			expectedRes: `{}`,
			expectedErr: nil,
		},
		{
			name:        "json with added key",
			a:           []byte(`{"foo": "bar"}`),
			b:           []byte(`{"foo": "bar", "baz": "qux"}`),
			opts:        nil,
			expectedRes: `{"added":{"baz":"qux"}}`,
			expectedErr: nil,
		},
		{
			name:        "json with removed key",
			a:           []byte(`{"foo": "bar", "baz": "qux"}`),
			b:           []byte(`{"foo": "bar"}`),
			opts:        nil,
			expectedRes: `{"removed":{"baz":"qux"}}`,
			expectedErr: nil,
		},
		{
			name:        "json with changed key",
			a:           []byte(`{"foo": "bar", "baz": "qux"}`),
			b:           []byte(`{"foo": "bar", "baz": "quux"}`),
			opts:        nil,
			expectedRes: `{"changed":{"baz":{"new":"quux","old":"qux"}}}`,
			expectedErr: nil,
		},
		{
			name:        "json with nested objects",
			a:           []byte(`{"foo": {"bar": "baz"}}`),
			b:           []byte(`{"foo": {"bar": "qux"}}`),
			opts:        nil,
			expectedRes: `{"changed":{"foo":{"bar":{"new":"qux","old":"baz"}}}}`,
			expectedErr: nil,
		},
		{
			name: "json with custom option keys",
			a:    []byte(`{"foo": "bar", "baz": "qux"}`),
			b:    []byte(`{"foo": "bar", "baz": "quux"}`),
			opts: &CompareJSONOptions{
				AddedKey:   "new",
				RemovedKey: "del",
				ChangedKey: "diff",
			},
			expectedRes: `{"diff":{"baz":{"new":"quux","old":"qux"}}}`,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := CompareJSON(tt.a, tt.b, tt.opts)
			if err != tt.expectedErr {
				t.Errorf("Expected error: %v, but got: %v", tt.expectedErr, err)
			}
			if res != tt.expectedRes {
				t.Errorf("Expected result: %v, but got: %v", tt.expectedRes, res)
			}
		})
	}
}
