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

type Node struct {
	ID              string
	Node            string
	Address         string
	Datacenter      string
	TaggedAddresses map[string]string
	Meta            map[string]string
	CreateIndex     uint64
	ModifyIndex     uint64
}

type CatalogService struct {
	ID                       string
	Node                     string
	Address                  string
	Datacenter               string
	TaggedAddresses          map[string]string
	NodeMeta                 map[string]string
	ServiceID                string
	ServiceName              string
	ServiceAddress           string
	ServiceTags              []string
	ServicePort              int
	ServiceEnableTagOverride bool
	CreateIndex              uint64
	ModifyIndex              uint64
}

type ACL struct {
	CreateIndex uint64
	ID          string
	ModifyIndex uint64
	Name        string
	Rules       string
	Type        string
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

// GetServices gets a list of KeyValues from Consul
func (c *Consul) GetServices() []CatalogService {
	result := c.makeServiceRequest()
	if c.debug {
		for i := range result {
			fmt.Println(result[i].printService())
		}
	}
	return result
}

// GetACLs gets a list of ACLs from Consul
func (c *Consul) GetACLs() []ACL {
	result := c.makeACLRequest()
	if c.debug {
		for i := range result {
			fmt.Println(result[i].printACL())
		}
	}
	return result
}

// GetNodes gets a list of Nodes from Consul
func (c *Consul) GetNodes() []Node {
	result := c.makeNodeRequest()
	if c.debug {
		for i := range result {
			fmt.Println(result[i].printNode())
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

func (c *Consul) makeServiceRequest() []CatalogService {
	var result []CatalogService
	cat := c.client.Catalog()

	services, _, err := cat.Services(nil)
	if err != nil {
		log.Fatal(err)
	}
	for i := range services {
		for n := range services[i] {
			fmt.Println(services[i][n])
			// service, _, err := cat.Service(services[i], "", nil)
		}

		// fmt.Println(services[i].ServicePort)
		// fmt.Println(services[i].ServiceID)

		// var service CatalogService
		// service.ServiceID = services[i].ServiceID
		// service.ServiceName = services[i].ServiceName
		// service.ServiceAddress = services[i].ServiceAddress
		// service.ServicePort = services[i].ServicePort
		// service.ServiceTags = services[i].ServiceTags
		// service.ServiceEnableTagOverride = services[i].ServiceEnableTagOverride
		// service.Address = services[i].Address
		// service.CreateIndex = services[i].CreateIndex
		// service.Datacenter = services[i].Datacenter
		// service.ID = services[i].ID
		// service.ModifyIndex = services[i].ModifyIndex
		// service.Node = services[i].Node
		// service.NodeMeta = services[i].NodeMeta
		// service.TaggedAddresses = services[i].TaggedAddresses

		// result = append(result, service)
	}
	return result
}

func (service *CatalogService) printService() string {
	return fmt.Sprintf("Id: %s; Address: %s; Name: %s; Port: %d; Tags: %v",
		service.ServiceID,
		service.ServiceAddress,
		service.ServiceName,
		service.ServicePort,
		service.ServicePort)
}

func (c *Consul) makeACLRequest() []ACL {
	var result []ACL
	acl := c.client.ACL()
	acls, _, err := acl.List(nil)
	if err != nil {
		log.Fatal(err)
	}
	for i := range acls {
		var acl ACL

		acl.CreateIndex = acls[i].CreateIndex
		acl.ID = acls[i].ID
		acl.ModifyIndex = acls[i].ModifyIndex
		acl.Name = acls[i].Name
		acl.Rules = acls[i].Rules
		acl.Type = acls[i].Type

		result = append(result, acl)
	}
	return result
}

func (acl *ACL) printACL() string {
	return fmt.Sprintf("Id: %s; Name: %s; Rules: %s; Type: %s; CreateIndex: %d; ModifyIndex: %d",
		acl.ID,
		acl.Name,
		acl.Rules,
		acl.Type,
		acl.CreateIndex,
		acl.ModifyIndex)
}

func (c *Consul) makeNodeRequest() []Node {
	var result []Node
	cat := c.client.Catalog()
	nodes, _, err := cat.Nodes(nil)
	if err != nil {
		log.Fatal(err)
	}
	for i := range nodes {
		var node Node
		node.ID = nodes[i].ID
		node.Address = nodes[i].Address
		node.CreateIndex = nodes[i].CreateIndex
		node.Datacenter = nodes[i].Datacenter
		node.Meta = nodes[i].Meta
		node.ModifyIndex = nodes[i].ModifyIndex
		node.Node = nodes[i].Node
		node.TaggedAddresses = nodes[i].TaggedAddresses

		result = append(result, node)
	}
	return result
}

func (node *Node) printNode() string {
	return fmt.Sprintf("Id: %s; Address: %s; Datacenter: %s; CreateIndex: %d; Meta: %v; ModifyIndex: %d; Node: %s; TaggedAddresses: %v",
		node.ID,
		node.Address,
		node.Datacenter,
		node.CreateIndex,
		node.Meta,
		node.ModifyIndex,
		node.Node,
		node.TaggedAddresses)
}
