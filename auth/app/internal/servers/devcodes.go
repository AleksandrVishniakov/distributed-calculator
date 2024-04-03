package servers

import "strings"

type DeveloperCode int

const (
	CodeLoginIsRequired    DeveloperCode = 400_001
	CodePasswordIsRequired DeveloperCode = 400_002
	CodeTooLongPassword    DeveloperCode = 400_003
	CodeInvalidCredentials DeveloperCode = 400_004

	CodeUserNotFound DeveloperCode = 404_001

	CodeUserAlreadyExists DeveloperCode = 409_001
)

func DevCodeFromErr(err error) (DeveloperCode, bool) {
	if err == nil {
		return DeveloperCode(0), false
	}

	return DevCodeFromMsg(err.Error())
}

func DevCodeFromMsg(msg string) (DeveloperCode, bool) {
	switch {
	case strings.Contains(msg, MsgLoginIsRequired):
		return CodeLoginIsRequired, true
	case strings.Contains(msg, MsgPasswordIsRequired):
		return CodePasswordIsRequired, true
	case strings.Contains(msg, MsgTooLongPassword):
		return CodeTooLongPassword, true
	case strings.Contains(msg, MsgInvalidCredentials):
		return CodeInvalidCredentials, true

	case strings.Contains(msg, MsgUserNotFound):
		return CodeUserNotFound, true

	case strings.Contains(msg, MsgUserAlreadyExists):
		return CodeUserAlreadyExists, true

	default:
		return DeveloperCode(0), false
	}
}
