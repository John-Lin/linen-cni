package ovsdb

import (
	"os/exec"
	"regexp"
)

func ovsdbUnixPath() string {
	var path string
	var re = regexp.MustCompile(`--remote=punix:(/.*openvswitch.*sock)`)
	cmd := "pgrep -f -a ovsdb-server"
	out, _ := exec.Command("sh", "-c", cmd).Output()
	if len(out) != 0 {
		str := string(out[:])
		strSubMatch := re.FindStringSubmatch(str)
		if len(strSubMatch) != 0 {
			path = strSubMatch[1]
		}
	}
	return path
}
