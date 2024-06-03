// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package timetypes_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/attr/xattr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
)

func TestDuration_Equals(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		currentDuration timetypes.GoDuration
		givenDuration   basetypes.StringValuable
		expectedMatch   bool
	}{
		"not equal - different durations": {
			currentDuration: timetypes.NewGoDurationValueFromStringMust("50s"),
			givenDuration:   timetypes.NewGoDurationValueFromStringMust("50m"),
			expectedMatch:   false,
		},
		"equal - exactly the same string": {
			currentDuration: timetypes.NewGoDurationValueFromStringMust("30h22m33s"),
			givenDuration:   timetypes.NewGoDurationValueFromStringMust("30h22m33s"),
			expectedMatch:   true,
		},
		"equal - same duration expressed differently": {
			currentDuration: timetypes.NewGoDurationValueFromStringMust("3h25m63s"),
			givenDuration:   timetypes.NewGoDurationValueFromStringMust("12363s"),
			expectedMatch:   true,
		},
		"error - not a Duration value": {
			currentDuration: timetypes.NewGoDurationValueFromStringMust("56s"),
			givenDuration:   basetypes.NewStringValue("abcdef"),
			expectedMatch:   false,
		},
	}
	for name, testCase := range testCases {
		name, testCase := name, testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			match := testCase.currentDuration.Equal(testCase.givenDuration)

			if testCase.expectedMatch != match {
				t.Errorf("Expected StringSemanticEquals to return: %t, but got: %t", testCase.expectedMatch, match)
			}
		})
	}
}

func TestDurationValidateAttribute(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		Duration      timetypes.GoDuration
		expectedDiags diag.Diagnostics
	}{
		"empty-struct": {
			Duration: timetypes.GoDuration{},
		},
		"null": {
			Duration: timetypes.NewGoDurationNull(),
		},
		"unknown": {
			Duration: timetypes.NewGoDurationUnknown(),
		},
		"valid duration": {
			Duration: timetypes.NewGoDurationValueFromStringMust("42h"),
		},
		"invalid duration": {
			Duration: timetypes.GoDuration{
				StringValue: basetypes.NewStringValue("nope"),
			},
			expectedDiags: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid Time Duration String Value",
					"A string value was provided that is not a valid Go Time Duration string format. "+
						`A duration string is a sequence of numbers, each with optional fraction and a unit suffix, such as "300ms", "-1.5h" or "2h45m". `+
						`Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".\n\n`+
						"Given Value: nope\n"+
						"Error: time: invalid duration \"nope\"",
				),
			},
		},
	}
	for name, testCase := range testCases {
		name, testCase := name, testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := xattr.ValidateAttributeResponse{}

			testCase.Duration.ValidateAttribute(
				context.Background(),
				xattr.ValidateAttributeRequest{
					Path: path.Root("test"),
				},
				&resp,
			)

			if diff := cmp.Diff(resp.Diagnostics, testCase.expectedDiags); diff != "" {
				t.Errorf("Unexpected diagnostics (-got, +expected): %s", diff)
			}
		})
	}
}

func TestDurationValidateParameter(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		Duration        timetypes.GoDuration
		expectedFuncErr *function.FuncError
	}{
		"empty-struct": {
			Duration: timetypes.GoDuration{},
		},
		"null": {
			Duration: timetypes.NewGoDurationNull(),
		},
		"unknown": {
			Duration: timetypes.NewGoDurationUnknown(),
		},
		"valid duration": {
			Duration: timetypes.NewGoDurationValueFromStringMust("42h"),
		},
		"invalid duration": {
			Duration: timetypes.GoDuration{
				StringValue: basetypes.NewStringValue("nope"),
			},
			expectedFuncErr: function.NewArgumentFuncError(
				0,
				"Invalid Go Time Duration String Value: "+
					"A string value was provided that is not a valid Go Time Duration string format. "+
					`A duration string is a sequence of numbers, each with optional fraction and a unit suffix, such as "300ms", "-1.5h" or "2h45m". `+
					`Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".\n\n`+
					"Given Value: nope\n"+
					"Error: time: invalid duration \"nope\"",
			),
		},
	}
	for name, testCase := range testCases {
		name, testCase := name, testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := function.ValidateParameterResponse{}

			testCase.Duration.ValidateParameter(
				context.Background(),
				function.ValidateParameterRequest{
					Position: int64(0),
				},
				&resp,
			)

			if diff := cmp.Diff(resp.Error, testCase.expectedFuncErr); diff != "" {
				t.Errorf("Unexpected diagnostics (-got, +expected): %s", diff)
			}
		})
	}
}

func TestDuration_ValueDuration(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		Duration         timetypes.GoDuration
		expectedDuration string
		expectedDiags    diag.Diagnostics
	}{
		"Duration string value is null ": {
			Duration: timetypes.NewGoDurationNull(),
			expectedDiags: diag.Diagnostics{
				diag.NewErrorDiagnostic(
					"Go Duration ValueDuration Error",
					"Duration string value is null",
				),
			},
		},
		"Duration string value is unknown ": {
			Duration: timetypes.NewGoDurationUnknown(),
			expectedDiags: diag.Diagnostics{
				diag.NewErrorDiagnostic(
					"Go Duration ValueDuration Error",
					"Duration string value is unknown",
				),
			},
		},
		"valid duration": {
			Duration:         timetypes.NewGoDurationValueFromStringMust("1h1m1s"),
			expectedDuration: "1h1m1s",
		},
	}
	for name, testCase := range testCases {
		name, testCase := name, testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			duration, diags := testCase.Duration.ValueGoDuration()
			expectedDuration, _ := time.ParseDuration(testCase.expectedDuration)

			if duration != expectedDuration {
				t.Errorf("Unexpected difference in time.Duration, got: %s, expected: %s", duration, expectedDuration)
			}

			if diff := cmp.Diff(diags, testCase.expectedDiags); diff != "" {
				t.Errorf("Unexpected diagnostics (-got, +expected): %s", diff)
			}
		})
	}
}
