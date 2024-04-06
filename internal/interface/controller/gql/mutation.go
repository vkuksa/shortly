package gql

import (
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/vkuksa/shortly/internal/link"
)

func initMutationType(uc *link.UseCase) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"shorten": &graphql.Field{
				Type:        linkModel,
				Description: "Shorten link",
				Args: graphql.FieldConfigArgument{
					"url": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					url := p.Args["url"].(string)
					link, err := uc.Shorten(p.Context, url)
					if err != nil {
						return nil, fmt.Errorf("shorten: %w", err)
					}

					return link, nil
				},
			},
		},
	})
}
