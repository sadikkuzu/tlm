package config

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
	"github.com/yusufcanb/tlm/shell"
	"net/url"
)

func (c *Config) Action(_ *cli.Context) error {
	var err error

	form := ConfigForm{
		host:    viper.GetString("llm.host"),
		explain: viper.GetString("llm.explain"),
		suggest: viper.GetString("llm.suggest"),
	}

	err = form.Run()
	if err != nil {
		return err
	}

	viper.Set("shell", form.shell)
	viper.Set("llm.host", form.host)
	viper.Set("llm.explain", form.explain)
	viper.Set("llm.suggest", form.suggest)

	err = viper.WriteConfig()
	if err != nil {
		return err
	}

	fmt.Println(shell.Ok() + " configuration saved")
	return nil
}

func (c *Config) Command() *cli.Command {
	return &cli.Command{
		Name:    "config",
		Aliases: []string{"c"},
		Usage:   "Configures tlm preferences.",
		Action:  c.Action,
		Subcommands: []*cli.Command{
			c.subCommandGet(),
			c.subCommandSet(),
		},
	}
}

func (c *Config) subCommandGet() *cli.Command {
	return &cli.Command{
		Name:      "get",
		Usage:     "get configuration by key",
		UsageText: "tlm config get <key>",
		Action: func(c *cli.Context) error {
			key := c.Args().Get(0)
			value := viper.GetString(key)

			if value == "" {
				fmt.Println(fmt.Sprintf("%s <%s> is not a tlm parameter", shell.Err(), key))
				return nil
			}

			fmt.Println(fmt.Sprintf("%s = %s", key, value))
			return nil
		},
	}
}

func (c *Config) subCommandSet() *cli.Command {
	return &cli.Command{
		Name:  "set",
		Usage: "set configuration",
		Action: func(c *cli.Context) error {
			key := c.Args().Get(0)

			switch key {
			case "llm.host":
				u, err := url.ParseRequestURI(c.Args().Get(1))
				if err != nil {
					return errors.New("Invalid url: " + c.Args().Get(1))
				}
				viper.Set(key, u.String())

			case "llm.suggest", "llm.explain":
				mode := c.Args().Get(1)
				if mode != "stable" && mode != "balanced" && mode != "creative" {
					return errors.New("Invalid mode: " + mode)
				}
				viper.Set(key, mode)
			default:
				fmt.Println(fmt.Sprintf("%s <%s> is not a tlm parameter", shell.Err(), key))
				return nil
			}

			viper.Set(key, c.Args().Get(1))
			err := viper.WriteConfig()
			if err != nil {
				return err
			}

			fmt.Println(fmt.Sprintf("%s = %s  %s", key, c.Args().Get(1), shell.Ok()))
			return nil
		},
	}
}
