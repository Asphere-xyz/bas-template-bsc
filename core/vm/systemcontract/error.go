package systemcontract

import "fmt"

var (
	errNotSupported   = fmt.Errorf("not supported")
	errInvalidCaller  = fmt.Errorf("invalid caller")
	errFailedToUnpack = fmt.Errorf("failed to unpack")
)
