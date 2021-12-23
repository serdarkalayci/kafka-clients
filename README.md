# kafka-clients
Kafka producer and consumer developed in Go.
The producer and cosumer uses [Segment's kafka-go library](https://github.com/segmentio/kafka-go).

The repository contains a docker compose file named [kafka-compose.yml](kafka-compose.yml) which can be used for development purposes.

By default kafka creates a topic automatically when a message arrives. But this topic will be created with 1 partition which means it can be consumed by just one consumer. To try with more than one consumer, the topic must be created with more partitions. To do this first enter interactive mode using

```sh
docker exec -it kafka-container-name sh
```

then running `kafka-topics.sh` 
```sh
cd bin
./kafka-topics.sh --delete --zookeeper my-cluster-zookeeper-client:2181 --topic my-topic
```

also a topic can be deleted by running
```sh
./kafka-topics.sh --create --zookeeper my-cluster-zookeeper-client:2181 --replication-factor 1 --partitions 10 --topic my-topic
```