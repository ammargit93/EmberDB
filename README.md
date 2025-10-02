# EmberDB

EmberDB is a minimal distributed key-value (KV) store implemented in Go. It demonstrates basic distributed system concepts, including leader-follower replication, peer discovery, and a TCP-based command interface.

## 🚀 Features

**🔑 Leader-Follower Replication**: Only the leader can perform SET commands; changes replicate to followers automatically.

**🌐 Peer Discovery**: Nodes register with the registry, which tracks all peers and the leader using SQLite3.

**💻 TCP Client Interface**: Interact with nodes using simple GET and SET commands.

**🔄 HTTP Replication**: Followers expose an HTTP endpoint for leader replication.

**🗄️ Persistent Registry**: Registry stores node info in SQLite3 for reliable discovery.

**📁 File Storage Commands**: SETFILE and GETFILE allow storing and retrieving file contents.

#

**To get started**:
```bash
git clone https://github.com/ammargit93/EmberDB.git
cd EmberDB
```

**Start the registry**
```bash
cd Registry
go run . # starts at localhost:5050
```


**Start the nodes**
```bash
cd cmd\emberbd
go run . :1010  # leader 
```


**Start the client**
```bash
cd client
go run . :1010 # connect with leader
```

```bash
ember> SET a 10
SET OK
ember> GET a
10
```
