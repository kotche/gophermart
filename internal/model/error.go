package model

import "fmt"

type ConflictLoginError struct {
	Err   error
	Login string
}

func (c ConflictLoginError) Error() string {
	return fmt.Sprintf("login %v already exists", c.Login)
}
