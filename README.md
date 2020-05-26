# Bureka

**:warning: WORK IN PROGRESS! :warning:**

A libp2p compatible implementation of the [Pastry DHT](http://rowstron.azurewebsites.net/PAST/pastry.pdf) in go.

Usage with libp2p:

```go
func main() {
    writer := node.NewWriter()
    d := dht.New(id, writer)
    
    n := node.New(d, host.Host, writer)
}
```