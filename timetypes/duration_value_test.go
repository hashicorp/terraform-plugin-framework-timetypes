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
		currentDuration timetypes.Duration
		givenDuration   basetypes.StringValuable
		expectedMatch   bool
	}{
		"not equal - different durations": {
			currentDuration: timetypes.NewDurationValueFromStringMust("50s"),
			givenDuration:   timetypes.NewDurationValueFromStringMust("50m"),
			expectedMatch:   false,
		},
		"equal - exactly the same string": {
			currentDuration: timetypes.NewDurationValueFromStringMust("30h22m33s"),
			givenDuration:   timetypes.NewDurationValueFromStringMust("30h22m33s"),
			expectedMatch:   true,
		},
		"equal - same duration expressed differently": {
			currentDuration: timetypes.NewDurationValueFromStringMust("3h25m63s"),
			givenDuration:   timetypes.NewDurationValueFromStringMust("12363s"),
			expectedMatch:   true,
		},
		"error - not a Duration value": {
			currentDuration: timetypes.NewDurationValueFromStringMust("56s"),
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
		Duration      timetypes.Duration
		expectedDiags diag.Diagnostics
	}{
		"empty-struct": {
			Duration: timetypes.Duration{},
		},
		"null": {
			Duration: timetypes.NewDurationNull(),
		},
		"unknown": {
			Duration: timetypes.NewDurationUnknown(),
		},
		"valid duration": {
			Duration: timetypes.NewDurationValueFromStringMust("42h"),
		},
		"invalid duration": {
			Duration: timetypes.Duration{
				StringValue: basetypes.NewStringValue("nope"),
			},
			expectedDiags: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid time duration String Value",
					"A string value was provided that is not valid time duration string format.\n\n"+
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
		Duration        timetypes.Duration
		expectedFuncErr *function.FuncError
	}{
		"empty-struct": {
			Duration: timetypes.Duration{},
		},
		"null": {
			Duration: timetypes.NewDurationNull(),
		},
		"unknown": {
			Duration: timetypes.NewDurationUnknown(),
		},
		"valid duration": {
			Duration: timetypes.NewDurationValueFromStringMust("42h"),
		},
		"invalid duration": {
			Duration: timetypes.Duration{
				StringValue: basetypes.NewStringValue("nope"),
			},
			expectedFuncErr: function.NewArgumentFuncError(
				0,
				"Invalid time duration String Value: "+
					"A string value was provided that is not valid time duration string format.\n\n"+
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
		Duration         timetypes.Duration
		expectedDuration string
		expectedDiags    diag.Diagnostics
	}{
		"Duration string value is null ": {
			Duration: timetypes.NewDurationNull(),
			expectedDiags: diag.Diagnostics{
				diag.NewErrorDiagnostic(
					"Duration ValueDuration Error",
					"Duration string value is null",
				),
			},
		},
		"Duration string value is unknown ": {
			Duration: timetypes.NewDurationUnknown(),
			expectedDiags: diag.Diagnostics{
				diag.NewErrorDiagnostic(
					"Duration ValueDuration Error",
					"Duration string value is unknown",
				),
			},
		},
		"valid duration": {
			Duration:         timetypes.NewDurationValueFromStringMust("1h1m1s"),
			expectedDuration: "1h1m1s",
		},
	}
	for name, testCase := range testCases {
		name, testCase := name, testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			duration, diags := testCase.Duration.ValueDuration()
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
