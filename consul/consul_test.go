package consul

import (
	"reflect"
	"testing"

	consulapi "github.com/hashicorp/consul/api"
)

func TestConsul_makeACLRequest(t *testing.T) {
	type fields struct {
		client *consulapi.Client
		debug  bool
	}
	tests := []struct {
		name   string
		fields fields
		want   []ACL
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Consul{
				client: tt.fields.client,
				debug:  tt.fields.debug,
			}
			if got := c.makeACLRequest(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Consul.makeACLRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
