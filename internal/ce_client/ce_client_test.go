package ce_client

import (
	arcalog "go.arcalot.io/log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func IntMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func TestIntMinBasic(t *testing.T) {
	assert.Equal(t, -2, IntMin(2, -2))
}

func TestNewCeClient(t *testing.T) {
	logger := arcalog.NewLogger(arcalog.LevelInfo, arcalog.NewNOOPLogger())
	testCases := map[string]struct {
		choice         string
		expectedResult bool
	}{
		"podman": {
			"podman",
			false,
		},
		"docker": {
			"docker",
			true,
		},
	}
	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			cec, err := NewCeClient(tc.choice)

			if err != nil {
				logger.Errorf("(%w)", err)
			}
			_, ok := cec.(ContainerEngineClient)
			assert.Equal(t, tc.expectedResult, ok)
		})
	}
}
