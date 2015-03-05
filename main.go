// The MIT License (MIT)

package main

import (
	"bytes"
	"database/sql"
	"encoding/binary"
	"flag"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"

	"github.com/boltdb/bolt"
	cznic "github.com/cznic/kv"
	"github.com/garyburd/redigo/redis"
	"github.com/peterbourgon/diskv"
	ledisConfig "github.com/siddontang/ledisdb/config"
	"github.com/siddontang/ledisdb/ledis"
	"github.com/steveyen/gkvlite"
	"github.com/syndtr/goleveldb/leveldb"
	// @todo https://github.com/zond/god
	// @todo https://github.com/HouzuoGuo/tiedot
	// @todo https://github.com/golang/groupcache ???
	"strings"
)

type (
	benchF func(b *testing.B)
	stores struct {
		bench []benchF
		run   *bool
	}
)

func main() {
	testAll := flag.Bool("all", false, "Run all tests")
	tests := []*stores{
		&stores{
			bench: []benchF{
				testMySQLWrite,
				testMySQLRead,
			},
			run: flag.Bool("mysql", false, ""),
		},
		&stores{
			bench: []benchF{
				testSQLiteWrite,
				testSQLiteRead,
			},
			run: flag.Bool("sqlite", false, ""),
		},
		&stores{
			bench: []benchF{
				testRedisWrite,
				testRedisRead,
			},
			run: flag.Bool("redis", false, ""),
		},
		&stores{
			bench: []benchF{
				testBoltWrite,
				testBoltRead,
			},
			run: flag.Bool("bolt", false, ""),
		},
		&stores{
			bench: []benchF{
				testGkvliteWrite,
				testGkvliteRead,
			},
			run: flag.Bool("gkvlite", false, ""),
		},
		&stores{
			bench: []benchF{
				testDiskvWrite,
				testDiskvRead,
			},
			run: flag.Bool("diskv", false, ""),
		},
		&stores{
			bench: []benchF{
				testCznicKvWrite,
				testCznicKvRead,
			},
			run: flag.Bool("cznickv", false, ""),
		},
		&stores{
			bench: []benchF{
				testLedisDbWrite,
				testLedisDbRead,
			},
			run: flag.Bool("ledisdb", false, ""),
		},
		&stores{
			bench: []benchF{
				testLevelDbWrite,
				testLevelDbRead,
			},
			run: flag.Bool("leveldb", false, ""),
		},
	}
	flag.Parse()

	ran := false
	for _, test := range tests {
		if *test.run == false && *testAll == false {
			continue
		}
		ran = true
		for _, f := range test.bench {
			res := testing.Benchmark(f)
			fmt.Printf("%s\t%s\t%s\n", funcName(f), res.String(), res.MemString())
		}
	}
	if ran == false {
		fmt.Println("No benchmark executed! Try -h switch for help.")
	}
	fmt.Println("\nDone!")
}

//func testPgWrite(b *testing.B) {
//	conn, err := pgx.Connect(pgx.ConnConfig{Host: "localhost", Port: 5432, Database: "test"})
//	isDoh(err)
//	conn.Exec("truncate table ids")
//	conn.Prepare("ids", "insert into ids values($1, $2)")
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		k, v := kv(i)
//		conn.Exec("ids", k, v)
//	}
//}

func sqlWrite(dn, dsn string, b *testing.B) *sql.DB {
	db, err := sql.Open(dn, dsn)
	isDoh(err)
	isDoh(db.Ping())

	_, err = db.Exec("DROP TABLE IF EXISTS `ids`")
	isDoh(err)
	_, err = db.Exec("CREATE TABLE `ids` (k1 int(10) NOT NULL,v1 int(10) NOT NULL, PRIMARY KEY (k1))")
	isDoh(err)
	stmt, err := db.Prepare("insert into ids values(?, ?)")
	isDoh(err)
	defer stmt.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = kv(i) // keep this because 3 additional allocs
		_, err := stmt.Exec(i, i)
		isDoh(err)
	}
	return db
}

func sqlRead(dn, dsn string, b *testing.B) {
	db := sqlWrite(dn, dsn, b)
	defer db.Close()
	stmt, err := db.Prepare("select v1 from `ids` where k1 = ?")
	defer stmt.Close()
	isDoh(err)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var v1 int
		err = stmt.QueryRow(i).Scan(&v1)
		isDoh(err)
		if v1 != i {
			b.Errorf("Expected: %d but got: %d", i, v1)
		}
	}
}

func testMySQLWrite(b *testing.B) {
	db := sqlWrite("mysql", "test:test@tcp(127.0.0.1:3306)/test", b)
	db.Close()
}

func testMySQLRead(b *testing.B) {
	sqlRead("mysql", "test:test@tcp(127.0.0.1:3306)/test", b)
}

func testSQLiteWrite(b *testing.B) {
	os.Remove("sqlite.db")
	defer os.Remove("sqlite.db")
	db := sqlWrite("sqlite3", "sqlite.db", b)
	db.Close()
}

func testSQLiteRead(b *testing.B) {
	os.Remove("sqlite.db")
	defer os.Remove("sqlite.db")
	sqlRead("sqlite3", "sqlite.db", b)
}

func redisWrite(b *testing.B) redis.Conn {
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	isDoh(err)
	conn.Do("flushdb")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k, v := kv(i)
		conn.Do("set", k, v)
	}
	conn.Do("save")
	return conn
}

