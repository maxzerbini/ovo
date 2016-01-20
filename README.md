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
## Node configuration

### The configuration file
The configuration file _serverconf.json_ is a JSON file that defines the addresses and ports used by OVO to listen for HTTP calls and cluster communications. The configuration file defines also other configurations parameters used by the server node.
These are the all the configuration parameters:
- *Name* is the unique node name, if omitted the node will generate a random one
- *Host* is the hostname or IP address of the HTTP listener
- *Port* is the port of the HTTP listener
- *APIHost* is the hostname or IP address used for inter-cluster communications
- *APIPort* is the port used for inter-cluster communications
- *Twins* is a list of node names of the cluster, the twins are the nodes used by the server to replicate its data
- *Stepbrothers* is a list of node names of the cluster, stepbrothers are the nodes to which the server requests to become a replica
- *Debug* is a flag that enables internal logging

This is a configuration file example
```JSON
{
	"ServerNode":
	{
		"Node":
		{
			"Name":"mizard",
			"Host":"192.168.1.102",
			"Port":5050,
			"APIHost":"192.168.1.102",
			"APIPort":5052
		},
		"Twins":[],
		"Stepbrothers":[]
	},
	"Debug":true
}
```
### Cluster configuration
OVO cluster can be formed by two or more nodes. Nodes can be added or removed without stopping the cluster activities. 
We must configure the node that is added to a cluster so that I can see at least another active node. This is done by providing a description (maybe partial) of the topology.
This sample configuration allows us to create a cluster formed by two nodes *mizard* and  *righel* and in which one is the twin of the other.
```JSON
{
  "ServerNode": {
    "Node": {
      "Name": "righel",
      "Host": "192.168.1.103",
      "Port": 5050,
      "APIHost": "192.168.1.103",
      "APIPort": 5052
    },
    "Twins": ["mizard"],
    "Stepbrothers": ["mizard"]
  },
  "Topology": {
    "Nodes": [
      {
        "Node": {
          "Name": "mizard",
          "Host": "192.168.1.102",
          "Port": 5050,
          "APIHost": "192.168.1.102",
          "APIPort": 5052
        }
      }
    ]
  },
  "Debug": true
}
```

### The temporary configuration file
Every time that the server starts or every time that the cluster topology changes the temporary configuration file is updated and saved.
The temporary configuration file resides in the same folder of the configuration file and has the same name but its extension is .temp .

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
- _DELETE /ovo/counters/:key_ delete the counter

## Client libraries

### Go client library
The Go OVO Client can connect a cluster of OVO nodes. The Go client source code can be found here https://github.com/maxzerbini/ovoclient .

### .Net cleint library
The .Net OVO Client can connect a cluster of OVO nodes and offers the same API of the Go client. The .Net client source code can be found here https://github.com/maxzerbini/ovodotnet .
The library can by downloaded via *Nuget.org* at https://www.nuget.org/packages/OVOdotNetClient/ or using the Nuget Package Manager.
```
PM> Install-Package OVOdotNetClient
```

### Java client library
The Java client library is under development.
The source code is at https://github.com/maxzerbini/ovojava .