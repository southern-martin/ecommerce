package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCanReturnTransition(t *testing.T) {
	t.Run("valid transitions", func(t *testing.T) {
		validCases := []struct {
			name string
			from ReturnStatus
			to   ReturnStatus
		}{
			{"requested to approved", ReturnStatusRequested, ReturnStatusApproved},
			{"requested to rejected", ReturnStatusRequested, ReturnStatusRejected},
			{"approved to shipped_back", ReturnStatusApproved, ReturnStatusShippedBack},
			{"shipped_back to received", ReturnStatusShippedBack, ReturnStatusReceived},
			{"received to refunded", ReturnStatusReceived, ReturnStatusRefunded},
		}

		for _, tc := range validCases {
			t.Run(tc.name, func(t *testing.T) {
				assert.True(t, CanReturnTransition(tc.from, tc.to),
					"expected transition from %s to %s to be valid", tc.from, tc.to)
			})
		}
	})

	t.Run("invalid transitions", func(t *testing.T) {
		invalidCases := []struct {
			name string
			from ReturnStatus
			to   ReturnStatus
		}{
			{"requested to shipped_back", ReturnStatusRequested, ReturnStatusShippedBack},
			{"requested to received", ReturnStatusRequested, ReturnStatusReceived},
			{"requested to refunded", ReturnStatusRequested, ReturnStatusRefunded},
			{"approved to approved", ReturnStatusApproved, ReturnStatusApproved},
			{"approved to rejected", ReturnStatusApproved, ReturnStatusRejected},
			{"rejected to approved", ReturnStatusRejected, ReturnStatusApproved},
			{"rejected to shipped_back", ReturnStatusRejected, ReturnStatusShippedBack},
			{"refunded to requested", ReturnStatusRefunded, ReturnStatusRequested},
			{"refunded to approved", ReturnStatusRefunded, ReturnStatusApproved},
			{"received to approved", ReturnStatusReceived, ReturnStatusApproved},
			{"shipped_back to approved", ReturnStatusShippedBack, ReturnStatusApproved},
		}

		for _, tc := range invalidCases {
			t.Run(tc.name, func(t *testing.T) {
				assert.False(t, CanReturnTransition(tc.from, tc.to),
					"expected transition from %s to %s to be invalid", tc.from, tc.to)
			})
		}
	})

	t.Run("unknown status returns false", func(t *testing.T) {
		assert.False(t, CanReturnTransition(ReturnStatus("unknown"), ReturnStatusApproved))
		assert.False(t, CanReturnTransition(ReturnStatusRequested, ReturnStatus("unknown")))
		assert.False(t, CanReturnTransition(ReturnStatus("foo"), ReturnStatus("bar")))
	})
}
