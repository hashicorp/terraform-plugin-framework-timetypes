// Copyright IBM Corp. 2023, 2025
// SPDX-License-Identifier: MPL-2.0

package timetypes_test

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
)

type DurationResourceModel struct {
	Duration timetypes.GoDuration `tfsdk:"duration"`
}

func ExampleGoDuration_ValueGoDuration() {
	// For example purposes, typically the data model would be populated automatically by Plugin Framework via Config, Plan or State.
	// https://developer.hashicorp.com/terraform/plugin/framework/handling-data/accessing-values
	data := DurationResourceModel{
		Duration: timetypes.NewGoDurationValueFromStringMust("1h2m3s"),
	}

	// Check that the duration data is known and able to be converted to time.Duration
	if !data.Duration.IsNull() && !data.Duration.IsUnknown() {
		t, diags := data.Duration.ValueGoDuration()
		if diags.HasError() {
			return
		}

		// Output: 1h2m3s
		fmt.Println(t.String())
	}
}
