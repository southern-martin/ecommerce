package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCanTransition(t *testing.T) {
	t.Run("valid transitions", func(t *testing.T) {
		validTransitions := []struct {
			name string
			from ShipmentStatus
			to   ShipmentStatus
		}{
			{"pending to label_created", StatusPending, StatusLabelCreated},
			{"pending to exception", StatusPending, StatusException},
			{"label_created to picked_up", StatusLabelCreated, StatusPickedUp},
			{"label_created to exception", StatusLabelCreated, StatusException},
			{"picked_up to in_transit", StatusPickedUp, StatusInTransit},
			{"picked_up to exception", StatusPickedUp, StatusException},
			{"in_transit to delivered", StatusInTransit, StatusDelivered},
			{"in_transit to exception", StatusInTransit, StatusException},
			{"exception to in_transit", StatusException, StatusInTransit},
			{"exception to delivered", StatusException, StatusDelivered},
		}

		for _, tc := range validTransitions {
			t.Run(tc.name, func(t *testing.T) {
				assert.True(t, CanTransition(tc.from, tc.to),
					"expected transition from %s to %s to be valid", tc.from, tc.to)
			})
		}
	})

	t.Run("invalid transitions", func(t *testing.T) {
		invalidTransitions := []struct {
			name string
			from ShipmentStatus
			to   ShipmentStatus
		}{
			{"delivered to pending", StatusDelivered, StatusPending},
			{"pending to delivered", StatusPending, StatusDelivered},
			{"pending to in_transit", StatusPending, StatusInTransit},
			{"pending to picked_up", StatusPending, StatusPickedUp},
			{"label_created to delivered", StatusLabelCreated, StatusDelivered},
			{"label_created to in_transit", StatusLabelCreated, StatusInTransit},
			{"picked_up to delivered", StatusPickedUp, StatusDelivered},
			{"picked_up to label_created", StatusPickedUp, StatusLabelCreated},
			{"in_transit to pending", StatusInTransit, StatusPending},
			{"in_transit to picked_up", StatusInTransit, StatusPickedUp},
		}

		for _, tc := range invalidTransitions {
			t.Run(tc.name, func(t *testing.T) {
				assert.False(t, CanTransition(tc.from, tc.to),
					"expected transition from %s to %s to be invalid", tc.from, tc.to)
			})
		}
	})

	t.Run("unknown status returns false", func(t *testing.T) {
		assert.False(t, CanTransition(ShipmentStatus("nonexistent"), StatusPending))
		assert.False(t, CanTransition(ShipmentStatus(""), StatusPending))
	})

	t.Run("same status returns false", func(t *testing.T) {
		statuses := []ShipmentStatus{
			StatusPending,
			StatusLabelCreated,
			StatusPickedUp,
			StatusInTransit,
			StatusDelivered,
			StatusException,
		}
		for _, s := range statuses {
			t.Run(string(s), func(t *testing.T) {
				assert.False(t, CanTransition(s, s),
					"expected transition from %s to %s to be invalid (same status)", s, s)
			})
		}
	})

	t.Run("delivered has no outgoing transitions", func(t *testing.T) {
		allStatuses := []ShipmentStatus{
			StatusPending, StatusLabelCreated, StatusPickedUp,
			StatusInTransit, StatusDelivered, StatusException,
		}
		for _, to := range allStatuses {
			assert.False(t, CanTransition(StatusDelivered, to),
				"delivered should not transition to %s", to)
		}
	})
}
