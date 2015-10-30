# ovo

OVO is an In-Memory Key/Value Storage

## Main features
- Multi-Master Cluster architecture
- The nodes can be added and removed dynamically
- Data's replications are done on many nodes (Twins) giving data's high availability
- Keys are strings but every kind of data values can be stored
- Auto-expiration if the TTL of the object is setted
- OVO supports data sharding on many cluster nodes using smart clients

The project is under development.

## Building OVO

```bash
$ go build -i github.com/maxzerbini/ovo
```

## Starting OVO
```bash
$ ovo -conf=./conf/serverconf.json
```

