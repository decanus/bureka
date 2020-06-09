# Bureka

**:warning: WORK IN PROGRESS! :warning:**

An implementation of the [Pastry DHT](http://rowstron.azurewebsites.net/PAST/pastry.pdf) in go. This package includes a libp2p compatible node making it easy to use in a libp2p network.

Usage with libp2p:

```go
func main() {
    d := dht.New(id)
    
    n := node.New(ctx.Background(), d, host.Host)
}
```
