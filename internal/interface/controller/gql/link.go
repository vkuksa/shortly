package gql

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/graphql-go/graphql"
	"github.com/vkuksa/shortly/internal/link"
)

type ErrorHandler interface {
	HandleGraphQLError(w http.ResponseWriter, r *http.Request, err error)
}

type LinkController struct {
	uc         *link.UseCase
	errhandler ErrorHandler

	schema graphql.Schema
}

func NewLinkController(uc *link.UseCase, eh ErrorHandler) (*LinkController, error) {
	var schema, err = graphql.NewSchema(
		graphql.SchemaConfig{
			Query:    initQueryType(uc),
			Mutation: initMutationType(uc),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("new schema: %w", err)
	}

	return &LinkController{uc: uc, errhandler: eh, schema: schema}, nil
}

func (c *LinkController) Register(router chi.Router) {
	router.Get("/graphql", c.gqlHandler)
}

func (c *LinkController) gqlHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if query == "" {
		c.errhandler.HandleGraphQLError(w, r, fmt.Errorf("empty query provided: %w", link.ErrBadInput))
		return
	}

	result, err := c.executeQuery(query)
	if err != nil {
		c.errhandler.HandleGraphQLError(w, r, fmt.Errorf("execute query %s failed: %w", query, err))
		return
	}

	c.writeJSONResponse(w, result, http.StatusOK)
}

func (c *LinkController) executeQuery(query string) (*graphql.Result, error) {
	result := graphql.Do(graphql.Params{
		Schema:        c.schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("graphql do: %v", result.Errors)
	}

	return result, nil
}

func (c *LinkController) writeJSONResponse(w http.ResponseWriter, obj any, status int) {
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(obj)
	if err != nil {
		slog.Error("writing response failed", slog.Any("component", "gql"), slog.Any("error", err))
		return
	}
}
