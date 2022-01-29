# Fort

Yet another distributed transactional key-value database

Fort is a distributed transactional key-value store, with a focus on:

* cluster level distributed transaction
* heterogeneous node aware data placement

### Technologies used:
* [Etcd](https://github.com/etcd-io/etcd) for storing metadata and notification
* [FlatBuffers](https://github.com/google/flatbuffers) for data serialization
* [Dragonboat](https://github.com/lni/dragonboat) for multi-group Raft
* [bbolt](https://github.com/etcd-io/bbolt) as storage driver for read-heavy segments
* [badger](https://github.com/dgraph-io/badger) as storage driver for write-heavy segments