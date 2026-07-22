package errors

import "errors"

var (
	ErrNotFound           = errors.New("resource not found")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
	ErrBadRequest         = errors.New("bad request")
	ErrConflict           = errors.New("resource already exists")
	ErrInternal           = errors.New("internal server error")
	ErrInvalidToken       = errors.New("invalid token")
	ErrUserNotFound       = errors.New("user not found")
	ErrTeamNotFound       = errors.New("team not found")
	ErrTaskNotFound       = errors.New("task not found")
	ErrNotMember          = errors.New("user is not a member of this team")
	ErrNotAdmin           = errors.New("user is not an admin")
	ErrNotOwner           = errors.New("user is not the owner")
	ErrAssigneeNotInTeam  = errors.New("assignee is not a member of the team")
	ErrInvalidStatus      = errors.New("invalid status")
	ErrInvalidRole        = errors.New("invalid role")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

