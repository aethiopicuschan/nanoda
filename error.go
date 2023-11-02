package nanoda

import "fmt"

type Error struct {
	Code ResultCode
	Msg  string
}

func (e Error) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Msg)
}

func (v *Voicevox) newError(code ResultCode) error {
	return Error{
		Code: code,
		Msg:  v.voicevoxErrorResultToMessage(code),
	}
}
