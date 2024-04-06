package gql

import (
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/vkuksa/shortly/internal/link"
)

func initQueryType(uc *link.UseCase) *graphql.Object {
	return graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"link": &graphql.Field{
					Type:        linkModel,
					Description: "Get link by uuid",
					Args: graphql.FieldConfigArgument{
						"uuid": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						uuid, ok := p.Args["uuid"].(string)
						if ok {
							link, err := uc.Retrieve(p.Context, uuid)
							if err != nil {
								return nil, fmt.Errorf("retrieve: %w", err)
							}

							return link, nil
						}
						return nil, nil
					},
				},
			},
		})
}
