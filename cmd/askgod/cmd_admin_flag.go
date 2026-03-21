package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v3"

	"github.com/nsec/askgod/api"
	"github.com/nsec/askgod/internal/utils"
)

func (c *client) cmdAdminAddFlag(ctx context.Context, cmd *cli.Command) error {
	flag := api.AdminFlagPost{}

	if cmd.NArg() > 0 {
		for _, arg := range cmd.Args().Slice() {
			err := setStructKey(&flag, arg)
			if err != nil {
				return err
			}
		}
	}

	err := c.queryStruct(ctx, "POST", "/flags", flag, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *client) cmdAdminDeleteFlag(ctx context.Context, cmd *cli.Command) error {
	if cmd.NArg() != 1 {
		_ = cli.ShowSubcommandHelp(cmd)

		return nil
	}

	err := c.queryStruct(ctx, "DELETE", "/flags/"+cmd.Args().Get(0), nil, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *client) cmdAdminImportFlags(ctx context.Context, cmd *cli.Command) error {
	if cmd.NArg() < 1 {
		_ = cli.ShowSubcommandHelp(cmd)

		return nil
	}

	// Flush all existing entries
	if cmd.Bool("flush") {
		reader := bufio.NewReader(os.Stdin)
		_, _ = fmt.Print("Flush all flags (yes/no): ") //nolint:forbidigo
		input, _ := reader.ReadString('\n')

		input = strings.TrimSuffix(input, "\n")
		if strings.TrimSpace(strings.ToLower(input)) != "yes" {
			return errors.New("user aborted flush operation")
		}

		err := c.queryStruct(ctx, "DELETE", "/flags?empty=1", nil, nil)
		if err != nil {
			return err
		}
	}

	// Read the file
	content, err := os.ReadFile(cmd.Args().Get(0))
	if err != nil {
		return err
	}

	// Parse the JSON file
	flags := []api.AdminFlag{}

	err = json.Unmarshal(content, &flags)
	if err != nil {
		return err
	}

	// Create the flags
	err = c.queryStruct(ctx, "POST", "/flags?bulk=1", flags, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *client) cmdAdminListFlags(ctx context.Context, _ *cli.Command) error {
	// Get the data
	resp := []api.AdminFlag{}

	err := c.queryStruct(ctx, "GET", "/flags", nil, &resp)
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Flag", "Value", "Return string", "Description", "Tags"})
	table.SetBorder(false)
	table.SetAutoWrapText(false)

	for _, entry := range resp {
		table.Append([]string{
			strconv.FormatInt(entry.ID, 10),
			entry.Flag,
			strconv.FormatInt(entry.Value, 10),
			entry.ReturnString,
			entry.Description,
			utils.PackTags(entry.Tags),
		})
	}

	table.Render()

	return nil
}

func (c *client) cmdAdminUpdateFlag(ctx context.Context, cmd *cli.Command) error {
	if cmd.NArg() < 1 {
		_ = cli.ShowSubcommandHelp(cmd)

		return nil
	}

	flag := api.AdminFlag{}

	err := c.queryStruct(ctx, "GET", "/flags/"+cmd.Args().Get(0), nil, &flag)
	if err != nil {
		return err
	}

	if cmd.NArg() > 1 {
		for _, arg := range cmd.Args().Slice()[1:] {
			err := setStructKey(&flag, arg)
			if err != nil {
				return err
			}
		}
	}

	err = c.queryStruct(ctx, "PUT", "/flags/"+cmd.Args().Get(0), flag.AdminFlagPut, nil)
	if err != nil {
		return err
	}

	return nil
}
