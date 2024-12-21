// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package godurationvalidator_test

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes/godurationvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func TestAtMostValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val         timetypes.GoDuration
		maximum     time.Duration
		expectError bool
	}
	tests := map[string]testCase{
		"unknown GoDuration": {
			val:     timetypes.NewGoDurationUnknown(),
			maximum: 2 * time.Second,
		},
		"null GoDuration": {
			val:     timetypes.NewGoDurationNull(),
			maximum: 2 * time.Second,
		},
		"valid duration as GoDuration": {
			val:     timetypes.NewGoDurationValue(1 * time.Second),
			maximum: 2 * time.Second,
		},
		"valid duration as GoDuration min": {
			val:     timetypes.NewGoDurationValue(2 * time.Second),
			maximum: 2 * time.Second,
		},
		"too large duration as GoDuration": {
			val:         timetypes.NewGoDurationValue(4 * time.Second),
			maximum:     2 * time.Second,
			expectError: true,
		},
	}

	for name, test := range tests {
		value, _ := test.val.ToStringValue(context.TODO())
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			request := validator.StringRequest{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    value,
			}
			response := validator.StringResponse{}
			godurationvalidator.AtMost(test.maximum).ValidateString(context.TODO(), request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})
	}
}
