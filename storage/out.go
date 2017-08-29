package storage

import (
	"log"

	"github.com/michaelgaida/consul-mirror/consul"
)

// GetKVs reads KVs from the storage and returns all with the highest version
func (db *Mssql) GetKVs() ([]consul.KV, error) {
	var kv = consul.KV{}
	var result = []consul.KV{}

	// Get all keys with the highest version
	rows, err := db.conn.Query("select DISTINCT a.flags, a.kvkey, a.lockindex, a.modifyindex, a.regex, a.session, a.kvvalue, a.datacenter from [consul].[dbo].[kv] a left outer join [consul].[dbo].[kv] b on a.datacenter = b.datacenter AND a.kvkey = b.kvkey AND a.version < b.version where b.kvkey is NULL")
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		err = rows.Scan(&kv.Flags, &kv.Key, &kv.LockIndex, &kv.ModifyIndex, &kv.Regex, &kv.Session, &kv.Value, &kv.Datacenter)
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, kv)
	}
	return result, nil
}

// GetKV reads KVs from the storage and returns all with the highest version
func (db *Mssql) getKV(key, dc string, version int) (consul.KV, error) {
	var kv = consul.KV{}

	// Get kv with the highest version
	rows, err := db.conn.Query("select flags, kvkey, lockindex, modifyindex, regex, session, kvvalue, datacenter from [consul].[dbo].[kv] where kvkey = ? AND datacenter = ? AND version = ?", key, dc, version)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		err = rows.Scan(&kv.Flags, &kv.Key, &kv.LockIndex, &kv.ModifyIndex, &kv.Regex, &kv.Session, &kv.Value, &kv.Datacenter)
		if err != nil {
			log.Fatal(err)
		}
	}
	return kv, nil
}

// getLatestVersion returns the highest version for a key in a dc
func (db *Mssql) getLatestVersion(key, dc string) int {
	v := 0
	version, err := db.conn.Prepare("select ISNULL(MAX(version), 0) from kv where kvkey = ? and datacenter = ?")
	defer version.Close()
	if err != nil {
		log.Fatal("Prepare statement for get highest version failed: ", err.Error())
	}

	versionres, err := version.Query(key, dc)
	if err != nil {
		log.Fatal("Get highest version failed: ", err.Error())
	}
	for versionres.Next() {
		err := versionres.Scan(&v)

		if err != nil {
			log.Fatal("Scan highest version failed: ", err.Error())
		}
	}
	return v
}
