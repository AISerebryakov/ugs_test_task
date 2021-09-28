package http

import (
	"fmt"

	"github.com/francoispqt/gojay"
)

const (
	objectsLimitExceededTitle = "objects_limit_exceeded"

	requestPartiallyCompletedTitle = "request_partially_completed"
)

type Warning struct {
	title string
	msg   string
}

func NewLimitExceededWarning() Warning {
	return Warning{
		title: objectsLimitExceededTitle,
		msg:   fmt.Sprintf("limit: %d objects", maxGettingObjects),
	}
}

func (warn *Warning) MakeEmpty() {
	warn.title = ""
	warn.msg = ""
}

func (warn Warning) IsEmpty() bool {
	return len(warn.title) == 0
}

func (warn Warning) MarshalJSONObject(enc *gojay.Encoder) {
	enc.AddStringKey(TitleKey, warn.title)
	enc.AddStringKey(MsgKey, warn.msg)
}

func (warn Warning) IsNil() bool {
	return warn.IsEmpty()
}

func (warn Warning) String() string {
	return warn.title + " " + warn.msg
}
