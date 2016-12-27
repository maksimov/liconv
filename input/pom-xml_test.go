package input

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserHomeDir(t *testing.T) {
	var tests = []struct {
		name     string
		fixture  func()
		expected string
	}{
		{"Windows", func() {
			savedFunc := isWindows
			os.Setenv("HOMEDRIVE", "C:")
			os.Setenv("HOMEPATH", "/Blah")
			isWindows = func() bool {
				isWindows = savedFunc
				return true
			}
		}, "C:/Blah"},
		{"WindowsUserProfile", func() {
			savedFunc := isWindows
			os.Setenv("USERPROFILE", "C:/Users/blah")
			isWindows = func() bool {
				isWindows = savedFunc
				return true
			}
		}, "C:/Users/blah"},
		{"Unix", func() {
			savedFunc := isWindows
			os.Setenv("HOME", "/home/blah")
			isWindows = func() bool {
				isWindows = savedFunc
				return false
			}
		}, "/home/blah"},
	}

	for _, tc := range tests {
		tc.fixture()
		actual := userHomeDir()
		assert.Equal(t, tc.expected, actual, fmt.Sprintf("Test case %s", tc.name))
		os.Clearenv()
	}
}
