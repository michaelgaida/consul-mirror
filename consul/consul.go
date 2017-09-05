package consul

import (
	"strconv"

	log "github.com/sirupsen/logrus"

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
	return &c
}

// GetKVs gets a list of KeyValues from Consul from all given DCs
func (c *Consul) GetKVs(key string, dcs []DC) ([]KV, error) {
	result, err := c.makeKVRequest(key, dcs)
	if err != nil {
		log.Error("Error making KV request: ", err.Error())
		return nil, err
	}
	for i := range result {
		log.Debug(result[i].printKV())
	}
	return result, nil
}

// GetDCs gets a list of available DCs from Consul
func (c *Consul) GetDCs(dc string) ([]DC, error) {
	result, err := c.makeDCRequest(dc)
	if err != nil {
		log.Error("Error making DC request: ", err.Error())
		return nil, err
	}
	for i := range result {
		log.Debug(result[i].printDC())
	}
	return result, nil
}

// WriteKVs writes KVs to consul
func (c *Consul) WriteKVs(kvs []KV, ignoreDC, dcprefix bool, setprefix string) error {
	for i := range kvs {
		log.Debug(kvs[i].printKV())
	}
	return c.writeKV(kvs, ignoreDC, dcprefix, setprefix)
}

// Equals checks if to kv are equal in tearms of content
func (kv *KV) Equals(ckv KV) bool {
	return (kv.Datacenter == ckv.Datacenter) &&
		(kv.Key == ckv.Key) &&
		(string(kv.Value[:]) == string(ckv.Value[:]))
}

// // GetServices gets a list of KeyValues from Consul
// func (c *Consul) GetServices() []CatalogService {
// 	result := c.makeServiceRequest()
// 	for i := range result {
// 		log.Debug(result[i].printService())
// 	}
// 	return result
// }

// // GetACLs gets a list of ACLs from Consul
// func (c *Consul) GetACLs() []ACL {
// 	result := c.makeACLRequest()
// 	for i := range result {
// 		log.Debug(result[i].printACL())
// 	}
// 	return result
// }

// // GetNodes gets a list of Nodes from Consul
// func (c *Consul) GetNodes() []Node {
// 	result := c.makeNodeRequest()
// 	for i := range result {
// 		log.Debug(result[i].printNode())
// 	}
// 	return result
// }
