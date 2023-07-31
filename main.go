package main

import (
	"github.com/joho/godotenv"
	logrus "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var log = logrus.WithField("prefix", "dcim")

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	godotenv.Load()
	DBInit()

	var rootCmd = &cobra.Command{Use: "dcim"}

	cc := &cobra.Command{
		Use:     "zone",
		Aliases: []string{"zones"},
		Short:   "Zone commands",
	}
	cc.AddCommand(NewListCommand(Zone{}))
	cc.AddCommand(NewCreateCommand(Zone{}))
	rootCmd.AddCommand(cc)

	cc = &cobra.Command{
		Use:     "rack",
		Aliases: []string{"racks"},
		Short:   "Rack commands",
	}
	cc.AddCommand(NewListCommand(Rack{}))
	cc.AddCommand(NewCreateCommand(Rack{}))
	rootCmd.AddCommand(cc)

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}

}
