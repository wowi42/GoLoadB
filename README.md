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


Because I like fun, and because I like to code, I did my own load balancer in Go, named [GoLB](https://github.com/wowi42/GoLB)
 Why that ? I didn't find the possibility to have round robin + sticky, and because each time, sticky was with the cookie. So, I did the choice to do it myself.

First, I can't do Direct Routing, so I did my proxy an application level. I used some of my own libraries, to do it. My first idea, was to put in a redis Table, values like this : RemoteAddr:server

But on this server, we didn't had enough ressources, so I did it without redis : everything in RAM. It's greedy, but can help at the beginning. To don't explode the memory, you can store it not like RemoteAddr:server but like this : RemoteIP:server. What is the difference ? RemoteAddr contain port too. Of course, I'm slowler than Nginx, but I respect this :
```configure the aforementioned load balancer to be "sticky" - the same host should hit the same webserver for repeat requests, only switching when a webserver goes down, and not switching back when the webserver goes back up```. GoLB don't switch back when the webserver is back.
There can be a problem, in my current GoLB : I didn't clean the ram, no TTL on keys, so if you use it in production, it can be a problem. But you don't need a lot of time to do it yourself : a general ttl cleaner, in a go routine, which send on a channel some instructions, and when the Handle Function catch this signal, just clean the map. 
