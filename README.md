# ovo

OVO is an In-Memory Key/Value Storage.

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

## Restful API
Clients can connect OVO using RESTful API. 

The available API set includes these endpoints:
- GET /ovo/keystorage gives the count of all the stored keys
- GET /ovo/keys gives the list of all the stored keys
- GET /ovo/keystorage/:key retrieves the object corresponding to key 
- POST /ovo/keystorage puts the body object in the storage
- PUT /ovo/keystorage same as POST
- DELETE /ovo/keystorage/:key removes the object from the storage
- GET /ovo/keystorage/:key/getandremove gets the object and removes it from the storage
- POST /ovo/keystorage/:key/updatevalueifequal updates the object with a new value if the input old value is equal to the stored object value 
- POST /ovo/keystorage/:key/updatekeyvalueifequal updates the object end the key with a new values if the input old value is equal to the stored object value
- POST /ovo/keystorage/:key/updatekey changes the key of an object 
- GET /ovo/cluster gets the cluster topology
- GET /ovo/cluster/me gets the node details

The Go OVO Client can connect a cluster of OVO nodes. The Go client source code can be found here https://github.com/maxzerbini/ovoclient .