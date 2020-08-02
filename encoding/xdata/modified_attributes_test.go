package xdata

import (
	"testing"
)

func TestModifiedAttribute_String(t *testing.T) {
	tests := []struct {
		input    ModifiedAttribute
		expected string
	}{
		{ModifiedAttribute(0), ""},
		{NotificationsEnabled, "Notifications Enabled"},
		{ActiveChecksEnabled, "Active Checks Enabled"},
		{PassiveChecksEnabled, "Passive Checks Enabled"},
		{EventHandlerEnabled, "Event Handler Enabled"},
		{FlapDetectionEnabled, "Flap Detection Enabled"},
		{FailurePredictionEnabled, "Failure Prediction Enabled"},
		{PerformanceDataEnabled, "Performance Data Enabled"},
		{ObsessiveHandlerEnabled, "Obsessive Handler Enabled"},
		{EventHandlerCommand, "Event Handler Command"},
		{CheckCommand, "Check Command"},
		{NormalCheckInterval, "Normal Check Interval"},
		{RetryCheckInterval, "Retry Check Interval"},
		{MaxCheckAttempts, "Max Check Attempts"},
		{FreshnessChecksEnabled, "Freshness Checks Enabled"},
		{CheckTimeperiod, "Check Timeperiod"},
		{CustomVariable, "Custom Variable"},
		{NotificationTimeperiod, "Notification Timeperiod"},
		{NotificationsEnabled | ActiveChecksEnabled, "Notifications Enabled, Active Checks Enabled"},
	}

	for _, test := range tests {
		s := test.input.String()
		if s != test.expected {
			t.Errorf("unexpected string, got: '%s', want: '%s'", s, test.expected)
		}
	}
}
