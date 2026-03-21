package main

import (
	"context"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/urfave/cli/v3"

	"github.com/nsec/askgod/api"
)

func (c *client) cmdAdminConfig(ctx context.Context, cmd *cli.Command) error {
	// Get the data
	resp := api.Config{}

	err := c.queryStruct(ctx, "GET", "/config", nil, &resp)
	if err != nil {
		return err
	}

	// Process any field update
	if cmd.NArg() > 0 {
		for _, arg := range cmd.Args().Slice() {
			err := setStructKey(&resp, arg)
			if err != nil {
				return err
			}
		}

		// Update the team
		err = c.queryStruct(ctx, "PUT", "/config", resp.ConfigPut, nil)
		if err != nil {
			return err
		}

		return nil
	}

	data, err := yaml.Marshal(&resp)
	if err != nil {
		return err
	}

	_, _ = fmt.Printf("%s", data) //nolint:forbidigo

	return nil
}
