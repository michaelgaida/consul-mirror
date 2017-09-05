package consul

import (
	"fmt"
	"strings"

	consulapi "github.com/hashicorp/consul/api"
	log "github.com/sirupsen/logrus"
)

func (c *Consul) makeKVRequest(key string, dcs []DC) ([]KV, error) {
	var result []KV
	kv := c.client.KV()

	for dcindex := range dcs {
		keys, _, err := c.client.KV().Keys(key, ":", &consulapi.QueryOptions{Datacenter: dcs[dcindex].Name})
		if err != nil {
			log.Error(err)
			return nil, err
		}
		for i := range keys {
			if !strings.HasPrefix(keys[i], "vault") {
				if keys[i] == "local/it-whitefalcon" {
					fmt.Print("------")
				}
				log.Debugf("BEFORE Export KV: %s, dc=%s", keys[i], dcs[dcindex].Name)
				kvc, _, err := kv.Get(keys[i], &consulapi.QueryOptions{Datacenter: dcs[dcindex].Name})
				if err != nil {
					log.Error(err)
					return nil, err
				}
				log.Debugf("Export KV: %s=%s, dc=%s", kvc.Key, kvc.Value, dcs[dcindex].Name)
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
				kv.Datacenter = dcs[dcindex].Name

				result = append(result, kv)
			}
		}
	}
	return result, nil
}

func (c *Consul) makeDCRequest(selectedDC string) ([]DC, error) {
	var result []DC
	cat := c.client.Catalog()
	dcs, err := cat.Datacenters()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	for i := range dcs {
		dc := DC{}
		dc.Name = dcs[i]
		if selectedDC == "" || selectedDC == dc.Name {
			result = append(result, dc)
		}
	}
	return result, nil
}

// func (c *Consul) makeServiceRequest() []CatalogService {
// 	var result []CatalogService
// 	cat := c.client.Catalog()
// 	var s = CatalogService{}

// 	services, _, err := cat.Services(nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	for i := range services {
// 		for n := range services[i] {
// 			fmt.Println(services[i][n])

// 			service, _, err := cat.Service(services[i][n], "", nil)
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 			for y := range service {
// 				s.Address = service[y].Address
// 				s.CreateIndex = service[y].CreateIndex
// 				s.Datacenter = service[y].Datacenter
// 				s.ID = service[y].ID
// 				s.ModifyIndex = service[y].ModifyIndex
// 				s.Node = service[y].Node
// 				s.NodeMeta = service[y].NodeMeta
// 				s.ServiceAddress = service[y].ServiceAddress
// 				s.ServiceEnableTagOverride = service[y].ServiceEnableTagOverride
// 				s.ServiceID = service[y].ServiceID
// 				s.ServiceName = service[y].ServiceName
// 				s.ServicePort = service[y].ServicePort
// 				s.ServiceTags = service[y].ServiceTags
// 			}

// 		}

// 		// fmt.Println(services[i].ServicePort)
// 		// fmt.Println(services[i].ServiceID)

// 		// var service CatalogService
// 		// service.ServiceID = services[i].ServiceID
// 		// service.ServiceName = services[i].ServiceName
// 		// service.ServiceAddress = services[i].ServiceAddress
// 		// service.ServicePort = services[i].ServicePort
// 		// service.ServiceTags = services[i].ServiceTags
// 		// service.ServiceEnableTagOverride = services[i].ServiceEnableTagOverride
// 		// service.Address = services[i].Address
// 		// service.CreateIndex = services[i].CreateIndex
// 		// service.Datacenter = services[i].Datacenter
// 		// service.ID = services[i].ID
// 		// service.ModifyIndex = services[i].ModifyIndex
// 		// service.Node = services[i].Node
// 		// service.NodeMeta = services[i].NodeMeta
// 		// service.TaggedAddresses = services[i].TaggedAddresses

// 		// result = append(result, service)
// 	}
// 	return result
// }

// func (c *Consul) makeACLRequest() []ACL {
// 	var result []ACL

// 	// c.client.ACL().

// 	acl := c.client.ACL()
// 	acls, _, err := acl.List(nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	for i := range acls {
// 		acl := ACL{}

// 		acl.CreateIndex = acls[i].CreateIndex
// 		acl.ID = acls[i].ID
// 		acl.ModifyIndex = acls[i].ModifyIndex
// 		acl.Name = acls[i].Name
// 		acl.Rules = acls[i].Rules
// 		acl.Type = acls[i].Type

// 		result = append(result, acl)
// 	}
// 	return result
// }

// func (c *Consul) makeNodeRequest() []Node {
// 	var result []Node
// 	cat := c.client.Catalog()
// 	nodes, _, err := cat.Nodes(nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	for i := range nodes {
// 		node := Node{}

// 		node.ID = nodes[i].ID
// 		node.Address = nodes[i].Address
// 		node.CreateIndex = nodes[i].CreateIndex
// 		node.Datacenter = nodes[i].Datacenter
// 		node.Meta = nodes[i].Meta
// 		node.ModifyIndex = nodes[i].ModifyIndex
// 		node.Node = nodes[i].Node
// 		node.TaggedAddresses = nodes[i].TaggedAddresses

// 		result = append(result, node)
// 	}
// 	return result
// }
