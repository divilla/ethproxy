package json_error

import "github.com/tidwall/sjson"

type (
	MapError map[string]interface{}

	MessageError struct {
		message string
	}
)

func NewMessage(message string) *MessageError {
	return &MessageError{
		message: message,
	}
}

func (e *MessageError) String() string {
	json := `{}`
	json, _ = sjson.Set(json, "message", e.message)

	return json
}

func (e *MessageError) Error() string {
	return e.String()
}