package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				ID:     "11111111-1111-1111-1111-123456789abc",
				Name:   "Cat Dog",
				Age:    25,                      // Valid Age
				Email:  "v@example.com",         // Valid Email
				Role:   "admin",                 // Valid Role
				Phones: []string{"12345678901"}, // Valid Phone
			},
			expectedErr: nil,
		},
		{
			in: User{
				ID:     "invalid-id-not-36-chars",
				Name:   "Invalid User",
				Age:    17,
				Email:  "invalid-email",
				Role:   "nonexistent-role",
				Phones: []string{"123456789"},
			},
			expectedErr: ValidationErrors{
				{"ID", ErrLen},
				{"Age", ErrMin},
				{"Email", ErrRegexp},
				{"Role", ErrIn},
				{"Phones", ErrLen},
			},
		},
		{
			in: App{
				Version: "12345",
			},
			expectedErr: nil,
		},
		{
			in: App{
				Version: "v1.0",
			},
			expectedErr: ValidationErrors{
				{"Version", ErrLen},
			},
		},
		{
			in: Response{
				Code: 200,
				Body: "OK",
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 800,
				Body: "Error",
			},
			expectedErr: ValidationErrors{
				{"Code", ErrIn},
			},
		},
		{
			in: Token{
				Header:    []byte("header"),
				Payload:   []byte("payload"),
				Signature: []byte("signature"),
			},
			expectedErr: nil,
		},
		{
			in:          struct{}{},
			expectedErr: nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			require.Equal(t, tt.expectedErr, err)
			_ = tt
		})
	}
}
