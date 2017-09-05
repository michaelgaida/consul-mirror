package storage

import (
	log "github.com/sirupsen/logrus"

	"github.com/michaelgaida/consul-mirror/consul"
)

// GetKVs reads KVs from the storage and returns all with the highest version
func (db *Mssql) GetKVs(prefix string) ([]consul.KV, error) {
	var kv = consul.KV{}
	var result = []consul.KV{}

	// Get all keys with the highest version
	rows, err := db.conn.Query("select DISTINCT a.flags, a.kvkey, a.lockindex, a.modifyindex, a.regex, a.session, a.kvvalue, a.datacenter from [consul].[dbo].[kv] a left outer join [consul].[dbo].[kv] b on a.datacenter = b.datacenter AND a.kvkey = b.kvkey AND a.version < b.version where b.kvkey is NULL AND a.kvkey like '" + prefix + "%' ")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&kv.Flags, &kv.Key, &kv.LockIndex, &kv.ModifyIndex, &kv.Regex, &kv.Session, &kv.Value, &kv.Datacenter)
		if err != nil {
			return nil, err
		}
		result = append(result, kv)
	}
	return result, nil
}

// GetKV reads KVs from the storage and returns all with the highest version
func (db *Mssql) getKV(key, dc string, version int) (*consul.KV, error) {
	var kv = consul.KV{}

	// Get kv with the highest version
	rows, err := db.conn.Query("select flags, kvkey, lockindex, modifyindex, regex, session, kvvalue, datacenter from [consul].[dbo].[kv] where kvkey = ? AND datacenter = ? AND version = ?", key, dc, version)
	defer rows.Close()
	if err != nil {
		log.Error("Error querying the KV with highest version", err.Error())
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&kv.Flags, &kv.Key, &kv.LockIndex, &kv.ModifyIndex, &kv.Regex, &kv.Session, &kv.Value, &kv.Datacenter)
		if err != nil {
			log.Error("Error scaning result set with the KV with highest version", err.Error())
			return nil, err
		}
	}
	return &kv, nil
}

// getLatestVersion returns the highest version for a key in a dc
func (db *Mssql) getLatestVersion(key, dc string) (int, error) {
	v := 0
	version, err := db.conn.Prepare("select ISNULL(MAX(version), 0) from kv where kvkey = ? and datacenter = ?")
	defer version.Close()
	if err != nil {
		log.Error("Prepare statement for get highest version failed: ", err.Error())
		return 0, err
	}

	versionres, err := version.Query(key, dc)
	if err != nil {
		log.Error("Get highest version failed: ", err.Error())
		return 0, err

	}
	for versionres.Next() {
		err := versionres.Scan(&v)

		if err != nil {
			log.Error("Scan highest version failed: ", err.Error())
			return 0, err
		}
	}
	return v, nil
}
