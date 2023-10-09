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

