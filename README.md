# Fort

Yet another distributed key-value database

Fort is a distributed key-value store, with a focus on:

* heterogeneous nodes/workloads aware segments placement
* live segments split/merge/migration with zero downtime

Potentially:
* Add a transaction layer to support distributed transaction.

### Technologies used:
* [Etcd](https://github.com/etcd-io/etcd) for storing metadata and notification
* [FlatBuffers](https://github.com/google/flatbuffers) for data serialization
* [Dragonboat](https://github.com/lni/dragonboat) for multi-group Raft
* [bbolt](https://github.com/etcd-io/bbolt) as storage driver for read-heavy segments
* [badger](https://github.com/dgraph-io/badger) as storage driver for write-heavy segments