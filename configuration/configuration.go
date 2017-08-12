package configuration

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
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
	Debug bool `json:"debug"`
}

// GetConfig reads the configuration from a configuration file and returns
// a configuration struct
func GetConfig(f string) *Struct {
	file, _ := os.Open(f)
	decoder := json.NewDecoder(file)
	config := Struct{}
	decerr := decoder.Decode(&config)
	if decerr != nil {
		fmt.Println("error parsing the configuration:", decerr)
	}
	return &config
}

// ValidateConfiguration validates a given configuration
func (c *Struct) ValidateConfiguration() int {
	if c.Database.Database == "" {
		log.Fatal("Database name is empty")
	}
	if c.Database.Port == 0 {
		return 2
	}
	return 0
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
	result = fmt.Sprintf("%sdebug:%t\n\n", result, c.Debug)

	return result
}
