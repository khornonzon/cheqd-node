package cli

import (
	"context"
	"github.com/cheqd/cheqd-node/x/resource/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

func CmdGetAllResourceVersions() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "all-resource-versions [collectionId] [name] [resourceType] [mimeType]",
		Short: "Query all resource versions",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			collectionId := args[0]
			name := args[1]
			resourceType := args[2]
			mimeType := args[3]

			params := &types.QueryGetAllResourceVersionsRequest{
				CollectionId: collectionId,
				Name:         name,
				ResourceType: resourceType,
				MimeType:     mimeType,
			}

			resp, err := queryClient.AllResourceVersions(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(resp)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
