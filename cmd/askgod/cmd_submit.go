package main

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"

	"github.com/nsec/askgod/api"
)

func (c *client) cmdSubmit(ctx context.Context, cmd *cli.Command) error {
	if cmd.NArg() != 1 {
		_ = cli.ShowCommandHelp(ctx, cmd, "submit")

		return nil
	}

	// Prepare the input
	flag := api.FlagPost{}
	flag.Flag = cmd.Args().Get(0)
	flag.Notes = cmd.String("notes")

	// Send the flag
	resp := api.Flag{}

	err := c.queryStruct(ctx, "POST", "/team/flags", flag, &resp)
	if err != nil {
		return err
	}

	// Process the points
	switch {
	case resp.Value < 0:
		_, _ = fmt.Printf("You shouldn't have sent that! You just lost your team %d points.\n", resp.Value*-1) //nolint:forbidigo

	case resp.Value == 0:
		_, _ = fmt.Print("You sent a valid flag, but no points have been granted.\n") //nolint:forbidigo

	default:
		_, _ = fmt.Printf("Congratulations, you score your team %d points!\n", resp.Value) //nolint:forbidigo
	}

	// And show any message we received
	if resp.ReturnString != "" {
		_, _ = fmt.Printf("Message: %s\n", resp.ReturnString) //nolint:forbidigo
	}

	return nil
}
