package keeper

import (
	"context"
	"crypto/sha256"
	"time"

	cheqdkeeper "github.com/cheqd/cheqd-node/x/cheqd/keeper"
	cheqdtypes "github.com/cheqd/cheqd-node/x/cheqd/types"
	cheqdutils "github.com/cheqd/cheqd-node/x/cheqd/utils"
	"github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CreateResource(goCtx context.Context, msg *types.MsgCreateResource) (*types.MsgCreateResourceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate corresponding DIDDoc exists
	namespace := k.cheqdKeeper.GetDidNamespace(ctx)
	did := cheqdutils.JoinDID(cheqdtypes.DidMethod, namespace, msg.Payload.CollectionId)

	// Validate Resource doesn't exist
	if k.HasResource(&ctx, msg.Payload.CollectionId, msg.Payload.Id) {
		return nil, types.ErrResourceExists.Wrap(msg.Payload.Id)
	}

	// Validate signatures
	didDocStateValue, err := k.cheqdKeeper.GetDid(&ctx, did)
	didDoc, err := didDocStateValue.UnpackDataAsDid()
	if err != nil {
		return nil, err
	}

	// We can use the same signers as for DID creation because didDoc stays the same
	signers := cheqdkeeper.GetSignerDIDsForDIDCreation(*didDoc)
	err = cheqdkeeper.VerifyAllSignersHaveAllValidSignatures(&k.cheqdKeeper, &ctx, map[string]cheqdtypes.StateValue{},
		msg.Payload.GetSignBytes(), signers, msg.Signatures)
	if err != nil {
		return nil, err
	}

	// Build Resource
	resource := msg.Payload.ToResource()

	resource.Checksum = string(sha256.New().Sum(resource.Data))
	resource.Created = time.Now().UTC().Format(time.RFC3339)
	// TODO: set version + update forward and backward links

	// Append backlink to didDoc
	didDocStateValue.Metadata.Resources = append(didDocStateValue.Metadata.Resources, resource.Id)
	err = k.cheqdKeeper.SetDid(&ctx, &didDocStateValue)
	if err != nil {
		return nil, types.ErrInternal.Wrapf(err.Error())
	}

	// Persist resource
	err = k.SetResource(&ctx, &resource)
	if err != nil {
		return nil, types.ErrInternal.Wrapf(err.Error())
	}

	// Build and return response
	return &types.MsgCreateResourceResponse{
		Resource: &resource,
	}, nil
}
