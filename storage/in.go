package storage

import (
	"log"
	"time"

	"github.com/michaelgaida/consul-mirror/consul"
)

func (db *Mssql) WriteKVs(kvs []consul.KV) {
	prep, err := db.conn.Prepare("insert into kv (timestamp, createIndex, flags, kvkey, lockindex, modifyIndex, regex, session, kvvalue, version) values (?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		log.Fatal("Prepare stmt failed: ", err.Error())
	}
	for i := range kvs {
		res, err := prep.Exec(
			time.Now(),
			kvs[i].CreateIndex,
			kvs[i].Flags,
			kvs[i].Key,
			kvs[i].LockIndex,
			kvs[i].ModifyIndex,
			kvs[i].Regex,
			kvs[i].Session,
			kvs[i].Value,
			0)
		// For the version we should gather the old version if available and check if the value changes
		if err != nil {
			log.Fatal("Exec into DB failed: ", err.Error())
		}
		if db.debug {
			lastID, err := res.LastInsertId()
			if err != nil {
				log.Fatal(err)
			}
			rowCnt, err := res.RowsAffected()
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("ID = %d, affected = %d\n", lastID, rowCnt)
		}
	}
}

func (db *Mssql) writeACLs(acls []consul.ACL) {
	prep, err := db.conn.Prepare("inset into acl (CreateIndex, ID, ModifyIndex, Name, Rules, Type) values (?,?,?,?,?,?)")
	if err != nil {
		log.Fatal("Prepare stmt failed: ", err.Error())
	}
	for i := range acls {
		res, err := prep.Exec(
			acls[i].CreateIndex,
			acls[i].ID,
			acls[i].ModifyIndex,
			acls[i].Name,
			acls[i].Rules,
			acls[i].Type)
		if err != nil {
			log.Fatal("Exec into DB failed: ", err.Error())
		}
		if db.debug {
			lastID, err := res.LastInsertId()
			if err != nil {
				log.Fatal(err)
			}
			rowCnt, err := res.RowsAffected()
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("ID = %d, affected = %d\n", lastID, rowCnt)
		}
	}
}

func (db *Mssql) writeServices(services []consul.CatalogService) {
	prep, err := db.conn.Prepare("inset into catalog_service (Address, CreateIndex, Datacenter, ID, ModifyIndex, Node, NodeMeta, ServiceAddress, ServiceEnableTagOverride, ServiceID, ServiceName, ServicePort, ServiceTags, TaggedAddresses) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		log.Fatal("Prepare stmt failed: ", err.Error())
	}
	for i := range services {
		res, err := prep.Exec(
			services[i].Address,
			services[i].CreateIndex,
			services[i].Datacenter,
			services[i].ID,
			services[i].ModifyIndex,
			services[i].Node,
			services[i].NodeMeta,
			services[i].ServiceAddress,
			services[i].ServiceEnableTagOverride,
			services[i].ServiceID,
			services[i].ServiceName,
			services[i].ServicePort,
			services[i].ServiceTags,
			services[i].TaggedAddresses)
		if err != nil {
			log.Fatal("Exec into DB failed: ", err.Error())
		}
		if db.debug {
			lastID, err := res.LastInsertId()
			if err != nil {
				log.Fatal(err)
			}
			rowCnt, err := res.RowsAffected()
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("ID = %d, affected = %d\n", lastID, rowCnt)
		}
	}
}

func (db *Mssql) writeNodes(nodes []consul.Node) {

	prep, err := db.conn.Prepare("inset into catalog_service (Address, CreateIndex, Datacenter, ID, Meta, ModifyIndex, Node, TaggedAddresses) values (?,?,?,?,?,?,?,?)")
	if err != nil {
		log.Fatal("Prepare stmt failed: ", err.Error())
	}
	for i := range nodes {
		res, err := prep.Exec(
			nodes[i].Address,
			nodes[i].CreateIndex,
			nodes[i].Datacenter,
			nodes[i].ID,
			nodes[i].Meta,
			nodes[i].ModifyIndex,
			nodes[i].Node,
			nodes[i].TaggedAddresses)
		if err != nil {
			log.Fatal("Exec into DB failed: ", err.Error())
		}
		if db.debug {
			lastID, err := res.LastInsertId()
			if err != nil {
				log.Fatal(err)
			}
			rowCnt, err := res.RowsAffected()
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("ID = %d, affected = %d\n", lastID, rowCnt)
		}
	}
}
