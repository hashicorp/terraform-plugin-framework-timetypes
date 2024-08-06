## 0.5.0 (August 06, 2024)

ENHANCEMENTS:

* timetypes: Implement StringSemanticEquals for GoDuration ([#75](https://github.com/hashicorp/terraform-plugin-framework-timetypes/issues/75))

## 0.4.0 (June 04, 2024)

BREAKING CHANGES:

* timetypes: Removed `Validate()` method from `RFC3339Type` type following deprecation of `xattr.TypeWithValidate` ([#58](https://github.com/hashicorp/terraform-plugin-framework-timetypes/issues/58))

NOTES:

* all: This Go module has been updated to Go 1.21 per the [Go support policy](https://go.dev/doc/devel/release#policy). It is recommended to review the [Go 1.21 release notes](https://go.dev/doc/go1.21) before upgrading. Any consumers building on earlier Go versions may experience errors ([#49](https://github.com/hashicorp/terraform-plugin-framework-timetypes/issues/49))

FEATURES:

* Support Go standard library time duration type (`GoDurationType` and `GoDuration`) ([#66](https://github.com/hashicorp/terraform-plugin-framework-timetypes/issues/66))

ENHANCEMENTS:

* timetypes: Added `ValidateAttribute()` method to `RFC3339` type, which supports validating an attribute value ([#58](https://github.com/hashicorp/terraform-plugin-framework-timetypes/issues/58))
* timetypes: Added `ValidateParameter()` method to `RFC3339` type, which supports validating a provider-defined function parameter value ([#58](https://github.com/hashicorp/terraform-plugin-framework-timetypes/issues/58))

## 0.3.0 (October 09, 2023)

ENHANCEMENTS:

* timetypes: Added `NewRFC3339TimePointerValue()` function, which supports creating a known value from a `*time.Time` ([#20](https://github.com/hashicorp/terraform-plugin-framework-timetypes/issues/20))

## 0.2.0 (August 08, 2023)

BREAKING CHANGES:

* timetypes: The `NewRFC3339Value` and `NewRFC3339PointerValue` functions now return `diag.Diagnostics` so an error diagnostic can be raised if the string is not RFC3339 formatted ([#7](https://github.com/hashicorp/terraform-plugin-framework-timetypes/issues/7))

ENHANCEMENTS:

* timetypes: Added `NewRFC3339ValueMust` and `NewRFC3339PointerValueMust` value creation functions, which panic if the string is not RFC3339 formatted ([#7](https://github.com/hashicorp/terraform-plugin-framework-timetypes/issues/7))
* timetypes: Added `NewRFC3339TimeValue()` function, which supports creating a known value from a `time.Time` ([#6](https://github.com/hashicorp/terraform-plugin-framework-timetypes/issues/6))

## 0.1.0 (July 28, 2023)

FEATURES:

* timetypes: Add new RFC3339 custom type implementation, representing an RFC 3339 timestamp string ([#2](https://github.com/hashicorp/terraform-plugin-framework-timetypes/issues/2))

