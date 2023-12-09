package defip

/* Keep this in sync with /usr/src/linux/include/linux/ipv6_route.h */

type routeTableFlag uint32

func (r routeTableFlag) Is(other routeTableFlag) bool { return r&other == other }

const (
	// rtfDefault indicates the default route, typically learned via Neighbor
	// Discovery (ND) protocol in IPv6
	rtfDefault routeTableFlag = 0x00010000

	// rtfAllOnLink indicates a fallback route when no routers are present on
	// the link
	rtfAllOnLink routeTableFlag = 0x00020000

	// rtfAddrConf indicates an address configuration route, usually set up by
	// Router Advertisements (RA)
	rtfAddrConf routeTableFlag = 0x00040000

	// rtfNoNextHop indicates a route that does not have a next hop defined,
	// implying direct delivery
	rtfNoNextHop routeTableFlag = 0x00200000

	// rtfExpires indicates a route that has a limited lifetime, after which it
	// expires
	rtfExpires routeTableFlag = 0x00400000

	// rtfCache indicates a cached route, often dynamically learned and stored
	// for efficiency
	rtfCache routeTableFlag = 0x01000000

	// rtfFlow indicates a flow significant route, used for advanced routing and
	// traffic control
	rtfFlow routeTableFlag = 0x02000000

	// rtfPolicy indicates a policy-based route, used for implementing policy
	// routing
	rtfPolicy routeTableFlag = 0x04000000

	// rtfLocal indicates a local route, used for routing local traffic within
	// the system
	rtfLocal routeTableFlag = 0x80000000
)

/* Keep this in sync with /usr/src/linux/include/linux/route.h */

const (
	// rtfUp indicates the route is usable
	rtfUp routeTableFlag = 0x0001

	// rtfGateway indicates the destination is a gateway
	rtfGateway routeTableFlag = 0x0002

	// rtfHost indicates a host-specific route (as opposed to a network route)
	rtfHost routeTableFlag = 0x0004

	// rtfReinstate indicates the route should be reinstated after a timeout
	rtfReinstate routeTableFlag = 0x0008

	// rtfDynamic indicates the route was created dynamically, typically by a
	// redirect
	rtfDynamic routeTableFlag = 0x0010

	// rtfModified indicates the route was modified dynamically, typically by a
	// redirect
	rtfModified routeTableFlag = 0x0020

	// rtfMTU indicates a specific MTU (Maximum Transmission Unit) for this
	// route
	rtfMTU routeTableFlag = 0x0040

	// rtfWindow indicates per-route TCP window clamping
	rtfWindow routeTableFlag = 0x0080

	// rtfIRTT indicates the initial round-trip time for this route
	rtfIRTT routeTableFlag = 0x0100

	// rtfReject indicates the route should be rejected
	rtfReject routeTableFlag = 0x0200
)

/* this is a 2.0.36 flag from /usr/src/linux/include/linux/route.h */

const (
	// rtfNotCache indicates a route that should not be cached. This is used
	// to prevent caching of specific routes, usually for routes that are
	// dynamic or special in nature.
	rtfNotCache routeTableFlag = 0x0400
)

func (r routeTableFlag) String() string {
	val := ""
	if r.Is(rtfUp) {
		val += "U"
	}
	if r.Is(rtfGateway) {
		val += "G"
	}
	if r.Is(rtfReject) {
		val += "!"
	}
	if r.Is(rtfHost) {
		val += "H"
	}
	if r.Is(rtfReinstate) {
		val += "R"
	}
	if r.Is(rtfDynamic) {
		val += "D"
	}
	if r.Is(rtfModified) {
		val += "M"
	}
	if r.Is(rtfDefault) {
		val += "d"
	}
	if r.Is(rtfAllOnLink) {
		val += "a"
	}
	if r.Is(rtfAddrConf) {
		val += "c"
	}
	if r.Is(rtfNoNextHop) {
		val += "o"
	}
	if r.Is(rtfExpires) {
		val += "e"
	}
	if r.Is(rtfCache) {
		val += "c"
	}
	if r.Is(rtfFlow) {
		val += "f"
	}
	if r.Is(rtfPolicy) {
		val += "p"
	}
	if r.Is(rtfLocal) {
		val += "l"
	}
	if r.Is(rtfMTU) {
		val += "u"
	}
	if r.Is(rtfWindow) {
		val += "w"
	}
	if r.Is(rtfIRTT) {
		val += "i"
	}
	if r.Is(rtfNotCache) {
		val += "n"
	}
	return val
}
