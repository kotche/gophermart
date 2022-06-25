package model

import "fmt"

type ConflictLoginError struct {
	Err   error
	Login string
}

func (conflict ConflictLoginError) Error() string {
	return fmt.Sprintf("login %v already exists", conflict.Login)
}

type AuthorizationError struct {
	Err error
}

func (auth AuthorizationError) Error() string {
	return auth.Err.Error()
}
