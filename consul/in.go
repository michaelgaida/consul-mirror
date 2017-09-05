package consul

import (
	"bytes"

	consulapi "github.com/hashicorp/consul/api"
	log "github.com/sirupsen/logrus"
)

func generateKey(dcprefix bool, setprefix string, kv KV) string {
	var b bytes.Buffer

	if dcprefix {
		b.WriteString(kv.Datacenter)
		b.WriteString("/")
	}
	if setprefix != "" {
		b.WriteString(setprefix)
		b.WriteString("/")
	}
	b.WriteString(kv.Key)

	return b.String()
}

func (c *Consul) writeKV(kvs []KV, ignoreDC, dcprefix bool, setprefix string) error {

	var kv = consulapi.KVPair{}

	for i := range kvs {

		kv.CreateIndex = kvs[i].CreateIndex
		kv.Flags = kvs[i].Flags
		kv.Key = generateKey(dcprefix, setprefix, kvs[i])
		kv.LockIndex = kvs[i].LockIndex
		kv.ModifyIndex = kvs[i].ModifyIndex
		kv.Session = kvs[i].Session
		kv.Value = kvs[i].Value

		writeopt := &consulapi.WriteOptions{}
		if !ignoreDC {
			writeopt = &consulapi.WriteOptions{Datacenter: kvs[i].Datacenter}
		}

		response, err := c.client.KV().Put(&kv, writeopt)
		if err != nil {
			if err.Error() != "Unexpected response code: 500 (No path to datacenter)" {
				return err
			}
			log.Errorf("%s [%s]", err.Error(), kvs[i].Datacenter)
		} else {
			log.Debugf("Wrote KV [%s : %s] (Duration: %v)", kv.Key, kv.Value, response.RequestTime)
		}
	}
	return nil
}

// func (c *Consul) writeACLs(acls []ACL) error {

// 	var acl = consulapi.ACLEntry{}

// 	for i := range acls {

// 		acl.CreateIndex = acls[i].CreateIndex

// 		_, response, err := c.client.ACL().Create(&acl, nil)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		log.Debug("Wrote ACL [%s] (Duration: %v)", acl.Name, response.RequestTime)
// 	}

// 	return nil
// }

// func (c *Consul) writeCatalogServices(services []CatalogService) error {
// 	// var service = consulapi.CatalogService{}
// 	var reg = consulapi.CatalogRegistration{}
// 	var service = consulapi.AgentService{}

// 	for i := range services {
// 		reg.Address = services[i].Address
// 		reg.Datacenter = services[i].Datacenter
// 		reg.ID = services[i].ID
// 		reg.Node = services[i].Node
// 		reg.NodeMeta = services[i].NodeMeta
// 		reg.TaggedAddresses = services[i].TaggedAddresses

// 		service.Address = services[i].ServiceAddress
// 		// service.CreateIndex = services[i].CreateIndex
// 		reg.Service = &service

// 		c.client.Catalog().Register(&reg, nil)

// 		// s := reg.Service()
// 		// service.Address = services[i].Address

// 		// dunno, response, err := c.client.Catalog().Register()
// 		// if err != nil {
// 		// 	log.Fatal(err)
// 		// }
// 		// if c.debug {
// 		// 	log.Printf("Wrote ACL [%s] (Duration: %v)", acl.Name, response.RequestTime)
// 		// }
// 		// fmt.Printf("dunno +++====+++ %s", dunno)
// 	}

// 	return nil

// }

// func (c *Consul) writeNodes(nodes []Node) error {
// 	var node = consulapi.Node{}

// 	for i := range nodes {
// 		node.Address = nodes[i].Address
// 		node.CreateIndex = nodes[i].CreateIndex
// 		node.Datacenter = nodes[i].Datacenter
// 		node.ID = nodes[i].ID
// 		node.Meta = nodes[i].Meta
// 		node.ModifyIndex = nodes[i].ModifyIndex
// 		node.Node = nodes[i].Node
// 		node.TaggedAddresses = nodes[i].TaggedAddresses

// 		// c.client.Catalog().Node
// 	}
// 	return nil
// }
