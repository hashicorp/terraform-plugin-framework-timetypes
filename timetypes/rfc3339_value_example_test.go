// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package timetypes_test

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
)

type TimeResourceModel struct {
	Timestamp timetypes.RFC3339 `tfsdk:"timestamp"`
}

func ExampleRFC3339_ValueRFC3339Time() {
	// For example purposes, typically the data model would be populated automatically by Plugin Framework via Config, Plan or State.
	// https://developer.hashicorp.com/terraform/plugin/framework/handling-data/accessing-values
	data := TimeResourceModel{
		Timestamp: timetypes.NewRFC3339ValueMust("2023-07-25T23:43:16Z"),
	}

	// Check that the RFC3339 data is known and able to be converted to time.Time
	if !data.Timestamp.IsNull() && !data.Timestamp.IsUnknown() {
		t, diags := data.Timestamp.ValueRFC3339Time()
		if diags.HasError() {
			return
		}

		// Output: 2023-07-25T23:43:16Z
		fmt.Println(t.Format(time.RFC3339))
	}
}
