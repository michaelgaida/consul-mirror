package consul

import consulapi "github.com/hashicorp/consul/api"

// Consul represents a consul connection
type Consul struct {
	client *consulapi.Client
	debug  bool
}

type KV struct {
	// Key is the name of the key. It is also part of the URL path when accessed
	// via the API.
	Key string

	// CreateIndex holds the index corresponding the creation of this KVPair. This
	// is a read-only field.
	CreateIndex uint64

	// ModifyIndex is used for the Check-And-Set operations and can also be fed
	// back into the WaitIndex of the QueryOptions in order to perform blocking
	// queries.
	ModifyIndex uint64

	// LockIndex holds the index corresponding to a lock on this key, if any. This
	// is a read-only field.
	LockIndex uint64

	// Flags are any user-defined flags on the key. It is up to the implementer
	// to check these values, since Consul does not treat them specially.
	Flags uint64

	// Value is the value for the key. This can be any value, but it will be
	// base64 encoded upon transport.
	Value []byte

	// Session is a string representing the ID of the session. Any other
	// interactions with this key over the same session must specify the same
	// session ID.
	Session string

	// Not implemented yet in consul API
	Regex string
}

type Node struct {
	ID              string
	Node            string
	Address         string
	Datacenter      string
	TaggedAddresses map[string]string
	Meta            map[string]string
	CreateIndex     uint64
	ModifyIndex     uint64
}

type CatalogService struct {
	ID                       string
	Node                     string
	Address                  string
	Datacenter               string
	TaggedAddresses          map[string]string
	NodeMeta                 map[string]string
	ServiceID                string
	ServiceName              string
	ServiceAddress           string
	ServiceTags              []string
	ServicePort              int
	ServiceEnableTagOverride bool
	CreateIndex              uint64
	ModifyIndex              uint64
}

type ACL struct {
	CreateIndex uint64
	ID          string
	ModifyIndex uint64
	Name        string
	Rules       string
	Type        string
}
