package gql

import "github.com/graphql-go/graphql"

var linkModel = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Product",
		Fields: graphql.Fields{
			"uuid": &graphql.Field{
				Type: graphql.String,
			},
			"url": &graphql.Field{
				Type: graphql.String,
			},
			"count": &graphql.Field{
				Type: graphql.Int,
			},
			"createdAt": &graphql.Field{
				Type: graphql.DateTime,
			},
			"expiresAt": &graphql.Field{
				Type: graphql.DateTime,
			},
		},
	},
)
