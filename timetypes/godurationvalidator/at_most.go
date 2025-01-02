// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package godurationvalidator

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = atMostValidator{}

// atMostValidator validates that a GoDuration Attribute's value is at most a certain value.
type atMostValidator struct {
	maximum time.Duration
}

// Description describes the validation in plain text formatting.
func (validator atMostValidator) Description(_ context.Context) string {
	return fmt.Sprintf("value must be at most %s", validator.maximum)
}

// MarkdownDescription describes the validation in Markdown formatting.
func (validator atMostValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

// ValidateString performs the validation.
func (v atMostValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	if d, err := time.ParseDuration(request.ConfigValue.ValueString()); err == nil && d > v.maximum {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			request.Path,
			v.Description(ctx),
			d.String(),
		))
	}
}

// AtMost returns an AttributeValidator which ensures that any configured
// attribute value:
//
//   - Is a duration
//   - Is less than or equal to the given maximum.
//
// Null (unconfigured) and unknown (known after apply) values are skipped.
func AtMost(maximum time.Duration) validator.String {
	return atMostValidator{
		maximum: maximum,
	}
}
