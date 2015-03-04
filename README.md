# Benchmark for all GoLang based Key/Value stores and clients.

Based on the idea of https://gist.github.com/karlseguin/0ba24030fb12b10b686b

Available persistent KV stores:

```
$ ./gokvbench -bolt -gkvlite -redis -diskv -cznickv -ledisdb
BoltDB	    3000	    461525 ns/op	   46333 B/op	      61 allocs/op
Redis	   20000	     82473 ns/op	     136 B/op	       8 allocs/op
gkvlite	   50000	     29106 ns/op	     292 B/op	       8 allocs/op
diskv	   10000	    131318 ns/op	     696 B/op	      18 allocs/op
cznickv	   10000	    156007 ns/op	    9432 B/op	      34 allocs/op
ledisdb	  200000	     10370 ns/op	     474 B/op	      11 allocs/op

Done!
```

Feel free to submit your own or improvements to the test via PR.

```
$ go run main.go -h
$ go run main.go -bolt -gkvlite -diskv -cznickv -ledisdb
```

Make sure that Redis runs on `127.0.0.1:6379` and if you test postgres: `localhost:5432`

