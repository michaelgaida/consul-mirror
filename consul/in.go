package consul

import (
	"fmt"
	"log"

	consulapi "github.com/hashicorp/consul/api"
)

func (c *Consul) writeKV(kvs []KV) error {

	var kv = consulapi.KVPair{}

	for i := range kvs {
		kv.CreateIndex = kvs[i].CreateIndex
		kv.Flags = kvs[i].Flags
		kv.Key = kvs[i].Key
		kv.LockIndex = kvs[i].LockIndex
		kv.ModifyIndex = kvs[i].ModifyIndex
		kv.Session = kvs[i].Session
		kv.Value = kvs[i].Value

		response, err := c.client.KV().Put(&kv, nil)
		if err != nil {
			log.Fatal(err)
		}
		if c.debug {
			log.Printf("Wrote KV [%s : %s] (Duration: %v)", kv.Key, kv.Value, response.RequestTime)
		}
	}

	return nil
}

func (c *Consul) writeACLs(acls []ACL) error {

	var acl = consulapi.ACLEntry{}

	for i := range acls {

		acl.CreateIndex = acls[i].CreateIndex

		dunno, response, err := c.client.ACL().Create(&acl, nil)
		if err != nil {
			log.Fatal(err)
		}
		if c.debug {
			log.Printf("Wrote ACL [%s] (Duration: %v)", acl.Name, response.RequestTime)
		}
		fmt.Printf("dunno +++====+++ %s", dunno)
	}

	return nil
}

func (c *Consul) writeCatalogServices(services []CatalogService) error {
	return nil
}

func (c *Consul) writeNodes(nodes []Node) error {
	return nil
}
