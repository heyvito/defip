package defip

import (
	"encoding/binary"
	"encoding/hex"
	"net/netip"
	"os"
	"slices"
	"strings"
)

var (
	routeV4 = "/proc/net/route"
	routeV6 = "/proc/net/ipv6_route"
)

/* ipv6_route:
00000000000000000000000000000000 00 00000000000000000000000000000000 00 00000000000000000000000000000000 ffffffff 00000001 00000001 00200200 lo
+------------------------------+ ++ +------------------------------+ ++ +------------------------------+ +------+ +------+ +------+ +------+ ++
|                                |  |                                |  |                                |        |        |        |        |
1                                2  3                                4  5                                6        7        8        9        10

  1. IPv6 destination network displayed in 32 hexadecimal chars without colons as separator
  2. IPv6 destination prefix length in hexadecimal
  3. IPv6 source network displayed in 32 hexadecimal chars without colons as separator
  4. IPv6 source prefix length in hexadecimal
  5. IPv6 next hop displayed in 32 hexadecimal chars without colons as separator
  6. Metric in hexadecimal
  7. Reference counter
  8. Use counter
  9. Flags
  10. Device name
*/

func ip6FromHex(in string) (ip netip.Addr, ok bool) {
	if len(in) != 32 {
		ok = false
		return
	}
	v, err := hex.DecodeString(in)
	if err != nil {
		ok = false
	}
	ip = netip.AddrFrom16([16]byte(v))
	ok = true
	return
}

func parseSingleRouteIPv6(fields []string) *NetRoute {
	dstNet, ok := ip6FromHex(fields[0])
	if !ok {
		return nil
	}
	nextHop, ok := ip6FromHex(fields[4])
	if !ok {
		return nil
	}
	rawFlags, err := hex.DecodeString(fields[8])
	if err != nil {
		return nil
	}
	flags := routeTableFlag(binary.BigEndian.Uint32(rawFlags))

	ifName := fields[9]
	return &NetRoute{
		Kind:        NetRouteKindV6,
		Destination: dstNet,
		Flags:       flags.String(),
		Netif:       ifName,
		Gateway:     nextHop,
	}
}

func getRoutesIPv6(source string) (NetRouteList, error) {
	f, err := os.ReadFile(source)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var routes NetRouteList
	lines := strings.Split(string(f), "\n")
	for _, v := range lines {
		v = strings.TrimSpace(v)
		if len(v) == 0 {
			continue
		}
		fields := strings.Fields(v)
		if len(fields) != 10 {
			return nil, &ErrInvalidRouteFileFormat{row: v}
		}
		item := parseSingleRouteIPv6(fields)
		if item == nil {
			return nil, &ErrInvalidRouteFileFormat{row: v}
		}
		routes = append(routes, *item)
	}

	return routes, nil
}

/*
Iface	Destination	Gateway 	Flags	RefCnt	Use	Metric	Mask		MTU	Window	IRTT
ens37	00000000	0101000A	0003	0	0	100	00000000	0	0	0
ens34	00000000	0101A8C0	0003	0	0	100	00000000	0	0	0
ens37	0000000A	00000000	0001	0	0	100	0000FFFF	0	0	0
ens37	0101000A	00000000	0005	0	0	100	FFFFFFFF	0	0	0
ens34	0001A8C0	00000000	0001	0	0	100	00FFFFFF	0	0	0
ens34	0101A8C0	00000000	0005	0	0	100	FFFFFFFF	0	0	0
*/

func ip4FromHex(in string) (ip netip.Addr, ok bool) {
	if len(in) != 8 {
		ok = false
		return
	}
	v, err := hex.DecodeString(in)
	if err != nil {
		ok = false
	}
	slices.Reverse(v)
	ip = netip.AddrFrom4([4]byte(v))
	ok = true
	return
}

func getRoutesIPv4(source string) (NetRouteList, error) {
	f, err := os.ReadFile(source)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var routes NetRouteList
	lines := strings.Split(string(f), "\n")
	if len(lines) < 1 {
		return nil, &ErrCantParse{}
	}
	fields := fieldSet(strings.Fields(lines[0]))
	ifNameIdx := fields.fieldIdx("Iface")
	dstNetIdx := fields.fieldIdx("Destination")
	gatewayIdx := fields.fieldIdx("Gateway")
	flagsIdx := fields.fieldIdx("Flags")

	if ifNameIdx == -1 || dstNetIdx == -1 || gatewayIdx == -1 || flagsIdx == -1 {
		return nil, &ErrCantParse{}
	}

	for _, v := range lines[1:] {
		v = strings.TrimSpace(v)
		if len(v) == 0 {
			continue
		}
		fields := strings.Fields(strings.TrimSpace(v))
		if len(fields) < 4 {
			return nil, &ErrInvalidRouteFileFormat{row: v}
		}
		dstNet, ok := ip4FromHex(fields[dstNetIdx])
		if !ok {
			return nil, &ErrInvalidRouteFileFormat{row: v}
		}
		gateway, ok := ip4FromHex(fields[gatewayIdx])
		if !ok {
			return nil, &ErrInvalidRouteFileFormat{row: v}
		}

		rawFlags, err := hex.DecodeString(fields[flagsIdx])
		if err != nil {
			return nil, &ErrInvalidRouteFileFormat{row: v}
		}
		flags := routeTableFlag(binary.BigEndian.Uint16(rawFlags))

		routes = append(routes, NetRoute{
			Kind:        NetRouteKindV4,
			Destination: dstNet,
			Flags:       flags.String(),
			Netif:       fields[ifNameIdx],
			Gateway:     gateway,
		})
	}

	return routes, nil
}
