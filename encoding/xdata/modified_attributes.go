package xdata

import (
	"strings"
)

type ModifiedAttribute int

const (
	NotificationsEnabled ModifiedAttribute = 1 << iota
	ActiveChecksEnabled
	PassiveChecksEnabled
	EventHandlerEnabled
	FlapDetectionEnabled
	FailurePredictionEnabled
	PerformanceDataEnabled
	ObsessiveHandlerEnabled
	EventHandlerCommand
	CheckCommand
	NormalCheckInterval
	RetryCheckInterval
	MaxCheckAttempts
	FreshnessChecksEnabled
	CheckTimeperiod
	CustomVariable
	NotificationTimeperiod
)

func (a ModifiedAttribute) String() string {
	attrs := make([]string, 0)
	if a&NotificationsEnabled != 0 {
		attrs = append(attrs, "Notifications Enabled")
	}

	if a&ActiveChecksEnabled != 0 {
		attrs = append(attrs, "Active Checks Enabled")
	}

	if a&PassiveChecksEnabled != 0 {
		attrs = append(attrs, "Passive Checks Enabled")
	}

	if a&EventHandlerEnabled != 0 {
		attrs = append(attrs, "Event Handler Enabled")
	}

	if a&FlapDetectionEnabled != 0 {
		attrs = append(attrs, "Flap Detection Enabled")
	}

	if a&FailurePredictionEnabled != 0 {
		attrs = append(attrs, "Failure Prediction Enabled")
	}

	if a&PerformanceDataEnabled != 0 {
		attrs = append(attrs, "Performance Data Enabled")
	}

	if a&ObsessiveHandlerEnabled != 0 {
		attrs = append(attrs, "Obsessive Handler Enabled")
	}

	if a&EventHandlerCommand != 0 {
		attrs = append(attrs, "Event Handler Command")
	}

	if a&CheckCommand != 0 {
		attrs = append(attrs, "Check Command")
	}

	if a&NormalCheckInterval != 0 {
		attrs = append(attrs, "Normal Check Interval")
	}

	if a&RetryCheckInterval != 0 {
		attrs = append(attrs, "Retry Check Interval")
	}

	if a&MaxCheckAttempts != 0 {
		attrs = append(attrs, "Max Check Attempts")
	}

	if a&FreshnessChecksEnabled != 0 {
		attrs = append(attrs, "Freshness Checks Enabled")
	}

	if a&CheckTimeperiod != 0 {
		attrs = append(attrs, "Check Timeperiod")
	}

	if a&CustomVariable != 0 {
		attrs = append(attrs, "Custom Variable")
	}

	if a&NotificationTimeperiod != 0 {
		attrs = append(attrs, "Notification Timeperiod")
	}

	return strings.Join(attrs, ", ")
}
