//go:build !linux

package netx

func binaryNewerThanTCPPort(string, int) bool {
	return false
}
