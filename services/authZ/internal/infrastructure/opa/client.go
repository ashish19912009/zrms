package opa

import (
	"context"
	"errors"
	"fmt"

	"github.com/open-policy-agent/opa/rego"
)

type Client struct {
	query rego.PreparedEvalQuery
}

func NewClient(policyPath string) (*Client, error) {
	query, err := rego.New(
		rego.Query("data.authz.allow"),
		rego.Load([]string{policyPath}, nil),
	).PrepareForEval(context.Background())

	if err != nil {
		return nil, fmt.Errorf("failed to prepare OPA query: %w", err)
	}

	return &Client{query: query}, nil
}

func (c *Client) Evaluate(ctx context.Context, input map[string]interface{}) (bool, error) {
	results, err := c.query.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		return false, fmt.Errorf("failed to evaluate policy: %w", err)
	}

	if len(results) == 0 {
		return false, errors.New("no results from OPA evaluation")
	}

	allowed, ok := results[0].Expressions[0].Value.(bool)
	if !ok {
		return false, errors.New("invalid result type from OPA")
	}

	return allowed, nil
}
