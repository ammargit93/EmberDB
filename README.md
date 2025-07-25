# 🔥 EmberDB - A Distributed In-Memory Key-Value Store

**EmberDB** is a lightweight, distributed, in-memory key-value database with basic replication and file and data caching. Designed for educational and experimental use, EmberDB showcases how to build a distributed system with leader election, peer registration, and data replication over TCP and HTTP.

---

## 🧠 Features

- 🧩 **Key-Value Store:** Simple `SET`, `GET`, `DEL`, `UPDATE` commands.
- 📁 **File Caching Support:** Upload and retrieve files using `SETFILE`, `GETFILE`, and `SAVE` commands.
- ⚡ **TCP Command Interface:** Interact with nodes using a custom TCP protocol.
- 📡 **Registry Server:** Maintains peer information and assigns the first node as leader.
- 🔁 **Leader Replication:** Only the leader node replicates `SET` operations to followers.
- 🌐 **HTTP SDK Gateway:** HTTP interface for clients (Go/Python supported).

---

### Setup

Start the Registry Server first
The registry server outputs the node ip currently stored. The first one is the leader.

```bash
cd Registry/
go run .
```

Start the Cache Server.

```bash
cd cmd/emberdb/
go run . :1010      #leader port
```

Connect to the Cache Server
```bash
cd client
go run . :1010      #leader port
```
A shell will open up.
```bash
ember> SET a 10
SET OK
ember> GET a
10
```

### SDK Support
to use the sdk you will have to start the api at /sdk/go-client/sdk_api.go

From then on you can simply follow the examples mentioned for go and Python.


This is just an MVP.