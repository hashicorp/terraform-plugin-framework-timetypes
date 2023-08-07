// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package timetypes_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
)

func TestRFC3339_StringSemanticEquals(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		currentRFC3339time timetypes.RFC3339
		givenRFC3339time   basetypes.StringValuable
		expectedMatch      bool
		expectedDiags      diag.Diagnostics
	}{
		"not equal - different dates": {
			currentRFC3339time: timetypes.NewRFC3339ValueMust("2023-07-25T23:43:16Z"),
			givenRFC3339time:   timetypes.NewRFC3339ValueMust("2023-07-26T23:43:16Z"),
			expectedMatch:      false,
		},
		"not equal - different times": {
			currentRFC3339time: timetypes.NewRFC3339ValueMust("2023-07-25T23:43:16Z"),
			givenRFC3339time:   timetypes.NewRFC3339ValueMust("2023-07-25T23:01:16Z"),
			expectedMatch:      false,
		},
		"not equal - different offset times": {
			currentRFC3339time: timetypes.NewRFC3339ValueMust("2023-07-25T23:43:16Z"),
			givenRFC3339time:   timetypes.NewRFC3339ValueMust("2023-07-25T23:43:16+03:00"),
			expectedMatch:      false,
		},
		"not equal - UTC time and local time": {
			currentRFC3339time: timetypes.NewRFC3339ValueMust("2023-07-25T23:43:16Z"),
			givenRFC3339time:   timetypes.NewRFC3339ValueMust("2023-07-25T20:43:16-03:00"),
			expectedMatch:      false,
		},
		"semantically equal - Z suffix and positive zero num offset": {
			currentRFC3339time: timetypes.NewRFC3339ValueMust("2023-07-25T23:43:16Z"),
			givenRFC3339time:   timetypes.NewRFC3339ValueMust("2023-07-25T23:43:16+00:00"),
			expectedMatch:      true,
		},
		"semantically equal - Z suffix and negative zero num offset": {
			currentRFC3339time: timetypes.NewRFC3339ValueMust("2023-07-25T23:43:16Z"),
			givenRFC3339time:   timetypes.NewRFC3339ValueMust("2023-07-25T23:43:16-00:00"),
			expectedMatch:      true,
		},
		"semantically equal - negative zero and positive zero num offset": {
			currentRFC3339time: timetypes.NewRFC3339ValueMust("2023-07-25T23:43:16-00:00"),
			givenRFC3339time:   timetypes.NewRFC3339ValueMust("2023-07-25T23:43:16+00:00"),
			expectedMatch:      true,
		},
		"semantically equal - byte for byte match": {
			currentRFC3339time: timetypes.NewRFC3339ValueMust("2023-07-25T23:43:16Z"),
			givenRFC3339time:   timetypes.NewRFC3339ValueMust("2023-07-25T23:43:16Z"),
			expectedMatch:      true,
		},
		"error - not given RFC3339 value": {
			currentRFC3339time: timetypes.NewRFC3339ValueMust("2023-07-25T23:43:16Z"),
			givenRFC3339time:   basetypes.NewStringValue("0000-00-00T00:00:00-00:00"),
			expectedMatch:      false,
			expectedDiags: diag.Diagnostics{
				diag.NewErrorDiagnostic(
					"Semantic Equality Check Error",
					"An unexpected value type was received while performing semantic equality checks. "+
						"Please report this to the provider developers.\n\n"+
						"Expected Value Type: timetypes.RFC3339\n"+
						"Got Value Type: basetypes.StringValue",
				),
			},
		},
	}
	for name, testCase := range testCases {
		name, testCase := name, testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			match, diags := testCase.currentRFC3339time.StringSemanticEquals(context.Background(), testCase.givenRFC3339time)

			if testCase.expectedMatch != match {
				t.Errorf("Expected StringSemanticEquals to return: %t, but got: %t", testCase.expectedMatch, match)
			}

			if diff := cmp.Diff(diags, testCase.expectedDiags); diff != "" {
				t.Errorf("Unexpected diagnostics (-got, +expected): %s", diff)
			}
		})
	}
}

func TestRFC3339_ValueRFC3339Time(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		RFC3339           timetypes.RFC3339
		expectedTimestamp string
		expectedDiags     diag.Diagnostics
	}{
		"RFC3339 string value is null ": {
			RFC3339: timetypes.NewRFC3339Null(),
			expectedDiags: diag.Diagnostics{
				diag.NewErrorDiagnostic(
					"RFC3339 ValueRFC3339Time Error",
					"RFC3339 string value is null",
				),
			},
		},
		"RFC3339 string value is unknown ": {
			RFC3339: timetypes.NewRFC3339Unknown(),
			expectedDiags: diag.Diagnostics{
				diag.NewErrorDiagnostic(
					"RFC3339 ValueRFC3339Time Error",
					"RFC3339 string value is unknown",
				),
			},
		},
		"valid RFC3339 Timestamp - Zulu suffix": {
			RFC3339:           timetypes.NewRFC3339ValueMust("2023-07-25T23:43:16Z"),
			expectedTimestamp: "2023-07-25T23:43:16Z",
		},
		"valid RFC3339 Timestamp - UTC offset ": {
			RFC3339:           timetypes.NewRFC3339ValueMust("2023-07-25T23:43:16-00:00"),
			expectedTimestamp: "2023-07-25T23:43:16-00:00",
		},
		"valid RFC3339 Timestamp - EDT offset ": {
			RFC3339:           timetypes.NewRFC3339ValueMust("2023-07-25T23:43:16-04:00"),
			expectedTimestamp: "2023-07-25T23:43:16-04:00",
		},
	}
	for name, testCase := range testCases {
		name, testCase := name, testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			rfc3339Time, diags := testCase.RFC3339.ValueRFC3339Time()
			expectedRFC3339Time, _ := time.Parse(time.RFC3339, testCase.expectedTimestamp)

			if rfc3339Time != expectedRFC3339Time {
				t.Errorf("Unexpected difference in time.Time, got: %s, expected: %s", rfc3339Time, expectedRFC3339Time)
			}

			if diff := cmp.Diff(diags, testCase.expectedDiags); diff != "" {
				t.Errorf("Unexpected diagnostics (-got, +expected): %s", diff)
			}
		})
	}
}
