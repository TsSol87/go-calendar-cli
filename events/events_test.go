package events

import (
	"testing"
)

func TestIsValidTitle(t *testing.T) {
	text := "Проверка валидности текста"
	err := IsValidTitle(text)
	if err != nil {
		t.Errorf("Expected no error for valid title %q, but got: %v", text, err)
	}
}

func TestIsValidTitle_Short(t *testing.T) {
	text := "Hi"
	err := IsValidTitle(text)
	if err == nil {
		t.Errorf("Expected an error for string short than 3 characters, got none")
	}
}

func TestIsValidTitle_Long(t *testing.T) {
	text := "This is a very long string that exceeds the maximum length of fifty characters and it should fail"
	err := IsValidTitle(text)
	if err == nil {
		t.Errorf("Expected an error for string longer than 50 characters, got none")
	}
}

func TestIsValidTitle_InvalidCharacters(t *testing.T) {
	text := "Invalid character!"
	err := IsValidTitle(text)
	if err == nil {
		t.Errorf("Expected an error for string with invalid characters, got none")
	}
}

func TestIsValidTitle_MixedCharacters(t *testing.T) {
	text := "Hello мир 123"
	err := IsValidTitle(text)
	if err != nil {
		t.Errorf("Expected no error for mixed characters, got an error")
	}
}

func TestIsValidTitle_OnlyNumbers(t *testing.T) {
	text := "1234567890"
	err := IsValidTitle(text)
	if err != nil {
		t.Errorf("Expected no error for string containing only numbers, got an error")
	}
}

func TestIsValidTitle_OnlySpaces(t *testing.T) {
	text := "   "
	err := IsValidTitle(text)
	if err != nil {
		t.Errorf("Expected no error for string containing only spaces, got an error")
	}
}

func TestIsValidTitle_EmptyString(t *testing.T) {
	text := ""
	err := IsValidTitle(text)
	if err == nil {
		t.Errorf("Expected error for empty string, got no error")
	}
}

func TestTimeParse_EmptyString(t *testing.T) {
	input := ""

	_, actualError := TimeParse(input)

	if actualError == nil {
		t.Errorf("Expected an error for empty string, but got none")
	}
}

func TestTimeParse_NonExistentDate(t *testing.T) {
	input := "2024-02-30 10:30"

	_, actualError := TimeParse(input)

	if actualError == nil {
		t.Errorf("Expected an error for non-existent date, but got none")
	}
}

func TestTimeParse_InvalidFormat(t *testing.T) {
	input := "27-10-2024 10:30"

	_, actualError := TimeParse(input)

	if actualError == nil {
		t.Errorf("Expected an error for invalid date format, but got none")
	}
}

func TestTimeParse_ValidFormat(t *testing.T) {
	input := "2024-10-27 10:30"
	_, actualError := TimeParse(input)
	if actualError != nil {
		t.Errorf("Expected no error for valid date format, but got error")
	}
}
