package main

import (
	"context"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/urfave/cli/v3"

	"github.com/nsec/askgod/api"
)

func (c *client) cmdStatus(ctx context.Context, _ *cli.Command) error {
	// Get the data
	resp := api.Status{}

	err := c.queryStruct(ctx, "GET", "", nil, &resp)
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(&resp)
	if err != nil {
		return err
	}

	_, _ = fmt.Printf("%s", data) //nolint:forbidigo

	return nil
}
