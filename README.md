# Benchmark for all GoLang based Key/Value stores and clients.

Based on the idea of https://gist.github.com/karlseguin/0ba24030fb12b10b686b

OSX 10.9.5; Processor  2.4 GHz Intel Core i5; Memory  8 GB 1600 MHz DDR3; Late 2013

```
$ ./gokvbench -all
BenchmarkMySQLWrite	   	 10000	    115209 ns/op	     200 B/op	      14 allocs/op
BenchmarkMySQLRead	   	 10000	    116643 ns/op	     384 B/op	      20 allocs/op
BenchmarkSQLiteWrite	  3000	    509226 ns/op	     184 B/op	      12 allocs/op
BenchmarkSQLiteRead	  	100000	     18932 ns/op	     448 B/op	      21 allocs/op
BenchmarkRedisWrite	   	 30000	     50255 ns/op	     136 B/op	       8 allocs/op
BenchmarkRedisRead	   	 30000	     53008 ns/op	     128 B/op	       9 allocs/op
BenchmarkBoltWrite		  5000	    333430 ns/op	   47422 B/op	      61 allocs/op
BenchmarkBoltRead	  	500000	      3143 ns/op	     648 B/op	      13 allocs/op
BenchmarkGkvliteWrite	100000	     24653 ns/op	     302 B/op	       8 allocs/op
BenchmarkGkvliteRead   1000000	      2037 ns/op	      24 B/op	       3 allocs/op
BenchmarkDiskvWrite	   	 20000	     77249 ns/op	     528 B/op	      18 allocs/op
BenchmarkDiskvRead	  	200000	    139040 ns/op	    2895 B/op	      21 allocs/op
BenchmarkCznicKvWrite    10000	    118868 ns/op	    9436 B/op	      34 allocs/op
BenchmarkCznicKvRead	 10000	    296773 ns/op	  149408 B/op	      54 allocs/op
BenchmarkLedisDbWrite	300000	      6279 ns/op	     423 B/op	      10 allocs/op
BenchmarkLedisDbRead	500000	      7200 ns/op	     801 B/op	      22 allocs/op
BenchmarkLevelDbWrite	300000	      5859 ns/op	     425 B/op	       9 allocs/op
BenchmarkLevelDbRead   1000000	      7522 ns/op	     725 B/op	      15 allocs/op

Done!
```

MySQL 5.6.23 via homebrew. Redis 2.8.19.

More will follow...

Feel free to submit your own or improvements to the test via PR.

Make sure that Redis runs on `127.0.0.1:6379` and if you test mysql: `test:test@127.0.0.1:3306`

The internal func `kv()` acquires 3 allocs so must subtract 3 allocs from the table above to get the real values.

Short discussion on [Twitter](https://twitter.com/schumacherfm/status/573060236166234112).

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
