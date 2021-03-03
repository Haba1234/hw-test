package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
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
			in: Response{Code: 100,
				Body: "rerferf",
			},
			expectedErr: ValidationErrors{ValidationError{
				Field: "Code",
				Err:   ErrFieldWrong,
			}},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			// Place your code here.
			err := Validate(tt.in)
			fmt.Printf("Test. Errors: %+v, %T\n", err, err)
			fmt.Printf("Test1. Errors: %+v, %T\n", tt.expectedErr, tt.expectedErr)
			if errors.Is(err, tt.expectedErr) {
				fmt.Println("Есть ошибка")
			}
			_ = tt
		})
	}
}
