package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Here we need to make all the links the
func (app *App) Migration07(ctx sdk.Context) {
	// TODO: We have to reset links first. Then we can use GetLastResourceVersionHeader
	// but only because resources in state are corted by creation time.
	// Also, we need to avoid loading all resources in memory.
	resourceList := app.resourceKeeper.GetAllResources(&ctx)
	for _, resource := range resourceList {
		previousResourceVersionHeader, found := app.resourceKeeper.GetLastResourceVersionHeader(&ctx, resource.Header.CollectionId, resource.Header.Name, resource.Header.ResourceType)
		if found {
			// Set links
			previousResourceVersionHeader.NextVersionId = resource.Header.Id
			resource.Header.PreviousVersionId = previousResourceVersionHeader.Id

			// Update previous version
			err := app.resourceKeeper.UpdateResourceHeader(&ctx, &previousResourceVersionHeader)
			if err != nil {
				return
			}
		}
	}
}
