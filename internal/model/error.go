package model

import "fmt"

type ConflictLoginError struct {
	Err   error
	Login string
}

func (conflict ConflictLoginError) Error() string {
	return fmt.Sprintf("login %v already exists", conflict.Login)
}

type AuthenticationError struct {
	Err error
}

func (auth AuthenticationError) Error() string {
	return auth.Err.Error()
}
