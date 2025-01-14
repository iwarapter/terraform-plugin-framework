package basetypes

import (
	"context"
	"fmt"
	"math/big"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/attr/xattr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Float64Typable extends attr.Type for float64 types.
// Implement this interface to create a custom Float64Type type.
type Float64Typable interface {
	xattr.TypeWithValidate

	// ValueFromFloat64 should convert the Float64 to a Float64Valuable type.
	ValueFromFloat64(context.Context, Float64Value) (Float64Valuable, diag.Diagnostics)
}

var _ Float64Typable = Float64Type{}

// Float64Type is the base framework type for a floating point number.
// Float64Value is the associated value type.
type Float64Type struct{}

// ApplyTerraform5AttributePathStep applies the given AttributePathStep to the
// type.
func (t Float64Type) ApplyTerraform5AttributePathStep(step tftypes.AttributePathStep) (interface{}, error) {
	return nil, fmt.Errorf("cannot apply AttributePathStep %T to %s", step, t.String())
}

// Equal returns true if the given type is equivalent.
func (t Float64Type) Equal(o attr.Type) bool {
	_, ok := o.(Float64Type)

	return ok
}

// String returns a human readable string of the type name.
func (t Float64Type) String() string {
	return "basetypes.Float64Type"
}

// TerraformType returns the tftypes.Type that should be used to represent this
// framework type.
func (t Float64Type) TerraformType(_ context.Context) tftypes.Type {
	return tftypes.Number
}

// Validate implements type validation.
func (t Float64Type) Validate(ctx context.Context, in tftypes.Value, path path.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	if in.Type() == nil {
		return diags
	}

	if !in.Type().Equal(tftypes.Number) {
		diags.AddAttributeError(
			path,
			"Float64 Type Validation Error",
			"An unexpected error was encountered trying to validate an attribute value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+
				fmt.Sprintf("Expected Number value, received %T with value: %v", in, in),
		)
		return diags
	}

	if !in.IsKnown() || in.IsNull() {
		return diags
	}

	var value *big.Float
	err := in.As(&value)

	if err != nil {
		diags.AddAttributeError(
			path,
			"Float64 Type Validation Error",
			"An unexpected error was encountered trying to validate an attribute value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+
				fmt.Sprintf("Cannot convert value to big.Float: %s", err),
		)
		return diags
	}

	_, accuracy := value.Float64()

	if accuracy != 0 {
		diags.AddAttributeError(
			path,
			"Float64 Type Validation Error",
			fmt.Sprintf("Value %s cannot be represented as a 64-bit floating point.", value),
		)
		return diags
	}

	return diags
}

// ValueFromFloat64 returns a Float64Valuable type given a Float64Value.
func (t Float64Type) ValueFromFloat64(_ context.Context, v Float64Value) (Float64Valuable, diag.Diagnostics) {
	return v, nil
}

// ValueFromTerraform returns a Value given a tftypes.Value.  This is meant to
// convert the tftypes.Value into a more convenient Go type for the provider to
// consume the data with.
func (t Float64Type) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if !in.IsKnown() {
		return NewFloat64Unknown(), nil
	}

	if in.IsNull() {
		return NewFloat64Null(), nil
	}

	var bigF *big.Float
	err := in.As(&bigF)

	if err != nil {
		return nil, err
	}

	f, accuracy := bigF.Float64()

	if accuracy != 0 {
		return nil, fmt.Errorf("Value %s cannot be represented as a 64-bit floating point.", bigF)
	}

	return NewFloat64Value(f), nil
}

// ValueType returns the Value type.
func (t Float64Type) ValueType(_ context.Context) attr.Value {
	// This Value does not need to be valid.
	return Float64Value{}
}
