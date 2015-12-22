# ovo

OVO is an In-Memory Distributed Cache and a Key/Value Storage.

## Main features

OVO is a distributed in-memory cache that supports data sharding on multiple instances and data replication.
OVO offers these features:
- Multi-Master cluster architecture
- The nodes can be added and removed without stopping the cluster
- Data are replicated on many nodes if the cluster is configured for replication (Twin nodes)
- Keys are strings but every kind of data values can be stored (JSON documents, XML, images, byte arrays, ...)
- Auto-expiration, data will be automatically removed from the storage if the TTL of the object is setted
- Atomic counters
- OVO supports data sharding on many cluster nodes using smart clients

The project is under development.

## Building OVO

```bash
$ go build -i github.com/maxzerbini/ovo
```

## Starting OVO
### Start a single node
```bash
$ ovo
```
### Start a three node cluster
```bash
$ ovo -conf=./conf/serverconf.json

$ ovo -conf=./conf/serverconf2.json

$ ovo -conf=./conf/serverconf3.json
```

## RESTful API
Clients can connect OVO using RESTful API. 

The available API set includes these endpoints:
- _GET /ovo/keystorage_ gives the count of all the stored keys
- _GET /ovo/keys gives_ the list of all the stored keys
- _GET /ovo/keystorage/:key_ retrieves the object corresponding to key 
- _POST /ovo/keystorage_ puts the body object in the storage
- _PUT /ovo/keystorage_ same as POST
- _DELETE /ovo/keystorage/:key_ removes the object from the storage
- _GET /ovo/keystorage/:key/getandremove_ gets the object and removes it from the storage
- _POST /ovo/keystorage/:key/updatevalueifequal_ updates the object with a new value if the input old value is equal to the stored object value 
- _POST /ovo/keystorage/:key/updatekeyvalueifequal_ updates the object end the key with a new values if the input old value is equal to the stored object value
- _POST /ovo/keystorage/:key/updatekey_ changes the key of an object 
- _GET /ovo/cluster_ gets the cluster topology
- _GET /ovo/cluster/me_ gets the node details
- _POST /ovo/counters_ sets the value of the counter
- _PUT /ovo/counters_ increments (or decrements) the value of the counter
- _GET /ovo/counters/:key_ gets the value of the counter

The Go OVO Client can connect a cluster of OVO nodes. The Go client source code can be found here https://github.com/maxzerbini/ovoclient .