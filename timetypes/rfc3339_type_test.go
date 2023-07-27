package timetypes_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
)

func TestRFC3339TypeValidate(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in            tftypes.Value
		expectedDiags diag.Diagnostics
	}{
		"empty-struct": {
			in: tftypes.Value{},
		},
		"null": {
			in: tftypes.NewValue(tftypes.String, nil),
		},
		"unknown": {
			in: tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
		},
		"valid RFC3339": {
			in: tftypes.NewValue(tftypes.String, "2023-07-25T20:43:16+00:00"),
		},
		"valid RFC3339 - Zulu": {
			in: tftypes.NewValue(tftypes.String, "2023-07-25T20:43:16Z"),
		},
		"valid RFC3339 - UTC Offset": {
			in: tftypes.NewValue(tftypes.String, "2023-07-25T20:43:16-05:00"),
		},
		"invalid RFC3339 - no date": {
			in: tftypes.NewValue(tftypes.String, "20:43:16-05:00"),
			expectedDiags: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid RFC3339 String Value",
					"A string value was provided that is not valid RFC3339 string format.\n\n"+
						"Given Value: 20:43:16-05:00\n"+
						"Error: parsing time \"20:43:16-05:00\" as \"2006-01-02T15:04:05Z07:00\": "+
						"cannot parse \"20:43:16-05:00\" as \"2006\"",
				),
			},
		},
		"invalid RFC3339 - no time": {
			in: tftypes.NewValue(tftypes.String, "2023-07-25T"),
			expectedDiags: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid RFC3339 String Value",
					"A string value was provided that is not valid RFC3339 string format.\n\n"+
						"Given Value: 2023-07-25T\n"+
						"Error: parsing time \"2023-07-25T\" as \"2006-01-02T15:04:05Z07:00\": "+
						"cannot parse \"\" as \"15\"",
				),
			},
		},
		"invalid RFC3339 - normal string": {
			in: tftypes.NewValue(tftypes.String, "notvalidrfc3339"),
			expectedDiags: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid RFC3339 String Value",
					"A string value was provided that is not valid RFC3339 string format.\n\n"+
						"Given Value: notvalidrfc3339\n"+
						"Error: parsing time \"notvalidrfc3339\" as \"2006-01-02T15:04:05Z07:00\": "+
						"cannot parse \"notvalidrfc3339\" as \"2006\"",
				),
			},
		},
		"wrong-value-type": {
			in: tftypes.NewValue(tftypes.Number, 123),
			expectedDiags: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"RFC3339 Time Type Validation Error",
					"An unexpected error was encountered trying to validate an attribute value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+
						"expected String value, received tftypes.Value with value: tftypes.Number<\"123\">",
				),
			},
		},
	}
	for name, testCase := range testCases {
		name, testCase := name, testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			diags := timetypes.RFC3339Type{}.Validate(context.Background(), testCase.in, path.Root("test"))
			println(diags.Errors())

			if diff := cmp.Diff(diags, testCase.expectedDiags); diff != "" {
				t.Errorf("Unexpected diagnostics (-got, +expected): %s", diff)
			}
		})
	}
}

func TestRFC3339TypeValueFromTerraform(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in          tftypes.Value
		expectation attr.Value
		expectedErr string
	}{
		"true": {
			in:          tftypes.NewValue(tftypes.String, "2023-07-25T20:43:16+00:00"),
			expectation: timetypes.NewRFC3339Value("2023-07-25T20:43:16+00:00"),
		},
		"unknown": {
			in:          tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
			expectation: timetypes.NewRFC3339Unknown(),
		},
		"null": {
			in:          tftypes.NewValue(tftypes.String, nil),
			expectation: timetypes.NewRFC3339Null(),
		},
		"wrongType": {
			in:          tftypes.NewValue(tftypes.Number, 123),
			expectedErr: "can't unmarshal tftypes.Number into *string, expected string",
		},
	}
	for name, testCase := range testCases {
		name, testCase := name, testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			got, err := timetypes.RFC3339Type{}.ValueFromTerraform(ctx, testCase.in)
			if err != nil {
				if testCase.expectedErr == "" {
					t.Fatalf("Unexpected error: %s", err)
				}
				if testCase.expectedErr != err.Error() {
					t.Fatalf("Expected error to be %q, got %q", testCase.expectedErr, err.Error())
				}
				return
			}
			if err == nil && testCase.expectedErr != "" {
				t.Fatalf("Expected error to be %q, didn't get an error", testCase.expectedErr)
			}
			if !got.Equal(testCase.expectation) {
				t.Errorf("Expected %+v, got %+v", testCase.expectation, got)
			}
			if testCase.expectation.IsNull() != testCase.in.IsNull() {
				t.Errorf("Expected null-ness match: expected %t, got %t", testCase.expectation.IsNull(), testCase.in.IsNull())
			}
			if testCase.expectation.IsUnknown() != !testCase.in.IsKnown() {
				t.Errorf("Expected unknown-ness match: expected %t, got %t", testCase.expectation.IsUnknown(), !testCase.in.IsKnown())
			}
		})
	}
}
