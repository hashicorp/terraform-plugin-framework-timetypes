// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package timetypes

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/attr/xattr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var (
	_ basetypes.StringValuable       = (*Duration)(nil)
	_ xattr.ValidateableAttribute    = (*Duration)(nil)
	_ function.ValidateableParameter = (*Duration)(nil)
)

// Duration represents a valid Go time duration string.
// See https://pkg.go.dev/time#ParseDuration for more details
type Duration struct {
	basetypes.StringValue
}

// Type returns an RFC3339Type.
func (d Duration) Type(_ context.Context) attr.Type {
	return DurationType{}
}

// Equal returns true if the given value is equivalent.
func (d Duration) Equal(o attr.Value) bool {
	other, ok := o.(Duration)

	if !ok {
		return false
	}

	// Strings are already validated at this point, ignoring errors
	thisDuration, _ := time.ParseDuration(d.ValueString())
	otherDuration, _ := time.ParseDuration(other.ValueString())

	return thisDuration == otherDuration
}

// ValidateAttribute implements attribute value validation. This type requires the value to be a String value that
// is valid Go time duration and utilizes the Go `time` library
func (d Duration) ValidateAttribute(ctx context.Context, req xattr.ValidateAttributeRequest, resp *xattr.ValidateAttributeResponse) {
	if d.IsUnknown() || d.IsNull() {
		return
	}

	if _, err := time.ParseDuration(d.ValueString()); err != nil {
		resp.Diagnostics.Append(diag.WithPath(req.Path, durationInvalidStringDiagnostic(d.ValueString(), err)))

		return
	}
}

// ValidateParameter implements provider-defined function parameter value validation. This type requires the value to
// be a String value that is a valid time duration and utilizes the Go `time` library
func (d Duration) ValidateParameter(ctx context.Context, req function.ValidateParameterRequest, resp *function.ValidateParameterResponse) {
	if d.IsUnknown() || d.IsNull() {
		return
	}

	if _, err := time.ParseDuration(d.ValueString()); err != nil {
		resp.Error = function.NewArgumentFuncError(
			req.Position,
			"Invalid time duration String Value: "+
				"A string value was provided that is not valid time duration string format.\n\n"+
				"Given Value: "+d.ValueString()+"\n"+
				"Error: "+err.Error(),
		)

		return
	}
}

// ValueDuration creates a new time.Duration instance with the time duration StringValue. A null or unknown value will produce an error diagnostic.
func (d Duration) ValueDuration() (time.Duration, diag.Diagnostics) {
	var diags diag.Diagnostics

	if d.IsNull() {
		diags.Append(diag.NewErrorDiagnostic("Duration ValueDuration Error", "Duration string value is null"))
		return time.Duration(0), diags
	}

	if d.IsUnknown() {
		diags.Append(diag.NewErrorDiagnostic("Duration ValueDuration Error", "Duration string value is unknown"))
		return time.Duration(0), diags
	}

	duration, err := time.ParseDuration(d.ValueString())
	if err != nil {
		diags.Append(diag.NewErrorDiagnostic("Duration ValueDuration Error", err.Error()))
		return time.Duration(0), diags
	}

	return duration, nil
}

// NewDurationNull creates an Duration with a null value. Determine whether the value is null via IsNull method.
func NewDurationNull() Duration {
	return Duration{
		StringValue: basetypes.NewStringNull(),
	}
}

// NewDurationUnknown creates an Duration with an unknown value. Determine whether the value is unknown via IsUnknown method.
func NewDurationUnknown() Duration {
	return Duration{
		StringValue: basetypes.NewStringUnknown(),
	}
}

// NewDurationValue creates an Duration with a known value.
func NewDurationValue(value time.Duration) Duration {
	return Duration{
		StringValue: basetypes.NewStringValue(value.String()),
	}
}

// NewDurationPointerValue creates an Duration with a null value if nil or
// a known value.
func NewDurationPointerValue(value *time.Duration) Duration {
	if value == nil {
		return NewDurationNull()
	}

	return Duration{
		StringValue: basetypes.NewStringValue(value.String()),
	}
}

// NewDurationValueFromString creates an Duration with a known value or raises an error
// diagnostic if the string is not Duration format.
func NewDurationValueFromString(value string) (Duration, diag.Diagnostics) {
	_, err := time.ParseDuration(value)

	if err != nil {
		// Returning an unknown value will guarantee that, as a last resort,
		// Terraform will return an error if attempting to store into state.
		return NewDurationUnknown(), diag.Diagnostics{durationInvalidStringDiagnostic(value, err)}
	}

	return Duration{
		StringValue: basetypes.NewStringValue(value),
	}, nil
}

// NewDurationValueFromStringMust creates an Duration with a known value or raises a panic
// if the string is not Duration format.
//
// This creation function is only recommended to create Duration values which
// either will not potentially affect practitioners, such as testing, or within
// exhaustively tested provider logic.
func NewDurationValueFromStringMust(value string) Duration {
	_, err := time.ParseDuration(value)

	if err != nil {
		panic(fmt.Sprintf("Invalid Duration String Value (%s): %s", value, err))
	}

	return Duration{
		StringValue: basetypes.NewStringValue(value),
	}
}

// NewDurationValueFromPointerString creates an Duration with a null value if nil, a known
// value, or raises an error diagnostic if the string is not Duration format.
func NewDurationValueFromPointerString(value *string) (Duration, diag.Diagnostics) {
	if value == nil {
		return NewDurationNull(), nil
	}

	return NewDurationValueFromString(*value)
}

// NewDurationValueFromPointerStringMust creates an Duration with a null value if nil, a
// known value, or raises a panic if the string is not Duration format.
//
// This creation function is only recommended to create Duration values which
// either will not potentially affect practitioners, such as testing, or within
// exhaustively tested provider logic.
func NewDurationValueFromPointerStringMust(value *string) Duration {
	if value == nil {
		return NewDurationNull()
	}

	return NewDurationValueFromStringMust(*value)
}
