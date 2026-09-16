//go:build !linux

package scheduler

// BinaryNewerThanListener is only implemented on Linux.
func BinaryNewerThanListener(string, int) bool {
	return false
}
