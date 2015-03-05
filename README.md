# Benchmark for all GoLang based Key/Value stores and clients.

Based on the idea of https://gist.github.com/karlseguin/0ba24030fb12b10b686b

Available persistent KV stores:

OSX 10.9.5; Processor  2.4 GHz Intel Core i5; Memory  8 GB 1600 MHz DDR3; Late 2013

```
$ ./gokvbench -redis -bolt
main.testRedisWrite	   20000	     67270 ns/op	     136 B/op	       8 allocs/op
main.testRedisRead	   20000	     70634 ns/op	     128 B/op	       9 allocs/op
main.testBoltWrite	    5000	    526057 ns/op	   47425 B/op	      61 allocs/op
main.testBoltRead	  500000	      3837 ns/op	     648 B/op	      13 allocs/op
main.testGkvliteWrite	   50000	     29655 ns/op	     292 B/op	       8 allocs/op
main.testGkvliteRead	 1000000	      2671 ns/op	      24 B/op	       3 allocs/op
main.testDiskvWrite	   10000	    109660 ns/op	     528 B/op	      18 allocs/op
main.testDiskvRead	  100000	     24315 ns/op	    2899 B/op	      21 allocs/op
main.testCznicKvWrite	   10000	    156615 ns/op	    9437 B/op	      34 allocs/op
main.testCznicKvRead	   10000	    372743 ns/op	  149407 B/op	      54 allocs/op
```

More will follow...

Old write test on old MacBook Air:

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

Short discussion on [Twitter](https://twitter.com/schumacherfm/status/573060236166234112)

### The MIT License (MIT)

Copyright (c) 2015 Cyrill Schumacher and contributors

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
