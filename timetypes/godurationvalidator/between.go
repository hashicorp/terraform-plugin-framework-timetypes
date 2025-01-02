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

var _ validator.String = betweenValidator{}

// betweenValidator validates that a GoDuration Attribute's value is in a range.
type betweenValidator struct {
	minimum, maximum time.Duration
}

// Description describes the validation in plain text formatting.
func (validator betweenValidator) Description(_ context.Context) string {
	return fmt.Sprintf("value must be between %s and %s", validator.minimum, validator.maximum)
}

// MarkdownDescription describes the validation in Markdown formatting.
func (validator betweenValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

// ValidateString performs the validation.
func (v betweenValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	if d, err := time.ParseDuration(request.ConfigValue.ValueString()); err == nil && (d < v.minimum || d > v.maximum) {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			request.Path,
			v.Description(ctx),
			d.String(),
		))
	}
}

// Between returns an AttributeValidator which ensures that any configured
// attribute value:
//
//   - Is a duration
//   - Is greater than or equal to the given minimum and less than or equal to the given maximum.
//
// Null (unconfigured) and unknown (known after apply) values are skipped.
func Between(minimum, maximum time.Duration) validator.String {
	if minimum > maximum {
		return nil
	}

	return betweenValidator{
		minimum: minimum,
		maximum: maximum,
	}
}
