package defip

func init() {
	getRoutes = func() (NetRouteList, error) {
		ip6List, err := getRoutesIPv6(routeV6)
		if err != nil {
			return nil, err
		}

		ip4List, err := getRoutesIPv4(routeV4)
		if err != nil {
			return nil, err
		}

		return append(ip4List, ip6List...), nil
	}
}
