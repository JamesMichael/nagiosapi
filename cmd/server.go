package cmd

import (
	"fmt"
	"os"

	"github.com/jamesmichael/nagiosapi/nagios/statusdata"
	"github.com/jamesmichael/nagiosapi/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	serverCmd = &cobra.Command{
		Use:   "server",
		Short: "Start HTTP API Server",
		Run:   serverCmdFunc,
	}
)

func init() {
	rootCmd.AddCommand(serverCmd)

	var apiAddr string
	serverCmd.Flags().StringVar(&apiAddr, "api.addr", "", "api address")
	viper.SetDefault("api_addr", ":3000")
	viper.BindPFlag("api_addr", serverCmd.Flags().Lookup("api.addr"))

	var nagiosStatusFile string
	serverCmd.Flags().StringVar(&nagiosStatusFile, "nagios.status-file", "", "path to status.dat")
	viper.SetDefault("nagios.status_file", "/var/log/nagios/status.dat")
	viper.BindPFlag("nagios.status_file", serverCmd.Flags().Lookup("nagios.status-file"))

	viper.SetDefault("app.production", true)
}

func serverCmdFunc(cmd *cobra.Command, args []string) {
	log := mustBuildLog()

	server := mustBuildAPIServer(
		log,
		mustBuildStatusHandler(
			log,
			mustBuildStatusRepo(log),
		),
	)

	server.ServeHTTP()
}

func mustBuildLog() *zap.Logger {
	var log *zap.Logger
	var err error
	if viper.GetBool("app.production") {
		log, err = zap.NewProduction()
	} else {
		log, err = zap.NewDevelopment()
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to create logger: %s", err)
		os.Exit(1)
	}

	return log
}

func mustBuildStatusRepo(l *zap.Logger) *statusdata.Repository {
	statusFile := viper.GetString("nagios.status_file")
	r, err := statusdata.NewRepository(statusFile,
		statusdata.WithLog(l),
	)
	if err != nil {
		l.Fatal("unable to read nagios status file",
			zap.String("filename", statusFile),
			zap.Error(err),
		)
	}

	return r
}

func mustBuildStatusHandler(l *zap.Logger, r *statusdata.Repository) *server.StatusHandler {
	h, err := server.NewStatusHandler(r)
	if err != nil {
		l.Fatal("unable to create status handler",
			zap.Error(err),
		)
	}

	return h
}

func mustBuildAPIServer(l *zap.Logger, h *server.StatusHandler) *server.Server {
	addr := viper.GetString("api.addr")
	s, err := server.NewServer(
		server.WithAddr(addr),
		server.WithStatusHandler(h),
	)
	if err != nil {
		l.Fatal("unable to start server",
			zap.String("addr", addr),
			zap.Error(err),
		)
	}

	return s
}
