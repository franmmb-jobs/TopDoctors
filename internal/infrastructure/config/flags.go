package config

import (
	"flag"
	"os"
	"strings"
)

var Flags struct {
	Config    string
	InTestEnv bool
}

func init() {
	isTestEnv()

	if !Flags.InTestEnv {
		flag.StringVar(&Flags.Config, "config", "configs/config.yml", "Route to YAML config file")

		flag.Parse()
	}

}

func isTestEnv() bool {

	if Flags.InTestEnv {
		return true
	}

	//Principal way to detect go test
	if flag.Lookup("test.v") != nil {
		Flags.InTestEnv = true
		return true
	}

	//Secondary ways to detect go test
	binary := os.Args[0]
	secdCheck := os.Getenv("APP_ENV") == "test" ||
		strings.HasSuffix(binary, ".test") ||
		strings.HasSuffix(binary, ".test.exe") ||
		strings.HasSuffix(binary, "test.exe")

	Flags.InTestEnv = secdCheck
	return secdCheck
}
