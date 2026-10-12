package todo

import (
	"strings"
	"unicode/utf8"
)

const MaxTitleLength = 200

func ValidateTitle(title string) (string, error) {
	title = strings.TrimSpace(title)

	if title == "" {
		return "", ErrEmptyTitle
	}
	if utf8.RuneCountInString(title) > MaxTitleLength {
		return "", ErrTitleTooLong
	}
	return title, nil
}
