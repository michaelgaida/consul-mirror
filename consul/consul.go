package consul

import (
	"fmt"
	"log"
	"strconv"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/michaelgaida/consul-mirror/configuration"
)

// Consul represents a consul connection
type Consul struct {
	client *consulapi.Client
	debug  bool
}

type KV struct {
	// Key is the name of the key. It is also part of the URL path when accessed
	// via the API.
	Key string

	// CreateIndex holds the index corresponding the creation of this KVPair. This
	// is a read-only field.
	CreateIndex uint64

	// ModifyIndex is used for the Check-And-Set operations and can also be fed
	// back into the WaitIndex of the QueryOptions in order to perform blocking
	// queries.
	ModifyIndex uint64

	// LockIndex holds the index corresponding to a lock on this key, if any. This
	// is a read-only field.
	LockIndex uint64

	// Flags are any user-defined flags on the key. It is up to the implementer
	// to check these values, since Consul does not treat them specially.
	Flags uint64

	// Value is the value for the key. This can be any value, but it will be
	// base64 encoded upon transport.
	Value []byte

	// Session is a string representing the ID of the session. Any other
	// interactions with this key over the same session must specify the same
	// session ID.
	Session string

	// Not implemented yet in consul API
	Regex string
}

type Service struct {
	id      string
	name    string
	address string
	port    int
	tags    []string
}

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

// GetKVs gets a list of KeyValues from Consul
func (c *Consul) GetKVs(key string) []KV {
	result := c.makeKVRequest(key)
	if c.debug {
		for i := range result {
			fmt.Println(result[i].printKV())
		}
	}
	return result
}

// GetKVs gets a list of KeyValues from Consul
func (c *Consul) GetServices() []Service {
	result := c.makeServiceRequest()
	if c.debug {
		for i := range result {
			fmt.Println(result[i].printService())
		}
	}
	return result
}

func (c *Consul) makeKVRequest(key string) []KV {
	var result []KV
	kv := c.client.KV()

	// var pkvp *consulapi.KVPair
	// pkvp.

	// kv.Put()

	keys, _, err := kv.Keys(key, ":", nil)
	if err != nil {
		log.Fatal(err)
	}
	for i := range keys {
		kvc, _, err := kv.Get(keys[i], nil)
		if err != nil {
			log.Fatal(err)
		}
		kv := KV{}
		kv.Key = kvc.Key
		kv.Value = kvc.Value
		kv.CreateIndex = kvc.CreateIndex
		kv.Flags = kvc.Flags
		kv.LockIndex = kvc.LockIndex
		kv.ModifyIndex = kvc.ModifyIndex
		kv.Session = kvc.Session
		// To be implemented in the API
		kv.Regex = ""

		result = append(result, kv)
	}
	return result
}

func (kv *KV) printKV() string {
	return fmt.Sprintf("Key: %s; Value: %s; RegEx: %s; CreateIndex: %d; Flags: %d; LockIndex: %d; ModifyIndex: %d, Session: %s",
		kv.Key,
		kv.Value,
		kv.Regex,
		kv.CreateIndex,
		kv.Flags,
		kv.LockIndex,
		kv.ModifyIndex,
		kv.Session)
}

func (c *Consul) makeServiceRequest() []Service {
	var result []Service
	cat := c.client.Catalog()
	services, _, err := cat.Service("consul", "", nil)
	if err != nil {
		log.Fatal(err)
	}
	for i := range services {
		fmt.Println(services[i].ServiceName)
		fmt.Println(services[i].ServicePort)
		fmt.Println(services[i].ServiceID)

		var service Service
		service.id = services[i].ServiceID
		service.name = services[i].ServiceName
		service.address = services[i].ServiceAddress
		service.port = services[i].ServicePort
		service.tags = services[i].ServiceTags

		result = append(result, service)
	}
	return result
}

func (service *Service) printService() string {
	return fmt.Sprintf("Id: %s; Address: %s; Name: %s; Port: %d; Tags: %v",
		service.id,
		service.address,
		service.name,
		service.port,
		service.tags)
}
