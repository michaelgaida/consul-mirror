package consul

// import (
// 	"reflect"
// 	"testing"

// 	"github.com/golang/mock/gomock"
// 	consulapi "github.com/hashicorp/consul/api"
// )

// func TestConsul_makeACLRequest(t *testing.T) {
// 	type ACL struct {
// 		CreateIndex uint64
// 		ID          string
// 		ModifyIndex uint64
// 		Name        string
// 		Rules       string
// 		Type        string
// 	}

// 	type consul struct {
// 		client *consulapi.Client
// 		debug  bool
// 	}
// 	type testsstruct struct {
// 		name   string
// 		fields consul
// 		want   []ACL
// 	}

// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	var client = consulapi.Client{}
// 	con := consul{&client, false}
// 	acl := []ACL{
// 		{1, "I am an ID", 2, "law", "dont cross red light", "schoen"},
// 	}

// 	tests := []testsstruct{
// 		{"First test", con, acl},
// 	}

// 	for _, tt := range tests {

// 		t.Run(tt.name, func(t *testing.T) {
// 			c := &Consul{
// 				client: tt.fields.client,
// 				debug:  tt.fields.debug,
// 			}
// 			mockIndex := c.makeACLRequest.NewMockIndex(ctrl)
// 			mockIndex.EXPECT().ConcreteRet().Return(tt.want)

// 			if got := c.makeACLRequest(); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("Consul.makeACLRequest() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestConsul_GetKVs(t *testing.T) {
// 	type fields struct {
// 		client *consulapi.Client
// 		debug  bool
// 	}
// 	type args struct {
// 		key string
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		args   args
// 		want   []KV
// 	}{
// 	// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			c := &Consul{
// 				client: tt.fields.client,
// 				debug:  tt.fields.debug,
// 			}
// 			if got := c.GetKVs(tt.args.key); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("Consul.GetKVs() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
