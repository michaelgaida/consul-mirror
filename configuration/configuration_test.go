package configuration

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// {
// 	"user": "",
// 	"password": "",
// 	"consul": {
// 		"acl": ""
// 	}
// }

// map[string]interface{}
// db := make(map[string]interface{})
// db["consul"] = make(map[string]interface{})
// db["consul"]["acl"] = ""

type database struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Database string `json:"database"`
	Port     int    `json:"port"`
}

type consul struct {
	ACL  string `json:"acl"`
	Host string `json:"host"`
	Port int    `json:"port"`
}

type fields struct {
	Database database `json:"database"`
	Consul   consul   `json:"consul"`
	LogLevel string   `json:"loglevel"`
}

func TestStruct_ValidateConfiguration(t *testing.T) {
	type teststruct struct {
		name string
		fields
		want error
	}
	dbconfigs := []database{
		{"athur", "Dont't Panic", "Earth", "Towl", 42},
		{"", "Dont't Panic", "Earth", "Towl", 42},
		{"athur", "", "Earth", "Towl", 42},
		{"athur", "Dont't Panic", "", "Towl", 42},
		{"athur", "Dont't Panic", "Earth", "", 42},
		{"athur", "Dont't Panic", "Earth", "Towl", 0},
	}

	consulconfigs := []consul{
		{"Harmless", "Vogons", 23},
		{"", "Vogons", 23},
		{"Harmless", "", 23},
		{"Harmless", "Vogons", 0},
	}

	tests := []teststruct{
		{"All Good", fields{dbconfigs[0], consulconfigs[0], "info"}, nil},
		{"Empty DB User", fields{dbconfigs[1], consulconfigs[0], "debug"}, errors.New("Database user is empty")},
		{"Empty DB Password", fields{dbconfigs[2], consulconfigs[0], "panic"}, errors.New("Database password is empty")},
		{"Empty DB Host", fields{dbconfigs[3], consulconfigs[0], "info"}, errors.New("Database host is empty")},
		{"Empty DB Name", fields{dbconfigs[4], consulconfigs[0], "info"}, errors.New("Database name is empty")},
		{"Empty DB Port", fields{dbconfigs[5], consulconfigs[0], "info"}, errors.New("Database port is 0")},
		{"Empty Consul ACL", fields{dbconfigs[0], consulconfigs[1], "info"}, nil},
		{"Empty Consul Host", fields{dbconfigs[0], consulconfigs[2], "info"}, errors.New("Consul host is empty")},
		{"Empty Consul Port", fields{dbconfigs[0], consulconfigs[3], "info"}, errors.New("Consul port is 0")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Struct{
				Database: tt.fields.Database,
				Consul:   tt.fields.Consul,
				LogLevel: tt.fields.LogLevel,
			}
			got := c.ValidateConfiguration()
			if (got != nil || tt.want != nil) && (got.Error() != tt.want.Error()) {
				t.Errorf("Struct.ValidateConfiguration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStruct_PrintDebug(t *testing.T) {
	type teststruct struct {
		name string
		fields
		want string
	}

	dbconfigs := []database{
		{"athur", "Dont't Panic", "Earth", "Towl", 42},
	}

	consulconfigs := []consul{
		{"Harmless", "Vogons", 23},
	}

	tests := []teststruct{
		{"All Good", fields{dbconfigs[0], consulconfigs[0], "debug"}, "42"},
		{"All Good", fields{dbconfigs[0], consulconfigs[0], "debug"}, "Vogons"},
		{"All Good", fields{dbconfigs[0], consulconfigs[0], "debug"}, "Harmless"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Struct{
				Database: tt.fields.Database,
				Consul:   tt.fields.Consul,
				LogLevel: tt.fields.LogLevel,
			}
			got := c.PrintDebug()
			assert.Contains(t, got, tt.want)
		})
	}
}

func TestGetConfig(t *testing.T) {
	type args struct {
		f string
	}
	tests := []struct {
		name string
		args args
		want *Struct
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetConfig(tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
