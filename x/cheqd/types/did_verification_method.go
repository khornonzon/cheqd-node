package types

import (
	"crypto/ed25519"
	"errors"
	"github.com/cheqd/cheqd-node/x/cheqd/utils"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/multiformats/go-multibase"
)

const (
	JsonWebKey2020             = "JsonWebKey2020"
	Ed25519VerificationKey2020 = "Ed25519VerificationKey2020"
)

var SupportedMethodTypes = []string{
	JsonWebKey2020,
	Ed25519VerificationKey2020,
}

var JwkMethodTypes = []string{
	JsonWebKey2020,
}

var MultibaseMethodTypes = []string{
	Ed25519VerificationKey2020,
}

func NewVerificationMethod(id string, type_ string, controller string, publicKeyJwk []*KeyValuePair, publicKeyMultibase string) *VerificationMethod {
	return &VerificationMethod{
		Id:                 id,
		Type:               type_,
		Controller:         controller,
		PublicKeyJwk:       publicKeyJwk,
		PublicKeyMultibase: publicKeyMultibase,
	}
}


// Helpers

func FindVerificationMethod(vms []VerificationMethod, id string) (VerificationMethod, bool) {
	for _, vm := range vms {
		if vm.Id == id {
			return vm, true
		}
	}

	return VerificationMethod{}, false
}

func GetVerificationMethodIds(vms []*VerificationMethod) []string {
	res := make([]string, len(vms))

	for i := range vms {
		res[i] = vms[i].Id
	}

	return res
}

func VerifySignature(vm VerificationMethod, message []byte, signature []byte) error {
	switch vm.Type {
	case Ed25519VerificationKey2020:
		_, keyBytes, err := multibase.Decode(vm.PublicKeyMultibase)
		if err != nil {
			return err
		}

		valid := ed25519.Verify(keyBytes, message, signature)
		if !valid {
			return ErrInvalidSignature.Wrapf("verification method: %s", vm.Id)
		}

		return nil

	case JsonWebKey2020:
		// TODO: Implement
		return nil

	default:
		return ErrSignatureNotSupported.Wrapf("verification method: %s, signature type: %s", vm.Id, vm.Type)
	}
}

// Validation

func (vm VerificationMethod) Validate(baseDid string, allowedNamespaces []string) error {
	return validation.ValidateStruct(&vm,
		validation.Field(&vm.Id, validation.Required, IsDIDUrl(allowedNamespaces, Empty, Empty, Required), HasPrefix(baseDid)),
		validation.Field(&vm.Controller, validation.Required, IsDID(allowedNamespaces)),
		validation.Field(&vm.Type, validation.Required, validation.In(utils.ToInterfaces(SupportedMethodTypes)...)),
		validation.Field(&vm.PublicKeyJwk,
			validation.When(utils.Contains(JwkMethodTypes, vm.Type), validation.Required).Else(validation.Empty),
		),
		validation.Field(&vm.PublicKeyMultibase,
			validation.When(utils.Contains(MultibaseMethodTypes, vm.Type), validation.Required, IsMultibase()).Else(validation.Empty),
		),
	)
}

func ValidVerificationMethod(baseDid string, allowedNamespaces []string) *CustomErrorRule {
	return NewCustomErrorRule(func(value interface{}) error {
		casted, ok := value.(VerificationMethod)
		if !ok {
			panic("ValidVerificationMethod must be only applied on verification methods")
		}

		return casted.Validate(baseDid, allowedNamespaces)
	})
}

func IsUniqueVerificationMethodList() *CustomErrorRule {
	return NewCustomErrorRule(func(value interface{}) error {
		casted, ok := value.([]*VerificationMethod)
		if !ok {
			panic("IsUniqueVerificationMethodList must be only applied on VM lists")
		}

		ids := GetVerificationMethodIds(casted)
		if !utils.IsUnique(ids) {
			return errors.New("there are verification method duplicates")
		}

		return nil
	})
}
