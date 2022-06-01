package trust

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/google/go-containerregistry/pkg/v1/remote/transport"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	RootPath             string
	ServerUrl            string `json:"server_url"`
	RootPassphrase       string `json:"root_passphrase"`
	RepositoryPassphrase string `json:"repository_passphrase"`
	Scopes               string `json:"scopes,omitempty"`
}

const (
	// configDir is the root path for configuration
	configDirEnv          = "NOTARY_CONFIG_DIR"
	configFileNameEnv     = "NOTARY_CONFIG_FILENAME"
	defaultConfigFileName = "gcr-config.json"
)

// ParseConfig read configfile (${configDir}/${configFileName})
// returns a Config object and error.
func ParseConfig(configDir string) (*Config, error) {
	if configDir == "" {
		configDir = os.Getenv(configDirEnv)
		if configDir == "" {
			configDir = filepath.Join(os.Getenv("HOME"), ".notary")
		}
	}
	if !filepath.IsAbs(configDir) {
		log.Warnf("config directory %s maybe wrong, not absolute path", configDir)
	}

	configFileName := os.Getenv(configFileNameEnv)
	if configFileName == "" {
		configFileName = defaultConfigFileName
	}

	configFilePath := filepath.Join(configDir, configFileName)
	configFile, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}

	c := new(Config)
	err = json.Unmarshal([]byte(configFile), c)
	if err != nil {
		return nil, err
	}
	c.Scopes = parseScopes(c)
	c.RootPath = configDir
	return c, nil
}

func parseScopes(config *Config) string {
	if config.Scopes == "" {
		return transport.PullScope
	}
	supportScopes := []string{
		transport.PullScope,
		transport.PushScope,
		transport.DeleteScope,
		transport.CatalogScope,
	}
	if !contains(supportScopes, config.Scopes) {
		log.Warnf("Scope %s is not supported. Supported: ['pull', 'push,pull', 'catalog'] ", config.Scopes)
		return transport.PullScope
	}
	return config.Scopes
}

func contains(sl []string, v string) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}
