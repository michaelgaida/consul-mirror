package consul

import "fmt"

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

func (service *CatalogService) printService() string {
	return fmt.Sprintf("Id: %s; Address: %s; Name: %s; Port: %d; Tags: %v",
		service.ServiceID,
		service.ServiceAddress,
		service.ServiceName,
		service.ServicePort,
		service.ServicePort)
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
