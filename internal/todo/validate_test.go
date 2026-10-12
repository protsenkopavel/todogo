package todo

import (
	"errors"
	"strings"
	"testing"
)

func TestValidateTitle(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr error
	}{
		{"обычное название", "купить молоко", "купить молоко", nil},
		{"пробелы по краям обрезаются", "  молоко  ", "молоко", nil},
		{"пустая строка", "", "", ErrEmptyTitle},
		{"только пробелы", "     ", "", ErrEmptyTitle},
		{"200 кириллических символов", strings.Repeat("я", 200), strings.Repeat("я", 200), nil},
		{"201 символ", strings.Repeat("я", 201), "", ErrTitleTooLong},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateTitle(tt.input)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}
