package qbus

import (
	"encoding/json"
	"errors"
	"fmt"
)

const (
	COMMAND_ERROR               = 0
	COMMAND_LOGIN               = 1
	COMMAND_LOGIN_RESPONSE      = 2
	COMMAND_GET_GROUPS          = 10
	COMMAND_GET_GROUPS_RESPONSE = 11
	COMMAND_SET_STATUS          = 12
	COMMAND_SET_STATUS_RESPONSE = 13
)

var (
	ErrBusy                   = errors.New("The controller is busy, please try again later")
	ErrSessionTimeout         = errors.New("The session timed out")
	ErrTooManyConnections     = errors.New("Too much devices are connected to the controller")
	ErrFailed                 = errors.New("The controller was unable to execute your command")
	ErrSessionStartFailure    = errors.New("Your session could not be started")
	ErrUnknownCommand         = errors.New("The command is unknown")
	ErrNoEQWebConfig          = errors.New("No EQOweb configuration found, please run System manager to upload and configure EQOweb.")
	ErrSystemManagerConnected = errors.New("System manager is still connected. Please close System manager to continue")
	ErrUndefined              = errors.New("Undefined error in the controller. Please try again later")
)

type Command struct {
	Type  int
	Value any
}

func NewLoginCommand(user, password string) *Command {
	return &Command{
		Type: COMMAND_LOGIN,
		Value: struct {
			Usr string
			Psw string
		}{
			Usr: user,
			Psw: password},
	}
}

func NewGetGroupsCommand() *Command {
	return &Command{Type: COMMAND_GET_GROUPS}
}

func NewSetStatusCommand(channelID int, status []int) *Command {
	return &Command{
		Type: COMMAND_SET_STATUS,
		Value: struct {
			Chnl int
			Val  []int
		}{
			Chnl: channelID,
			Val:  status},
	}
}

type Response struct {
	Type  int
	Value json.RawMessage
}

// Error returns the appropriate error if the response
// is an error response (Type==0), nil if not
func (r *Response) Error() (err error) {
	if r.Type != COMMAND_ERROR {
		return
	}
	errorResponse := &struct {
		Error int
	}{}
	if err = json.Unmarshal(r.Value, errorResponse); err != nil {
		return
	}
	switch errorResponse.Error {
	case 1:
		err = ErrBusy
	case 2:
		err = ErrSessionTimeout
	case 3:
		err = ErrTooManyConnections
	case 4:
		err = ErrFailed
	case 5:
		err = ErrSessionStartFailure
	case 6:
		err = ErrUnknownCommand
	case 7:
		err = ErrNoEQWebConfig
	case 8:
		err = ErrSystemManagerConnected
	case 255:
		err = ErrUndefined
	default:
		err = fmt.Errorf("The controller returned unknown errorcode %d", errorResponse.Error)
	}

	return
}

func (r *Response) GetLoginResponse() (sessionID string, err error) {
	if err = r.Error(); err != nil {
		return
	}
	if r.Type != COMMAND_LOGIN_RESPONSE {
		err = fmt.Errorf("Response Type %d is not a Login Response", r.Type)
	}
	responseValue := &struct {
		Rsp bool
		Id  string
	}{}
	if err = json.Unmarshal(r.Value, responseValue); err != nil {
		return
	}
	sessionID = responseValue.Id
	return
}

func (r *Response) GetGroupsResponse() (groups []Group, err error) {
	if err = r.Error(); err != nil {
		return
	}
	if r.Type != COMMAND_GET_GROUPS_RESPONSE {
		err = fmt.Errorf("Response Type %d is not a GetGroups Response", r.Type)
	}
	responseValue := &struct {
		Groups []Group
	}{}
	if err = json.Unmarshal(r.Value, responseValue); err != nil {
		return
	}
	groups = responseValue.Groups
	return
}

func (r *Response) SetStatusResponse() (channelID int, status []int, err error) {
	if err = r.Error(); err != nil {
		return
	}
	if r.Type != COMMAND_SET_STATUS_RESPONSE {
		err = fmt.Errorf("Response Type %d is not a SetStatus Response", r.Type)
	}
	responseValue := &struct {
		Chnl int
		Val  []int
	}{}
	if err = json.Unmarshal(r.Value, responseValue); err != nil {
		return
	}
	channelID = responseValue.Chnl
	status = responseValue.Val
	return
}
