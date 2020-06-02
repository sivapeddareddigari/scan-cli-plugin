package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/docker/cli/cli-plugins/manager"
	"github.com/docker/cli/cli-plugins/plugin"
	"github.com/docker/cli/cli/command"
	"github.com/docker/docker-scan/internal"
	"github.com/docker/docker-scan/internal/provider"
	"github.com/spf13/cobra"
)

func main() {
	plugin.Run(func(dockerCli command.Cli) *cobra.Command {
		return newScanCmd(dockerCli)
	}, manager.Metadata{
		SchemaVersion: "0.1.0",
		Vendor:        "Docker Inc.",
		Version:       internal.Version,
	})
}

func newScanCmd(_ command.Cli) *cobra.Command {
	var (
		auth        string
		showVersion bool
	)
	cmd := &cobra.Command{
		Short:       "Docker Scan",
		Long:        `A tool to scan your docker image`,
		Use:         "scan [OPTIONS] IMAGE",
		Annotations: map[string]string{},
		RunE: func(cmd *cobra.Command, args []string) error {
			scanProvider := provider.NewSnykProvider()
			if showVersion {
				version, err := internal.FullVersion(scanProvider)
				if err != nil {
					return err
				}
				fmt.Println(version)
				return nil
			}
			if len(args) != 1 {
				//nolint:golint
				return fmt.Errorf(`"docker run" requires at least 1 argument.
See 'docker scan --help'.`)
			}
			if auth != "" {
				fmt.Println("Authenticating to Snyk using", auth)
				c := exec.Command("snyk", "auth", auth)
				c.Stdin = os.Stdin
				c.Stdout = os.Stdout
				c.Stderr = os.Stderr
				if err := c.Run(); err != nil {
					return err
				}
				fmt.Println("Authenticated")
				fmt.Println()
			}
			c := exec.Command("snyk", "test", "--docker", args[0])
			c.Stdin = os.Stdin
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			return c.Run()
		},
	}
	cmd.Flags().StringVar(&auth, "auth", "", "Use snyk API token to authenticate on snyk.io")
	cmd.Flags().BoolVar(&showVersion, "version", false, "Display version of scan plugin and snyk cli")
	return cmd
}
