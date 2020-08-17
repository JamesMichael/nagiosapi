package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/jamesmichael/nagiosapi/nagios/cmd"
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
	viper.SetDefault("api.addr", ":3000")
	viper.BindPFlag("api.addr", serverCmd.Flags().Lookup("api.addr"))

	var nagiosStatusFile string
	serverCmd.Flags().StringVar(&nagiosStatusFile, "nagios.status-file", "", "path to status.dat")
	viper.SetDefault("nagios.status_file", "/var/log/nagios/status.dat")
	viper.BindPFlag("nagios.status_file", serverCmd.Flags().Lookup("nagios.status-file"))

	var externalCommandsFile string
	serverCmd.Flags().StringVar(&externalCommandsFile, "nagios.external-commands-file", "", "path to nagios commands file")
	viper.SetDefault("nagios.external_commands_file", "/usr/local/nagios/var/rw/nagios.cmd")
	viper.BindPFlag("nagios.external_commands_file", serverCmd.Flags().Lookup("nagios.external-commands-file"))

	viper.SetDefault("app.production", true)
}

func serverCmdFunc(cmd *cobra.Command, args []string) {
	log := mustBuildLog()

	server := mustBuildAPIServer(
		log,
		mustBuildCommandsHandler(log),
		mustBuildStatusHandler(
			log,
			mustBuildStatusRepo(log),
		),
	)

	server.ServeHTTP()
}

func mustBuildCommandsHandler(l *zap.Logger) *server.CommandsHandler {
	commandsFile := viper.GetString("nagios.external_commands_file")
	ecf, err := cmd.NewWriter(
		cmd.WithFilename(commandsFile),
		cmd.WithLogger(l),
		cmd.WithNonBlocking(true),
	)
	if err != nil {
		l.Fatal("unable to open new external commands writer",
			zap.String("filename", commandsFile),
			zap.Error(err),
		)
	}

	go ecf.Run()

	h, err := server.NewCommandsHandler(ecf)
	if err != nil {
		l.Fatal("unable to create commands handler",
			zap.Error(err),
		)
	}
	return h
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

	opts := []statusdata.RepositoryOpt{
		statusdata.WithLog(l),
	}

	if viper.GetBool("nagios.reload_status_file") {
		opts = append(opts, statusdata.WithRefresh(time.Duration(viper.GetInt("nagios.reload_interval"))*time.Second))
	}
	r, err := statusdata.NewRepository(statusFile, opts...)
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

func mustBuildAPIServer(l *zap.Logger, ch *server.CommandsHandler, sh *server.StatusHandler) *server.Server {
	addr := viper.GetString("api.addr")
	s, err := server.NewServer(
		server.WithAddr(addr),
		server.WithCommandsHandler(ch),
		server.WithLog(l),
		server.WithStatusHandler(sh),
	)
	if err != nil {
		l.Fatal("unable to start server",
			zap.String("addr", addr),
			zap.Error(err),
		)
	}

	return s
}
