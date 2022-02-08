package infiniteBitmask

import "errors"

var (
	ErrValuesMismatched = errors.New("values do not belong to the same generator")
	ErrParseString      = errors.New("cannot parse string")
)
