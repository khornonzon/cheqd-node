package keeper

import (
	"crypto/ed25519"
	"errors"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"github.com/cheqd/cheqd-node/x/cheqd/types"
)

func FindPublicKey(signer types.Signer, id string) (ed25519.PublicKey, error) {
	for _, authentication := range signer.Authentication {
		if authentication == id {
			for _, vm := range signer.VerificationMethod {
				if vm.Id == id {
					return base58.Decode(vm.PublicKeyMultibase[1:]), nil
				}
			}

			msg := fmt.Sprintf("Verification Method %s not found", id)
			return nil, errors.New(msg)
		}
	}

	msg := fmt.Sprintf("Authentication %s not found", id)
	return nil, errors.New(msg)
}

func FindVerificationMethod(vms []*types.VerificationMethod, id string) *types.VerificationMethod {
	for _, vm := range vms {
		if vm.Id == id {
			return vm
		}
	}

	return nil
}
