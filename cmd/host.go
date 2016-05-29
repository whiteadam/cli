package cmd

import (
	"github.com/codegangsta/cli"
	"github.com/rancher/go-rancher/client"
)

func HostCommand() cli.Command {
	return cli.Command{
		Name:      "hosts",
		ShortName: "host",
		Usage:     "Operations on hosts",
		Action:    hostLs,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "quiet,q",
				Usage: "Only display IDs",
			},
			cli.StringFlag{
				Name:  "format",
				Usage: "'json' or Custom format: {{.Id}} {{.Name}",
			},
		},
	}
}

type HostsData struct {
	ID          string
	Host        client.Host
	IPAddresses []client.IpAddress
}

func hostLs(ctx *cli.Context) error {
	c, err := GetClient(ctx)
	if err != nil {
		return err
	}

	env, err := GetEnvironment(c)
	if err != nil {
		return err
	}

	collection, err := c.Host.List(&client.ListOpts{
		Filters: map[string]interface{}{
			"accountId": env.Id,
		},
	})
	if err != nil {
		return err
	}

	writer := NewTableWriter([][]string{
		{"ID", "Host.Id"},
		{"HOSTNAME", "Host.Hostname"},
		{"STATE", "Host.State"},
		{"IP", "{{ips .IpAddresses}}"},
	}, ctx)

	defer writer.Close()

	for _, item := range collection.Data {
		ips := client.IpAddressCollection{}
		err := c.GetLink(item.Resource, "ipAddresses", &ips)
		if err != nil {
			return err
		}

		writer.Write(&HostsData{
			ID:          item.Id,
			Host:        item,
			IPAddresses: ips.Data,
		})
	}

	return writer.Err()
}