package consul

import (
	"fmt"
	"log"
)

func (c *Consul) makeKVRequest(key string) []KV {
	var result []KV
	kv := c.client.KV()

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

func (c *Consul) makeACLRequest() []ACL {
	var result []ACL

	// c.client.ACL().

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
