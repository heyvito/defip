package defip

import (
	"os/exec"
	"strings"
)

func init() {
	getRoutes = func() (NetRouteList, error) {
		cmd := exec.Command("netstat", "-rn")
		output, err := cmd.CombinedOutput()
		if err != nil {
			return nil, err
		}
		parser := newNetstatParser()
		for _, line := range strings.Split(string(output), "\n") {
			if err = parser.feed(line); err != nil {
				return nil, err
			}
		}
		return parser.netData, nil
	}
}
