GoLB
====

GoLB is a Go Sticky Round Robin Load balancer.

```
{
    "BackServers": [
        "127.0.0.1:8081",
        "127.0.0.1:8082"
    ],
    "Log": {
        "Folder": "/tmp/golb/"
    },
    "LogColor": true,
    "Name": "edge42",
    "Server": {
        "Hostname": "127.0.0.1",
        "Port": "8080"
    }
}

```

TTL is in second : it's for the expiration of keys in redis
