//go:build !(darwin || linux)

package defip

func init() {
	getRoutes = func() (NetRouteList, error) {
		return nil, ErrNotImplemented{}
	}
}
