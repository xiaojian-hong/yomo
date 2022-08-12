package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/yomorun/cli/serverless"
)

func parseURL(url string, opts *serverless.Options) error {
	if url == "" {
		url = "localhost:9000"
	}
	splits := strings.Split(url, ":")
	if len(splits) != 2 {
		return fmt.Errorf(`the format of url "%s" is incorrect, it should be "host:port", f.e. localhost:9000`, url)
	}
	host := splits[0]
	port, err := strconv.Atoi(splits[1])
	if err != nil {
		return err
	}
	opts.Host = host
	opts.Port = port
	return nil
}

func getViperName(name string) string {
	return "yomo_sfn_" + strings.ReplaceAll(name, "-", "_")
}

func bindViper(cmd *cobra.Command) *viper.Viper {
	v := viper.New()

	// bind environment variables
	v.AllowEmptyEnv(true)
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		name := getViperName(f.Name)
		v.BindEnv(name)
		v.SetDefault(name, f.DefValue)
	})

	return v
}

func loadViperValue(cmd *cobra.Command, v *viper.Viper, p *string, name string) {
	f := cmd.Flag(name)
	if !f.Changed {
		*p = v.GetString(getViperName(name))
	}
}
