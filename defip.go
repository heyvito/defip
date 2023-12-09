package defip

import (
	"cmp"
	"fmt"
	"net"
	"net/netip"
	"slices"
	"strings"
)

type NetRouteKind uint8

const (
	NetRouteKindV4 NetRouteKind = iota + 1
	NetRouteKindV6
)

func (n NetRouteKind) String() string {
	switch n {
	case NetRouteKindV4:
		return "IPv4"
	case NetRouteKindV6:
		return "IPv6"
	}
	panic("Invalid NetRouteKind")
}

type NetRoute struct {
	Kind        NetRouteKind
	Destination netip.Addr
	Flags       string
	Netif       string
	Gateway     netip.Addr
}

func (n NetRoute) HasFlags(flags ...string) bool {
	for _, v := range flags {
		if !strings.Contains(n.Flags, v) {
			return false
		}
	}

	return true
}

type NetRouteList []NetRoute

var filterRoute = func(r *NetRoute) bool {
	return r.HasFlags("U", "G") &&
		!r.HasFlags("H")
}

func (n NetRouteList) FindDefaults(kind NetRouteKind) []NetRoute {
	var result []NetRoute

	for _, v := range n {
		if v.Kind == kind && filterRoute(&v) {
			result = append(result, v)
		}
	}

	return result
}

var getRoutes func() (NetRouteList, error) = nil

// FindRoutes returns a list of detected routes to default gateways
func FindRoutes() (NetRouteList, error) {
	return getRoutes()
}

func filter[S interface{ ~[]E }, E any](set S, fn func(i E) bool) S {
	var r S
	for _, v := range set {
		if fn(v) {
			r = append(r, v)
		}
	}

	return r
}

// FindDefaultIP attempts to find an IP of given NetRouteKind that's most likely
// connected to wider network. Returns ErrNoIP in case no IP with the given kind
// can be detected.
func FindDefaultIP(kind NetRouteKind) (*netip.Addr, error) {
	routes, err := FindRoutes()
	if err != nil {
		panic(err)
	}

	routes = filter(routes, func(i NetRoute) bool {
		return i.HasFlags("U", "G")
	})

	ifaces := map[string]bool{}
	for _, v := range routes {
		ifaces[v.Netif] = true
	}

	var addrs []netip.Addr
	for name := range ifaces {
		iface, err := net.InterfaceByName(name)
		if err != nil {
			return nil, fmt.Errorf("could not get interface `%s': %w", name, err)
		}

		ips, err := iface.Addrs()
		if err != nil {
			return nil, fmt.Errorf("could not get IPs for interface `%s': %w", name, err)
		}

		for _, v := range ips {
			rawAdd, ok := v.(*net.IPNet)
			if !ok {
				continue
			}

			if err != nil {
				continue
			}

			var add netip.Addr
			if v4 := rawAdd.IP.To4(); v4 != nil {
				if kind != NetRouteKindV4 {
					continue
				}
				add = netip.AddrFrom4([4]byte(v4))
			} else {
				if kind != NetRouteKindV6 {
					continue
				}
				add = netip.AddrFrom16([16]byte(rawAdd.IP))
			}
			add = add.WithZone(name)
			addrs = append(addrs, add)
		}
	}

	if ip := selectIP(kind, addrs); ip != nil {
		return ip, nil
	}

	return nil, ErrNoIP
}

type ipWeight struct {
	addr   netip.Addr
	weight int
}

var ulaEnd = netip.MustParseAddr("fdff:ffff:ffff:ffff:ffff:ffff:ffff:ffff")
var ulaStart = netip.MustParseAddr("fd00::")

func isULA(addr netip.Addr) bool {
	return addr.Compare(ulaStart) >= 0 && addr.Compare(ulaEnd) <= 0
}

func sortWeighted(list []netip.Addr) {
	if len(list) == 0 {
		return
	}

	weightList := make([]ipWeight, len(list))
	for i, v := range list {
		weight := 0

		if isULA(v) {
			weight += 2
		}

		if v.IsPrivate() {
			weight += 1
		}
		if v.IsGlobalUnicast() {
			weight += 1
		}
		weightList[i].weight = weight
		weightList[i].addr = v

		fmt.Printf("IP %s has weight %d\n", v, weight)
	}

	slices.SortFunc(weightList, func(a, b ipWeight) int {
		return cmp.Compare(b.weight, a.weight)
	})

	for i, v := range weightList {
		list[i] = v.addr
	}

}

func selectIP(kind NetRouteKind, list []netip.Addr) *netip.Addr {
	list = filter(list, func(i netip.Addr) bool {
		return (kind == NetRouteKindV6 && i.Is6()) ||
			(kind == NetRouteKindV4 && i.Is4())
	})

	sortWeighted(list)

	return &list[0]
}
