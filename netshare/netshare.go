package netshare

import (
	"fmt"
	"github.com/ContainX/docker-volume-netshare/netshare/drivers"
	log "github.com/Sirupsen/logrus"
	"github.com/docker/go-plugins-helpers/volume"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
)

const (
	UsernameFlag     = "username"
	PasswordFlag     = "password"
	  TCPFlag          = "tcp"
	PortFlag         = "port"
VerboseFlag      = "verbose"
OptionsFlag      = "options"
	BasedirFlag      = "basedir"


	)

var (
	rootCmd = &cobra.Command{
		Use:              "docker-volume-netshare",
		Short:            "NFS and CIFS - Docker volume driver plugin",
		Long:             NetshareHelp,
		PersistentPreRun: setupLogger,
	}

	cifsCmd = &cobra.Command{
		Use:   "cifs",
		Short: "run plugin in CIFS mode",
		Run:   execCIFS,
	}

	)

func Execute() {
	setupFlags()
	rootCmd.Long = fmt.Sprintf(NetshareHelp, Version, BuildDate)
	rootCmd.AddCommand(cifsCmd)
	rootCmd.Execute()
}

func setupFlags() {
	rootCmd.PersistentFlags().StringVar(&baseDir, BasedirFlag, filepath.Join(volume.DefaultDockerRootDirectory, PluginAlias), "Mounted volume base directory")
	rootCmd.PersistentFlags().Bool(TCPFlag, false, "Bind to TCP rather than Unix sockets.  Can also be set via NETSHARE_TCP_ENABLED")
	rootCmd.PersistentFlags().String(PortFlag, ":8877", "TCP Port if --tcp flag is true.  :PORT for all interfaces or ADDRESS:PORT to bind.")
	rootCmd.PersistentFlags().Bool(VerboseFlag, false, "Turns on verbose logging")

	cifsCmd.Flags().StringP(UsernameFlag, "u", "", "Username to use for mounts.  Can also set environment NETSHARE_CIFS_USERNAME")
	cifsCmd.Flags().StringP(PasswordFlag, "p", "", "Password to use for mounts.  Can also set environment NETSHARE_CIFS_PASSWORD")
	cifsCmd.Flags().StringP(OptionsFlag, "o", "", "Options passed to Cifs mounts (ex: nounix,uid=433)")

}

func setupLogger(cmd *cobra.Command, args []string) {
	if verbose, _ := cmd.Flags().GetBool(VerboseFlag); verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}

func execCIFS(cmd *cobra.Command, args []string) {
	user := typeOrEnv(cmd, UsernameFlag, EnvSambaUser)
	pass := typeOrEnv(cmd, PasswordFlag, EnvSambaPass)
		options, _ := cmd.Flags().GetString(OptionsFlag)

	creds := drivers.NewCifsCredentials(user, pass, domain, security, fileMode, dirMode)

	d := drivers.NewCIFSDriver(rootForType(drivers.CIFS), creds, netrc, options)
	if len(user) > 0 {
		startOutput(fmt.Sprintf("CIFS :: %s, opts: %s", creds, options))
	} else {
		startOutput(fmt.Sprintf("CIFS :: netrc: %s, opts: %s", netrc, options))
	}
	start(drivers.CIFS, d)
}

func startOutput(info string) {
	log.Infof("== docker-volume-netshare :: Version: %s - Built: %s ==", Version, BuildDate)
	log.Infof("Starting %s", info)
}

func typeOrEnv(cmd *cobra.Command, flag, envname string) string {
	val, _ := cmd.Flags().GetString(flag)
	if val == "" {
		val = os.Getenv(envname)
	}
	return val
}

func rootForType(dt drivers.DriverType) string {
	return filepath.Join(baseDir, dt.String())
}

func start(dt drivers.DriverType, driver volume.Driver) {
	h := volume.NewHandler(driver)
	if isTCPEnabled() {
		addr := os.Getenv(EnvTCPAddr)
		if addr == "" {
			addr, _ = rootCmd.PersistentFlags().GetString(PortFlag)
		}
		fmt.Println(h.ServeTCP(dt.String(), addr, nil))
	} else {
		fmt.Println(h.ServeUnix(dt.String(), syscall.Getgid()))
	}
}

func isTCPEnabled() bool {
	if tcp, _ := rootCmd.PersistentFlags().GetBool(TCPFlag); tcp {
		return tcp
	}

	if os.Getenv(EnvTCP) != "" {
		ev, _ := strconv.ParseBool(os.Getenv(EnvTCP))
		fmt.Println(ev)

		return ev
	}
	return false
}
