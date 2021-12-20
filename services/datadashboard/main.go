package main

import (
	"context"
	"flag"
	goflag "flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

var version = "?"

var (
	rootCmd = &cobra.Command{
		Use:   "datadashboard",
		Short: "datadashboard for reference",
		Long:  "A datadashboard for reference in a monorepo",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			goflag.Parse()
		},
	}
)

func init() {
	rootCmd.AddCommand(RunCmd())
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
}

func run(quit context.Context, done chan error) {
	for {
		select {
		case <-quit.Done():
			done <- nil
			return
		}
	}
}

func RunCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "run datadashboard",
		Long:  "Run the datadashboard as a long-running service",
		RunE: func(cmd *cobra.Command, args []string) error {
			defer func(begin time.Time) {
				glog.Infof("stop datadashboard after %v", time.Since(begin))
			}(time.Now())

			glog.Infof("Start datadashboard on %v", time.Now())

			quit, cancel := context.WithCancel(context.TODO())
			done := make(chan error)
			go run(quit, done)

			go func() {
				sigch := make(chan os.Signal)
				signal.Notify(sigch, syscall.SIGINT, syscall.SIGTERM)
				glog.Info(<-sigch)
				cancel()
			}()

			return <-done
		},
	}

	return cmd
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		glog.Fatalf("%v", err)
	}
}
