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
	rows, err := db.conn.Query("select DISTINCT a.flags, a.kvkey, a.lockindex, a.modifyindex, a.regex, a.session, a.kvvalue from [consul].[dbo].[kv] a left outer join [consul].[dbo].[kv] b on a.datacenter = b.datacenter AND a.kvkey = b.kvkey AND a.version < b.version where b.kvkey is NULL")
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		err = rows.Scan(&kv.Flags, &kv.Key, &kv.LockIndex, &kv.ModifyIndex, &kv.Regex, &kv.Session, &kv.Value)
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, kv)
	}
	return result, nil
}
