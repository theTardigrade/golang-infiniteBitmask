package infiniteBitmask

import "errors"

var (
	ErrValuesMismatched  = errors.New("values do not belong to the same generator")
	ErrLoadString        = errors.New("cannot load string")
	ErrPairUninitialized = errors.New("pair has not been initialized")
)
