package validation

import (
	"errors"
	"net/url"
	"regexp"
)

var usernameRe = regexp.MustCompile(`^[a-zA-Z0-9_]{3,32}$`)

func ValidateUsername(u string) error {
	if !usernameRe.MatchString(u) {
		return errors.New("username must be 3–32 chars, alphanumeric or underscore")
	}
	return nil
}

func ValidatePassword(p string) error {
	if len(p) < 6 {
		return errors.New("password must be at least 6 characters")
	}
	return nil
}

func ValidateAdTitle(t string) error {
	if len(t) == 0 || len(t) > 100 {
		return errors.New("title must be 1–100 characters")
	}
	return nil
}

func ValidateAdText(txt string) error {
	if len(txt) == 0 || len(txt) > 2000 {
		return errors.New("text must be 1–2000 characters")
	}
	return nil
}

func ValidatePrice(price float64) error {
	if price < 0 || price > 1_000_000 {
		return errors.New("price must be between 0 and 1,000,000")
	}
	return nil
}

func ValidateImageURL(u string) error {
	parsed, err := url.ParseRequestURI(u)
	if err != nil {
		return errors.New("invalid image URL")
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return errors.New("image URL must be http or https")
	}
	return nil
}
