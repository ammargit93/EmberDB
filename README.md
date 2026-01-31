# EmberDB v2

EmberDB is a minimal key-value (KV) store implemented in Go. It demonstrates basic system concepts, including basic KV API, Periodic snapshot persistence, file handling.

## 🚀 Features

**🔑 Periodic Snapshot**: Writes data to disk every n seconds, specified with the --snapshot flag.

**💻 REST API Interface**: Interact with nodes using simple Set, Get, Update, Delete commands.

**🔄 Minimal SDK**: Go client library to interact with the API.

**🗄️ Server in Fiber**: Server is built using the Fiber web framework.

**📄 Write-Ahead Log**: Writes a sequential log with commands executed.

**📁 File Storage Commands**: Files can be stored as byte arrays.

#

**To get started**:
```bash
git clone https://github.com/ammargit93/EmberDB.git
cd EmberDB
```

**Start the server**
```bash
cd cmd
go run .
```

alternatively
```bash
go run . --snapshot 10s  # saves a snapshot every 10 seconds, fallbacks to 5s if flag not provided.
```

**Start the client**
Connect using postman or curl

```curl
# Set a key
curl -X POST http://localhost:9182/set \
  -H "Content-Type: application/json" \
  -d '{"namespace":"users","key":"username","value":"john doe"}'

# Get a key
curl http://localhost:9182/get/users/username

# Update a key
curl -X PATCH http://localhost:9182/update \
  -H "Content-Type: application/json" \
  -d '{"namespace":"users","key":"username","value":"jane"}'

# Delete a key
curl -X DELETE http://localhost:9182/delete/users/username

```

### Upcoming 
- Failure detection and crash recovery
- Enhanced file handling
- Better client libraries for multiple languages
- Better Data structures for kv store (skiplist)
- Distributed systems behavior