func testRedisWrite(b *testing.B) {
	conn := redisWrite(b)
	conn.Close()
}

func testRedisRead(b *testing.B) {
	conn := redisWrite(b)
	defer conn.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k, v := kv(i)
		rv, err := conn.Do("get", k)
		if err != nil {
			b.Error(err)
		}
		if rvb, ok := rv.([]byte); ok {
			bc(b, v, rvb)
		} else {
			b.Errorf("Failed to convert %s", rv)
		}
	}

}

func boltWrite(b *testing.B) (*bolt.DB, []byte) {
	os.Remove("bolt.db")
	db, err := bolt.Open("bolt.db", 0600, nil)
	isDoh(err)
	bucket := []byte("MAIN")
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
	return db, bucket
}

func testBoltWrite(b *testing.B) {
	db, _ := boltWrite(b)
	defer db.Close()
}

func testBoltRead(b *testing.B) {
	db, bucket := boltWrite(b)
	defer db.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k, v := kv(i)
		db.View(func(tx *bolt.Tx) error {
			bb := tx.Bucket(bucket)
			bv := bb.Get(k)
			bc(b, v, bv)
			return nil
		})
	}
}

func gkvliteWrite(b *testing.B) (*os.File, *gkvlite.Store, *gkvlite.Collection) {
	os.Remove("test.gkvlite")
	f, err := os.Create("test.gkvlite")
	isDoh(err)

	s, err := gkvlite.NewStore(f)

	defer f.Sync()
	defer s.Flush()
	c := s.SetCollection("MAIN", nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k, v := kv(i)
		c.Set(k, v)
	}
	return f, s, c
}

func testGkvliteWrite(b *testing.B) {
	f, s, _ := gkvliteWrite(b)
	defer f.Close()
	defer s.Close()
}

func testGkvliteRead(b *testing.B) {
	f, s, c := gkvliteWrite(b)
	defer f.Close()
	defer s.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k, v := kv(i)
		v2, err := c.Get(k)
		if err != nil {
			b.Error(err)
		}
		bc(b, v, v2)
	}
}

func diskvWrite(b *testing.B) *diskv.Diskv {
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
	return d
}

func testDiskvWrite(b *testing.B) {
	_ = diskvWrite(b)
}

func testDiskvRead(b *testing.B) {
	d := diskvWrite(b)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k, v := kvs(i)
		v2, err := d.Read(k)
		if err != nil {
			b.Error(err)
		}
		bc(b, v, v2)
	}
}

func cznicKvWrite(b *testing.B) *cznic.DB {
	os.Remove("cznic.db")
	o := &cznic.Options{
	//		VerifyDbBeforeOpen:  true,
	//		VerifyDbAfterOpen:   true,
	//		VerifyDbBeforeClose: true,
	//		VerifyDbAfterClose:  true,
	}
	db, err := cznic.Create("cznic.db", o)
	isDoh(err)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k, v := kv(i)
		isDoh(db.Set(k, v))
	}
	return db
}

func testCznicKvWrite(b *testing.B) {
	db := cznicKvWrite(b)
	db.Close()
}

func testCznicKvRead(b *testing.B) {
	db := cznicKvWrite(b)
	defer db.Close()
	for i := 0; i < b.N; i++ {
		k, v := kv(i)
		v2, err := db.Get(nil, k)
		if err != nil {
			b.Error(err)
		}
		bc(b, v, v2)
	}
}

func ledisDbWrite(b *testing.B) (*ledis.Ledis, *ledis.DB) {
	dataDir := "ledis-test"
	os.RemoveAll(dataDir)
	cfg := ledisConfig.NewConfigDefault()
	cfg.DataDir = dataDir
	l, err := ledis.Open(cfg)
	isDoh(err)
	db, err := l.Select(0)
	isDoh(err)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k, v := kv(i)
		isDoh(db.Set(k, v))
	}
	return l, db
}

func testLedisDbWrite(b *testing.B) {
	l, _ := ledisDbWrite(b)
	l.Close()
}

func testLedisDbRead(b *testing.B) {
	l, db := ledisDbWrite(b)
	defer l.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k, v := kv(i)
		v2, err := db.Get(k)
		if err != nil {
			b.Error(err)
		}
		bc(b, v, v2)
	}
}

func leveldbWrite(b *testing.B) *leveldb.DB {
	dataDir := "leveldb-test"
	os.RemoveAll(dataDir)
	db, err := leveldb.OpenFile(dataDir, nil)
	isDoh(err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k, v := kv(i)
		isDoh(db.Put(k, v, nil))
	}
	return db
}

func testLevelDbWrite(b *testing.B) {
	db := leveldbWrite(b)
	db.Close()
}

func testLevelDbRead(b *testing.B) {
	db := leveldbWrite(b)
	defer db.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k, v := kv(i)
		v2, err := db.Get(k, nil)
		if err != nil {
			b.Error(err)
		}
		bc(b, v, v2)
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

func funcName(i interface{}) string {
	fn := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	return strings.Replace(fn, "main.test", "Benchmark", -1)
}

func bc(b *testing.B, expected, actual []byte) {
	if bytes.Compare(expected, actual) != 0 {
		b.Fatal("Expected %s got %s", expected, actual)
	}
}
