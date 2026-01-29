package sys

import (
	"errors"

	"github.com/AiratS/micro_as_bigtech_course/platform_common/pkg/sys/codes"
)

type commonError struct {
	msg  string
	code codes.Code
}

func NewCommonError(msg string, code codes.Code) *commonError {
	return &commonError{msg, code}
}

func (c *commonError) Error() string {
	return c.msg
}

func (c *commonError) Code() codes.Code {
	return c.code
}

func IsCommonError(err error) bool {
	var ce *commonError
	return errors.As(err, &ce)
}

func GetCommonError(err error) *commonError {
	var ce *commonError
	if !errors.As(err, &ce) {
		return nil
	}

	return ce
}
