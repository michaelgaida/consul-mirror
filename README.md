Consul Mirror
===============
[![Build Status](https://travis-ci.org/michaelgaida/consul-mirror.svg?branch=master)](https://travis-ci.org/michaelgaida/consul-mirror)
[![Coverage Status](https://coveralls.io/repos/github/michaelgaida/consul-mirror/badge.svg?branch=master)](https://coveralls.io/github/michaelgaida/consul-mirror?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/michaelgaida/consul-mirror)](https://goreportcard.com/report/github.com/michaelgaida/consul-mirror)
[![GoDoc](https://godoc.org/github.com/michaelgaida/consul-mirror?status.svg)](https://godoc.org/github.com/michaelgaida/consul-mirror)


consul-mirror is a tool to mirror a single or multi datacenter setup of consul KV-Store into a database. Consul with the KV Store and service discovery is a great tool to bring your setup to a next scalable and flexible level. Starting with consul this tool can help you to have more visability of your Store as well as give you some more fallbacks.

## Why another tool

consul-template provides features that I could not get with `consul-replicate`, `consul kv export` or `consul snapshot`:

* **Export in Database** - Having a copy of the consul KV Store in a Database enables you to inverstigate on key structures and query on values. This can create more visability if the structure is correct or correct used. Investigation on a table where you can query every field.

* **Import from Database** - This tool makes it easy to export your DB KV store into consul. This is specially for the transition time very helpful.

* **Easy partial copies** - For developers which want to have a selected subset of the KV store in a DB or even exported and imported in their local consul dev setup. Datacenter can be exported but ignored doing the import, that will merge the multidatacenter KV storage into one storage with an optional key prefix of the original datacenter.

* **Versioning** - Mirrors can be incremental and versioned. KV Store changes will not overwrite old values but create a new entry with increased version.

* **Backup** - Because KV are versioned and timestamped you can rollback your consul data to any given time or version in the past.

* **Fallback** - Specially in the initial time introducing Consul it is nice to have a fallback to a technologie which is already established in your setup.


## DB
### Engine
Currently supported is MSSQL. Feel free to support more engines and get them in with a pull request

### Structure
```
CREATE TABLE [kv] ( 
	[id] INT IDENTITY ( 1, 1 )  NOT NULL, 
	[timestamp] DATETIME NOT NULL, 
	[createindex] NUMERIC( 20 ) NULL, 
	[flags] NUMERIC( 20 ) NULL, 
	[kvkey] VARCHAR( 255 ) NOT NULL, 
	[kvvalue] VARCHAR( max ) NOT NULL, 
	[lockindex] NUMERIC( 20 ) NULL, 
	[modifyindex] NUMERIC( 20 ) NULL, 
	[regex] VARCHAR( 255 ) NULL, 
	[session] VARCHAR( 255 ) NULL, 
	[version] INT NOT NULL, 
	[acl] VARCHAR( 50 ) NULL, 
	[datacenter] VARCHAR( 50 ) NULL, 
	[deleted] BIT DEFAULT '0' NOT NULL )
GO;
```

## Contributing
To build and install consul-mirror locally clone the repository:

```shell
$ git clone git@github.com:michaelgaida/consul-mirror.git
```

To compile the `consul-mirror` binary for your local machine:

```shell
$ make
```

This will compile the `consul-mirror` binary into `bin/consul-mirror` as
well as your `$GOPATH` and run the test suite.

If you just want to run the tests:

```shell
$ make test
```

Or to run get the code coverage:

```shell
make cov
```