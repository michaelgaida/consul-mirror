package consul

import (
	"log"
	"strconv"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/michaelgaida/consul-mirror/configuration"
)

// GetConsul gets a
func GetConsul(config *configuration.Struct) *Consul {
	clientConfig := consulapi.DefaultConfig()
	clientConfig.Address = config.Consul.Host + ":" + strconv.Itoa(config.Consul.Port)
	clientConfig.Token = config.Consul.ACL
	consul, err := consulapi.NewClient(clientConfig)
	if err != nil {
		log.Fatal(err)
	}

	c := Consul{}
	c.client = consul
	c.debug = config.Debug
	return &c
}

// GetKVs gets a list of KeyValues from Consul from all given DCs
func (c *Consul) GetKVs(key string, dcs []DC) []KV {
	result := c.makeKVRequest(key, dcs)
	if c.debug {
		for i := range result {
			log.Println(result[i].printKV())
		}
	}
	return result
}

// GetDCs gets a list of available DCs from Consul
func (c *Consul) GetDCs() []DC {
	result := c.makeDCRequest()
	if c.debug {
		for i := range result {
			log.Println(result[i].printDC())
		}
	}
	return result
}

// GetServices gets a list of KeyValues from Consul
func (c *Consul) GetServices() []CatalogService {
	result := c.makeServiceRequest()
	if c.debug {
		for i := range result {
			log.Println(result[i].printService())
		}
	}
	return result
}

// GetACLs gets a list of ACLs from Consul
func (c *Consul) GetACLs() []ACL {
	result := c.makeACLRequest()
	if c.debug {
		for i := range result {
			log.Println(result[i].printACL())
		}
	}
	return result
}

// GetNodes gets a list of Nodes from Consul
func (c *Consul) GetNodes() []Node {
	result := c.makeNodeRequest()
	if c.debug {
		for i := range result {
			log.Println(result[i].printNode())
		}
	}
	return result
}

// WriteKVs writes KVs to consul
func (c *Consul) WriteKVs(kvs []KV, keepDC bool) error {
	if c.debug {
		for i := range kvs {
			log.Println(kvs[i].printKV())
		}
	}
	return c.writeKV(kvs, keepDC)
}
