package config

import (
	"encoding/json"
	"flag"
	"github.com/devops-salt/src/util"
	"io/ioutil"
	"log"
)

// Configure options
type Configure struct {
	LogFile        string       `json:"logger_logfile"`
	LogLevel       string       `json:"logger_level"`
	LogMaxSize     int          `json:"logger_max_size"`
	LogMaxBackups  int          `json:"logger_max_backups"`
	LogMaxDuration int          `json:"logger_max_duration"`
	UDPIPv4        string       `json:"udp_ipv4"`
	UDPPort        int          `json:"udp_port"`
	HTTPPort       int          `json:"http_port"`
	CallbackUrl    string       `json:"callback_url"`
	RunMode        string        `json:"run_mode"`
	DownloadUrl    string        `json:"download_url"`
}

var cfg Configure

// GetConfigure configure
func GetConfigure() Configure {
	return cfg
}

// GetLogFile gets logfile path
func GetLogFile() string {
	return cfg.LogFile
}

// GetLogMaxSize gets the max size of log file
func GetLogMaxSize() int {
	return cfg.LogMaxSize
}

func GetLogLevel() string {
	return cfg.LogLevel
}

func GetCallbackUrl() string {
	return cfg.CallbackUrl
}

func GetDownloadUrl() string {
	return cfg.DownloadUrl
}

func GetRunMod() string {
	return cfg.RunMode
}

// GetLogMaxBackups gets the max backups of log file
func GetLogMaxBackups() int {
	return cfg.LogMaxBackups
}

// GetLogMaxDuration gets the max duration of log file
func GetLogMaxDuration() int {
	return cfg.LogMaxDuration
}

// GetTCPPort gets udp listening port
func GetUDPPort() int {
	return cfg.UDPPort
}

// GetTCPIPv4 gets udp listening ip
func GetUDPIPv4() string {
	return cfg.UDPIPv4
}

// GetHTTPPort gets HTTP listening port
func GetHTTPPort() int {
	return cfg.HTTPPort
}

// GetSource gets Source
func GetSource() [] string {
	return source
}

// configfile config file path
var configfile string

// PidFile PID file path
//var PidFile string

// params client or server
var server string

func GetServer()  string{
	return server
}

// source match
var source []string

func init() {
	flag.StringVar(&configfile, "configfile", "C:\\golang\\src\\github.com\\devops-salt\\etc\\salt_client.json", "configure file for salt")
	////flag.StringVar(&PidFile, "pidfile", "/usr/local/salt/var/run/salt.pid", "pidfile for ngb_kinton")
	flag.StringVar(&server, "server", "client", "server for salt")
	flag.Parse()
	source = util.GetLocalIP()

	data, err := ioutil.ReadFile(configfile)
	if err != nil {
		log.Fatal("[Error] ", err.Error())
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		log.Fatal("[Error] ", err.Error())
	}
}
