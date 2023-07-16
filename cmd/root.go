package cmd

import (
	"fmt"

	"github.com/atscan/atr/util/version"
	"github.com/spf13/cobra"
)

var (
	// Used for flags.
	//cfgFile    string
	workingDir string

	rootCmd = &cobra.Command{
		Use:   "atr",
		Short: "AT Protocol IPLD-CAR Repository toolkit",
		//Long: `.`,
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	//cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&workingDir, "working-dir", "C", ".", "")

	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
	rootCmd.AddCommand(versionCmd)
	//rootCmd.AddCommand(ShowCmd)
}

/*func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		//viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		//home, err := os.UserHomeDir()
		//cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".atr")
	}

	//viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}*/

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version.Version(""))
	},
}
