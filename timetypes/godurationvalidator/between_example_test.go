// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package godurationvalidator_test

import (
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes/godurationvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func ExampleBetween() {
	// Used within a Schema method of a DataSource, Provider, or Resource
	_ = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"example_attr": schema.StringAttribute{
				Required:   true,
				CustomType: timetypes.GoDurationType{},
				Validators: []validator.String{
					// Validate duration value must be at least 5 seconds and at most 60 seconds
					godurationvalidator.Between(5*time.Second, 60*time.Second),
				},
			},
		},
	}
}
