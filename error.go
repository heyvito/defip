package defip

import (
	"fmt"
	"runtime"
)

// ErrCantParse is returned if the route table is garbage.
type ErrCantParse struct{}

// ErrNotImplemented is returned if your operating system
// is not supported by this package. Please raise an issue
// to request support.
type ErrNotImplemented struct{}

// ErrInvalidRouteFileFormat is returned if the format
// of /proc/net/route is unexpected on Linux systems.
// Please raise an issue.
type ErrInvalidRouteFileFormat struct {
	row string
}

func (*ErrCantParse) Error() string {
	return "can't parse route table"
}

func (*ErrNotImplemented) Error() string {
	return "not implemented for OS: " + runtime.GOOS
}

func (e *ErrInvalidRouteFileFormat) Error() string {
	return fmt.Sprintf("invalid row %q in route file", e.row)
}

// ErrNoIP indicates that the library could not obtain an IP matching the
// provided kind.
var ErrNoIP = fmt.Errorf("could not find IP matching provided kind")
