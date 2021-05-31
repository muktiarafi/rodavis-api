package api

import (
	"encoding/json"
	"io"
)

func Bind(r io.Reader, v interface{}) error {
	err := json.NewDecoder(r).Decode(v)
	if err != nil {
		exc := NewSingleMessageException(
			EINVALID,
			"",
			"Invalid Payload",
			err,
		)
		return exc
	}
	return nil
}
