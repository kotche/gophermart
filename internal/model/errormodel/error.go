package errormodel

import "fmt"

type ConflictLoginError struct {
	Login string
}

func (conflict ConflictLoginError) Error() string {
	return fmt.Sprintf("login %v already exists", conflict.Login)
}

type AuthenticationError struct{}

func (auth AuthenticationError) Error() string {
	return "invalid login/password"
}

type OrderAlreadyUploadedCurrentUserError struct{}

func (o OrderAlreadyUploadedCurrentUserError) Error() string {
	return "the order has already been uploaded by the current user"
}

type OrderAlreadyUploadedAnotherUserError struct{}

func (o OrderAlreadyUploadedAnotherUserError) Error() string {
	return "the order has already been uploaded by another user"
}

type CheckLuhnError struct{}

func (c CheckLuhnError) Error() string {
	return "invalid order number format"
}
