package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrap(t *testing.T) {
	tests := map[string]struct {
		err           error
		msg           string
		args          []any
		expectedError string
	}{
		"NilError_ReturnNil": {
			err: nil,
		},
		"WithMessageNoArgs": {
			err:           errors.New("failure"),
			msg:           "this operation failed",
			args:          nil,
			expectedError: "this operation failed: failure",
		},
		"WithMessageAndArgs": {
			err:           errors.New("failure"),
			msg:           "this operation failed with %s",
			args:          []any{"argument"},
			expectedError: "this operation failed with argument: failure",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := Wrap(tt.err, tt.msg, tt.args...)
			if tt.expectedError != "" {
				assert.EqualError(t, result, tt.expectedError)
			} else {
				assert.NoError(t, result)
			}
		})
	}
}
