package dataloadgen

import "errors"

var ErrNotFound = errors.New("Record not found via dataloader")

type ErrorMap[KeyT comparable] map[KeyT]error

func (e ErrorMap[KeyT]) Error() string {
	str := ""
	for _, v := range e {
		str += v.Error()
	}
	return str
}
