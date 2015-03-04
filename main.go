package main

/*
@karlseguin:
bolt      5000	    277963 ns/op
redis    30000	     48081 ns/op
pg       10000	    149691 ns/op

Yes, the Bolt transactions could be batched. But so too could the PG transactions,
and the Redis work could be pipelined. And that isn't always a workable solution.

@schumacherfm:
- requirement that data must be persisted some how
*/

import (
	"encoding/binary"
	"os"
	"strconv"
	"testing"

	"flag"

	"fmt"

	"github.com/boltdb/bolt"
	cznic "github.com/cznic/kv"
	"github.com/garyburd/redigo/redis"
	"github.com/jackc/pgx"
	"github.com/peterbourgon/diskv"
	ledisConfig "github.com/siddontang/ledisdb/config"
	"github.com/siddontang/ledisdb/ledis"
	"github.com/steveyen/gkvlite"
	// @todo https://github.com/zond/god
	// @todo https://github.com/syndtr/goleveldb
	// @todo https://github.com/HouzuoGuo/tiedot
	// @todo mysql
	// @todo sqlite3
	// @todo https://github.com/golang/groupcache ???
)

type (
	stores struct {
		name  string
		bench func(b *testing.B)
		run   *bool
	}
)

func main() {
	testAll := flag.Bool("all", false, "Run all tests")
	tests := []*stores{
		&stores{
			name:  "BoltDB",
			bench: testBolt,
			run:   flag.Bool("bolt", false, ""),
		},
		&stores{
			name:  "Redis",
			bench: testRedis,
			run:   flag.Bool("redis", false, ""),
		},
		&stores{
			name:  "PostgreSQL",
			bench: testPG,
			run:   flag.Bool("postgres", false, ""),
		},
		&stores{
			name:  "gkvlite",
			bench: testGkvlite,
			run:   flag.Bool("gkvlite", false, ""),
		},
		&stores{
			name:  "diskv",
			bench: testDiskv,
			run:   flag.Bool("diskv", false, ""),
		},
		&stores{
			name:  "cznickv",
			bench: testCznicKv,
			run:   flag.Bool("cznickv", false, ""),
		},
		&stores{
			name:  "ledisdb",
			bench: testLedisDb,
			run:   flag.Bool("ledisdb", false, ""),
		},
	}
	flag.Parse()

	ran := false
	for _, test := range tests {
		if *test.run == false && *testAll == false {
			continue
		}
		ran = true
		res := testing.Benchmark(test.bench)
		fmt.Printf("%s\t%s\t%s\n", test.name, res.String(), res.MemString())
	}
	if ran == false {
		fmt.Println("No benchmark executed! Try -h switch for help.")
	}
	fmt.Println("\nDone!")
}

func testPG(b *testing.B) {
	conn, err := pgx.Connect(pgx.ConnConfig{Host: "localhost", Port: 5432, Database: "test"})
	isDoh(err)
	conn.Exec("truncate table ids")
	conn.Prepare("ids", "insert into ids values($1, $2)")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k, v := kv(i)
		conn.Exec("ids", k, v)
	}
}

func testRedis(b *testing.B) {
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	isDoh(err)
	defer conn.Close()
	conn.Do("flushdb")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k, v := kv(i)
		conn.Do("set", k, v)
	}
	conn.Do("save")
}

func testBolt(b *testing.B) {
	os.Remove("bolt.db")
	db, err := bolt.Open("bolt.db", 0600, nil)
	isDoh(err)
	bucket := []byte("MAIN")
	defer db.Close()
	db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists(bucket)
		return nil
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k, v := kv(i)
		db.Update(func(tx *bolt.Tx) error {
			return tx.Bucket(bucket).Put(k, v)
		})
	}
}

func testGkvlite(b *testing.B) {
	os.Remove("test.gkvlite")
	f, err := os.Create("test.gkvlite")
	isDoh(err)
	defer f.Close()
	s, err := gkvlite.NewStore(f)
	defer s.Close()
	defer f.Sync()
	defer s.Flush()
	c := s.SetCollection("MAIN", nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k, v := kv(i)
		c.Set(k, v)
	}
}

func testDiskv(b *testing.B) {
	dir := "pbdiskv"
	os.RemoveAll(dir)
	d := diskv.New(diskv.Options{
		BasePath:     dir,
		Transform:    func(s string) []string { return []string{} },
		CacheSizeMax: 1024 * 1024, // 1MB
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k, v := kvs(i)
		isDoh(d.Write(k, v))

	}
}

func testCznicKv(b *testing.B) {
	os.Remove("cznic.db")
	o := &cznic.Options{
	//		VerifyDbBeforeOpen:  true,
	//		VerifyDbAfterOpen:   true,
	//		VerifyDbBeforeClose: true,
	//		VerifyDbAfterClose:  true,
	}
	db, err := cznic.Create("cznic.db", o)
	isDoh(err)
	defer db.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k, v := kv(i)
		isDoh(db.Set(k, v))
	}
}

func testLedisDb(b *testing.B) {
	dataDir := "ledis-test"
	os.RemoveAll(dataDir)
	cfg := ledisConfig.NewConfigDefault()
	cfg.DataDir = dataDir
	l, err := ledis.Open(cfg)
	isDoh(err)
	defer l.Close()
	db, err := l.Select(0)
	isDoh(err)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k, v := kv(i)
		isDoh(db.Set(k, v))
	}

}

func kv(i int) ([]byte, []byte) {
	k := []byte(strconv.Itoa(i))
	v := make([]byte, 8)
	binary.LittleEndian.PutUint64(v, uint64(i))
	return k, v
}

func kvs(i int) (string, []byte) {
	v := make([]byte, 8)
	binary.LittleEndian.PutUint64(v, uint64(i))
	return strconv.Itoa(i), v
}

func isDoh(err error) {
	if err != nil {
		panic(err)
	}
}
