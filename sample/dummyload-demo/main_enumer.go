// GENERATED CODE -- DO NOT EDIT

package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/maargenton/go-cli/pkg/enumer/enum"
)

// ---------------------------------------------------------------------------
// WorkloadType

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Rsa2048-(0)]
	_ = x[Rsa4096-(1)]
	_ = x[EcdsaP256-(2)]
	_ = x[EcdsaP284-(3)]
	_ = x[EcdsaP521-(4)]
	_ = x[Ed25519-(5)]
}

var _ enum.Type = (*WorkloadType)(nil)

var WorkloadTypeValues = []enum.Value{
	{
		Name:     "rsa2048",
		GoName:   "Rsa2048",
		AltNames: []string{"rsa2048", "Rsa2048"},
		Value:    Rsa2048,
	},
	{
		Name:     "rsa4096",
		GoName:   "Rsa4096",
		AltNames: []string{"rsa4096", "Rsa4096"},
		Value:    Rsa4096,
	},
	{
		Name:     "ecdsa-p256",
		GoName:   "EcdsaP256",
		AltNames: []string{"ecdsa-p256", "EcdsaP256", "ecdsaP256", "ecdsa_p256"},
		Value:    EcdsaP256,
	},
	{
		Name:     "ecdsa-p284",
		GoName:   "EcdsaP284",
		AltNames: []string{"ecdsa-p284", "EcdsaP284", "ecdsaP284", "ecdsa_p284"},
		Value:    EcdsaP284,
	},
	{
		Name:     "ecdsa-p521",
		GoName:   "EcdsaP521",
		AltNames: []string{"ecdsa-p521", "EcdsaP521", "ecdsaP521", "ecdsa_p521"},
		Value:    EcdsaP521,
	},
	{
		Name:     "ed25519",
		GoName:   "Ed25519",
		AltNames: []string{"ed25519", "Ed25519"},
		Value:    Ed25519,
	},
}

func (v WorkloadType) EnumValues() []enum.Value {
	return WorkloadTypeValues
}

func (v WorkloadType) String() string {
	switch v {
	case Rsa2048:
		return "rsa2048"
	case Rsa4096:
		return "rsa4096"
	case EcdsaP256:
		return "ecdsa-p256"
	case EcdsaP284:
		return "ecdsa-p284"
	case EcdsaP521:
		return "ecdsa-p521"
	case Ed25519:
		return "ed25519"
	}
	return "WorkloadType(" + strconv.FormatInt(int64(v), 10) + ")"
}

func ParseWorkloadType(s string) (WorkloadType, error) {
	switch strings.ToLower(s) {
	case "rsa2048":
		return Rsa2048, nil
	case "rsa4096":
		return Rsa4096, nil
	case "ecdsa-p256", "ecdsap256", "ecdsa_p256":
		return EcdsaP256, nil
	case "ecdsa-p284", "ecdsap284", "ecdsa_p284":
		return EcdsaP284, nil
	case "ecdsa-p521", "ecdsap521", "ecdsa_p521":
		return EcdsaP521, nil
	case "ed25519":
		return Ed25519, nil
	}
	return 0, fmt.Errorf("invalid WorkloadType value '%v'", s)
}

func (v *WorkloadType) Set(s string) error {
	vv, err := ParseWorkloadType(s)
	if err != nil {
		return err
	}
	*v = vv
	return nil
}

func (v WorkloadType) MarshalText() (text []byte, err error) {
	return []byte(v.String()), nil
}

func (v *WorkloadType) UnmarshalText(text []byte) error {
	return v.Set(string(text))
}

// WorkloadType
// ---------------------------------------------------------------------------
