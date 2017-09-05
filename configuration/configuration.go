package configuration

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/urfave/cli"
)

// Struct that represents the configuration
type Struct struct {
	Database struct {
		User     string `json:"user"`
		Password string `json:"password"`
		Host     string `json:"host"`
		Database string `json:"database"`
		Port     int    `json:"port"`
	} `json:"database"`
	Consul struct {
		ACL  string `json:"acl"`
		Host string `json:"host"`
		Port int    `json:"port"`
	} `json:"consul"`
	LogLevel string `json:"loglevel"`
}

func (c *Struct) setDefaults() {
	c.Consul.Port = 8500
	c.Consul.Host = "localhost"
	c.Database.Port = 1433
	c.LogLevel = "info"
}

// GetConfig reads the configuration from a configuration file and returns
// a configuration struct
func GetConfig(f string) *Struct {
	file, _ := os.Open(f)
	decoder := json.NewDecoder(file)
	config := Struct{}
	config.setDefaults()
	decerr := decoder.Decode(&config)
	if decerr != nil {
		fmt.Println("error parsing the configuration:", decerr)
	}
	return &config
}

// ValidateConfiguration validates a given configuration
func (c *Struct) ValidateConfiguration() error {
	if c.Database.User == "" {
		return errors.New("Database user is empty")
	}
	if c.Database.Password == "" {
		return errors.New("Database password is empty")
	}
	if c.Database.Host == "" {
		return errors.New("Database host is empty")
	}
	if c.Database.Database == "" {
		return errors.New("Database name is empty")
	}
	if c.Database.Port == 0 {
		return errors.New("Database port is 0")
	}
	if c.Consul.Host == "" {
		return errors.New("Consul host is empty")
	}
	if c.Consul.Port == 0 {
		return errors.New("Consul port is 0")
	}
	return nil
}

// OverwriteConfig updates the configuration with cli context
func (c *Struct) OverwriteConfig(cli *cli.Context) {
	if cli.GlobalString("token") != "" {
		c.Consul.ACL = cli.GlobalString("token")
	}
	if cli.GlobalString("dbpassword") != "" {
		c.Database.Password = cli.GlobalString("dbpassword")
	}
	if cli.GlobalString("loglevel") != "" {
		c.LogLevel = cli.GlobalString("loglevel")
	}
}

// PrintDebug returns a string representation of the configuration
func (c *Struct) PrintDebug() string {
	var result string
	result = fmt.Sprintf("=== Configuration ===\n")
	result = fmt.Sprintf("%sdatabase.port:%d\n", result, c.Database.Port)
	result = fmt.Sprintf("%sdatabase.sserver:%s\n", result, c.Database.Host)
	result = fmt.Sprintf("%ssdatabase.user:%s\n", result, c.Database.User)
	result = fmt.Sprintf("%ssdatabase.database:%s\n\n", result, c.Database.Database)
	result = fmt.Sprintf("%sconsul.acl:%s\n", result, c.Consul.ACL)
	result = fmt.Sprintf("%sconsul.host:%s\n", result, c.Consul.Host)
	result = fmt.Sprintf("%sconsul.port:%d\n\n", result, c.Consul.Port)
	result = fmt.Sprintf("%sloglevel:%s\n\n", result, c.LogLevel)

	return result
}
