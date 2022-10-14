/*
Copyright © 2022 Darren O'Connor mellow.drifter@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package tracerpki

import (
	"fmt"
	"os"

	"github.com/mellowdrifter/tracerpki/pkg/tracerpki"
	"github.com/spf13/cobra"
)

const version = "0.1"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "tracerpki",
	Version: version,
	Short:   "traceroute replacement",
	Args:    cobra.ExactArgs(1),
	Long: `traceroute replacement that gives RPKI, AS number, and AS name per hop

Usage:
  traceroute [ -46dFITnreAUDV ] [ -f first_ttl ] [ -m max_ttl ] [ -N squeries ] [ -p port ] [ -t tos ] [ -q nqueries ] [ -s src_addr ] host 
Options:
  -F  --dont-fragment         Do not fragment packets

`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := getOptions(cmd)
		//TODO: Move this into get Options
		opts.Location = args[0]
		fmt.Printf("options are %+v\n", opts)
		tracerpki.Trace(*opts)
	},
}

func getOptions(cmd *cobra.Command) *tracerpki.Args {
	var opt tracerpki.Args

	v6, err := cmd.Flags().GetBool("ipv6")
	if err != nil {
		panic(err)
	}

	v4, err := cmd.Flags().GetBool("ipv4")
	if err != nil {
		panic(err)
	}

	max_ttl, err := cmd.Flags().GetUint("max_ttl")
	if err != nil {
		panic(err)
	}

	opt.V6 = v6
	opt.V4 = v4
	opt.MaxTTL = max_ttl

	return &opt
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("ipv6", "6", false, "Use IPv6 (default)")
	rootCmd.Flags().BoolP("ipv4", "4", true, "Use IPv4")
	rootCmd.Flags().UintP("max_ttl", "m", 30, "Set the max number of hops (max TTL to be reached). Default is 30")
}
