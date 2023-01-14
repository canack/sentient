package models

import "strings"

type ResponseMessage string

func (r ResponseMessage) String() string {
	return string(r)
}

func (r ResponseMessage) Pretty() string {
	sanitized := strings.Replace(r.String(), "A: ", "", 1)
	return sanitized
}
