package safeid

import (
	"testing"
)

func TestIsGeneric(t *testing.T) {
	tt := []struct {
		name       string
		f          func() bool
		expOutcome bool
	}{
		{
			"generic",
			IsGeneric[Generic],
			true,
		},
		{
			"custom",
			IsGeneric[test],
			false,
		},
		{
			"custom",
			IsGeneric[empty],
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ok := tc.f()
			assertEqual(t, tc.expOutcome, ok)
		})
	}
}
