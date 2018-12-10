package nsupdate

import (
	"fmt"
	"path"
	"regexp"
	"strings"
)

const (
	// BINDMAN_NAMESERVER_ADDRESS environment variable identifier for the nameserver address
	BINDMAN_NAMESERVER_ADDRESS = "BINDMAN_NAMESERVER_ADDRESS"

	// BINDMAN_NAMESERVER_PORT environment variable identifier for the nameserver port
	BINDMAN_NAMESERVER_PORT = "BINDMAN_NAMESERVER_PORT"

	// BINDMAN_NAMESERVER_KEYFILE environment variable identifier for the nameserver key name
	BINDMAN_NAMESERVER_KEYFILE = "BINDMAN_NAMESERVER_KEYFILE"

	// BINDMAN_NAMESERVER_ZONE environment variable identifier for the zone to be managed
	BINDMAN_NAMESERVER_ZONE = "BINDMAN_NAMESERVER_ZONE"

	// BINDMAN_MODE defines the execution mode of bindman (DEBUG|PROD); defaults to PROD
	BINDMAN_MODE = "BINDMAN_MODE"
)

// check tests if a NSUpdate setup is ok; returns a set of error strings in case something is not right
func (nsu *NSUpdate) check() (success bool, errs []string) {
	success = true
	errMsg := "The environment variable %s cannot be empty"
	if nsu.Server == "" {
		errs = append(errs, fmt.Sprintf(errMsg, BINDMAN_NAMESERVER_ADDRESS))
		success = false
	}

	if nsu.KeyFile == "" {
		errs = append(errs, fmt.Sprintf(errMsg, BINDMAN_NAMESERVER_KEYFILE))
		success = false
	}

	if nsu.Zone == "" {
		errs = append(errs, fmt.Sprintf(errMsg, BINDMAN_NAMESERVER_ZONE))
		success = false
	}

	m := `K.*\.\+157\+.*\.key`
	if succ, _ := regexp.MatchString(m, nsu.KeyFile); !succ {
		errs = append(errs, fmt.Sprintf("Environment variable %s did not match the regex %v: %s", BINDMAN_NAMESERVER_KEYFILE, m, nsu.KeyFile))
		success = false
	}

	// TODO: Test connection

	return success, errs
}

// getKeyFilePath joins the base path with key file name
func (nsu *NSUpdate) getKeyFilePath() string {
	return path.Join(nsu.BasePath, nsu.KeyFile)
}

// getSubdomainName we expect names to come in the format subdomain.zone. This function returns the subdomain part
func (nsu *NSUpdate) getSubdomainName(name string) string {
	str := strings.Replace(name, nsu.Zone, "", 1)
	return strings.TrimSuffix(str, ".")
}

// checkName checks if the name is in the expected format: subdomain.zone
func (nsu *NSUpdate) checkName(name string) (bool, error) {
	m := fmt.Sprintf(".*\\.%s", nsu.Zone)
	return regexp.MatchString(m, name)
}
