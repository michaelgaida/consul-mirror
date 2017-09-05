package storage

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/michaelgaida/consul-mirror/consul"
)

// WriteKVs writes a KV array to a MSSQL table
// TODO: Only write if value (in dc) changed
func (db *Mssql) WriteKVs(kvs []consul.KV, ignoreDCs, incremental bool) error {

	insert, err := db.conn.Prepare("insert into kv (timestamp, createIndex, flags, kvkey, lockindex, modifyindex, regex, session, kvvalue, version, datacenter) values (?,?,?,?,?,?,?,?,?,?,?)")
	defer insert.Close()
	if err != nil {
		log.Error("Prepare stmt failed: ", err.Error())
		return err
	}

	for i := range kvs {

		dc := ""
		if !ignoreDCs {
			dc = kvs[i].Datacenter
		}

		v, err := db.getLatestVersion(kvs[i].Key, dc)
		if err != nil {
			log.Error("Error getting latest version: ", err.Error())
			return err
		}

		// If incremental is true we only want to write a entry if the key for dc changed
		modified, err := db.kvIsModified(kvs[i], v)
		if err != nil {
			log.Error("Error checking kv is modiefied: ", err.Error())
			return err
		}
		if (!incremental) || (modified) {
			v++

			log.Infof("Write KV %s=%s (version %d)\n", kvs[i].Key, kvs[i].Value, v)

			res, err := insert.Exec(
				time.Now(),
				kvs[i].CreateIndex,
				kvs[i].Flags,
				kvs[i].Key,
				kvs[i].LockIndex,
				kvs[i].ModifyIndex,
				kvs[i].Regex,
				kvs[i].Session,
				kvs[i].Value,
				v,
				dc)
			if err != nil {
				log.Error("Exec into DB failed: ", err.Error())
				return err

			}
			lastID, err := res.LastInsertId()
			if err != nil {
				log.Error(err)
				return err

			}
			rowCnt, err := res.RowsAffected()
			if err != nil {
				log.Error(err)
				return err
			}
			log.Debugf("ID = %d, affected = %d\n", lastID, rowCnt)
		}
	}
	return nil
}

func (db *Mssql) kvIsModified(kv consul.KV, version int) (bool, error) {
	dbkv, err := db.getKV(kv.Key, kv.Datacenter, version)
	if err != nil {
		log.Error(err)
		return false, err
	}
	return !kv.Equals(*dbkv), nil
}

// func (db *Mssql) writeACLs(acls []consul.ACL) {
// 	prep, err := db.conn.Prepare("inset into acl (CreateIndex, ID, ModifyIndex, Name, Rules, Type) values (?,?,?,?,?,?)")
// 	if err != nil {
// 		log.Fatal("Prepare stmt failed: ", err.Error())
// 	}
// 	for i := range acls {
// 		res, err := prep.Exec(
// 			acls[i].CreateIndex,
// 			acls[i].ID,
// 			acls[i].ModifyIndex,
// 			acls[i].Name,
// 			acls[i].Rules,
// 			acls[i].Type)
// 		if err != nil {
// 			log.Fatal("Exec into DB failed: ", err.Error())
// 		}

// 		lastID, err := res.LastInsertId()
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		rowCnt, err := res.RowsAffected()
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		log.Debug("ID = %d, affected = %d\n", lastID, rowCnt)

// 	}
// }

// func (db *Mssql) writeServices(services []consul.CatalogService) {
// 	prep, err := db.conn.Prepare("inset into catalog_service (Address, CreateIndex, Datacenter, ID, ModifyIndex, Node, NodeMeta, ServiceAddress, ServiceEnableTagOverride, ServiceID, ServiceName, ServicePort, ServiceTags, TaggedAddresses) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
// 	if err != nil {
// 		log.Fatal("Prepare stmt failed: ", err.Error())
// 	}
// 	for i := range services {
// 		res, err := prep.Exec(
// 			services[i].Address,
// 			services[i].CreateIndex,
// 			services[i].Datacenter,
// 			services[i].ID,
// 			services[i].ModifyIndex,
// 			services[i].Node,
// 			services[i].NodeMeta,
// 			services[i].ServiceAddress,
// 			services[i].ServiceEnableTagOverride,
// 			services[i].ServiceID,
// 			services[i].ServiceName,
// 			services[i].ServicePort,
// 			services[i].ServiceTags,
// 			services[i].TaggedAddresses)
// 		if err != nil {
// 			log.Fatal("Exec into DB failed: ", err.Error())
// 		}

// 		lastID, err := res.LastInsertId()
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		rowCnt, err := res.RowsAffected()
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		log.Debug("ID = %d, affected = %d\n", lastID, rowCnt)
// 	}
// }

// func (db *Mssql) writeNodes(nodes []consul.Node) {

// 	prep, err := db.conn.Prepare("inset into catalog_service (Address, CreateIndex, Datacenter, ID, Meta, ModifyIndex, Node, TaggedAddresses) values (?,?,?,?,?,?,?,?)")
// 	if err != nil {
// 		log.Fatal("Prepare stmt failed: ", err.Error())
// 	}
// 	for i := range nodes {
// 		res, err := prep.Exec(
// 			nodes[i].Address,
// 			nodes[i].CreateIndex,
// 			nodes[i].Datacenter,
// 			nodes[i].ID,
// 			nodes[i].Meta,
// 			nodes[i].ModifyIndex,
// 			nodes[i].Node,
// 			nodes[i].TaggedAddresses)
// 		if err != nil {
// 			log.Fatal("Exec into DB failed: ", err.Error())
// 		}
// 		lastID, err := res.LastInsertId()
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		rowCnt, err := res.RowsAffected()
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		log.Debug("ID = %d, affected = %d\n", lastID, rowCnt)
// 	}
// }
