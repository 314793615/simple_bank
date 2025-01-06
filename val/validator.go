package val

import (
	"errors"
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidUserName = regexp.MustCompile(`^[a-zA-z0-9_]+$`).MatchString
	
)

func ValidString(s string, minLen int, maxLen int) error {
	if len(s) > maxLen || len(s) < minLen {
		return fmt.Errorf("string length must between %d-%d", minLen, maxLen)
	}
	return nil
}

func ValidUserName(username string) error {
	err := ValidString(username, 3, 100)
	if err != nil {
		return err
	}
	if ok := isValidUserName(username); !ok {
		return errors.New("the username must contains lower case, bing case, or underline")
	}
	return nil
}

func ValidFullName(fullname string) error {
	err := ValidString(fullname, 3, 100)
	if err != nil {
		return err
	}
	return nil
}

func ValidPassword(password string) error {
	err := ValidString(password, 6, 10)
	if err != nil {
		return err
	}
	return nil
}

func ValidEmail(e string) error {
	err := ValidString(e, 3, 100)
	if err != nil {
		return err
	}

	if _, err = mail.ParseAddress(e); err != nil {
		return errors.New("the email address is not valid")
	}
	return nil
}