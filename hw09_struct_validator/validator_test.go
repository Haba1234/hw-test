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
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
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
			in: Response{
				Code: 100,
				Body: "test string",
			},
			expectedErr: ValidationErrors{ValidationError{
				Field: "Code",
				Err:   ErrValidate,
			}},
		},
		{
			in: Response{
				Code: 200,
				Body: "test string",
			},
			expectedErr: ValidationErrors{},
		},
		{
			in: User{
				ID:     "3452452345234",
				Name:   "Ivan",
				Age:    16,
				Email:  "test.ru",
				Role:   "admi",
				Phones: []string{"01234567891", "12345"},
				meta:   nil,
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "ID",
					Err:   ErrValidate,
				},
				ValidationError{
					Field: "Age",
					Err:   ErrValidate,
				},
				ValidationError{
					Field: "Email",
					Err:   ErrValidate,
				},
				ValidationError{
					Field: "Role",
					Err:   ErrValidate,
				},
				ValidationError{
					Field: "Phones",
					Err:   ErrValidate,
				},
			},
		},
		{
			in: User{
				ID:     "2345609456847rth456ifhr560rfhrtfn563",
				Name:   "Ivan",
				Age:    20,
				Email:  "ivan@test.ru",
				Role:   "stuff",
				Phones: []string{"01234567891"},
				meta:   nil,
			},
			expectedErr: ValidationErrors{},
		},
		{
			in: App{Version: "0.01"},
			expectedErr: ValidationErrors{ValidationError{
				Field: "Version",
				Err:   ErrValidate,
			}},
		},
		{
			in:          App{Version: "10.01"},
			expectedErr: ValidationErrors{},
		},
		{
			in: Token{
				Header:    []byte{0, 1, 5},
				Payload:   nil,
				Signature: nil,
			},
			expectedErr: ValidationErrors{},
		},
		{
			in:          "empty",
			expectedErr: ValidationErrors{},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			require.Equal(t, tt.expectedErr, err, "test failed")
			_ = tt
		})
	}
}
