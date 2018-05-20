package sqlcom

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	. "github.com/smartystreets/goconvey/convey"
)

func GetDb() *DB {
	db, err := sql.Open("mysql", "test:test@tcp(localhost:3306)/test")
	if err != nil {
		panic(err)

	}
	return &DB{db}
}

type Item struct {
	Id     int64                    `db:"id"`
	Name   string                   `db:"name"`
	Count  int                      `db:"count"`
	Money  float64                  `db:"money"`
	Data   string                   `db:"data"`
	List   []map[string]interface{} `db:"list,json"`
	Images []string                 `db:"images,json"`
	Ext    map[string]interface{}   `db:"ext,json"`
	Time   int64                    `db:"time,time"`
}

func TestQueryTo(t *testing.T) {
	db := GetDb()
	defer db.Close()

	Convey("create test table", t, func() {
		_, err := db.Exec("DROP TABLE IF EXISTS test;")
		So(err, ShouldBeNil)
		_, err = db.Exec(`CREATE TABLE test (
			id INT NOT NULL AUTO_INCREMENT,
			name VARCHAR(128),
			count INT,
			money DOUBLE,
			data VARBINARY(32),
			list JSON,
			images JSON,
			ext JSON,
			time DATETIME,
			PRIMARY KEY (id) 
		)`)
		So(err, ShouldBeNil)

		Convey("insert data", func() {
			now := time.Now()
			_, err := db.Exec("INSERT INTO test(name,count,money,data,list,images,ext,time) values"+
				"(?,?,?,?,?,?,?,FROM_UNIXTIME(?)),(?,?,?,?,?,?,?,FROM_UNIXTIME(?))",
				"name1", 1, 1.1, []byte("binary"), `[{"key": "val"}]`, `["url1","url2"]`, `{"num": 1}`, now.Unix(),
				"name2", 2, 1.2, []byte("binary2"), `[{"key": 2}]`, `["url1","url2"]`, `{"num": 1}`, now.Unix())
			So(err, ShouldBeNil)

			Convey("query to map slice", func() {
				var maps []map[string]interface{}
				err := db.QueryTo(&maps, "SELECT * FROM test WHERE id < ?;", 3)
				So(err, ShouldBeNil)
				So(len(maps), ShouldEqual, 2)
				So(len(maps[0]), ShouldEqual, 9)
				So(maps[0]["name"], ShouldEqual, "name1")
				So(maps[1]["count"], ShouldEqual, 2)
				So(maps[0]["money"], ShouldEqual, 1.1)
				So(maps[1]["data"], ShouldEqual, "binary2")
				So(maps[0]["list"], ShouldEqual, `[{"key": "val"}]`)
				So(maps[1]["images"], ShouldEqual, `["url1", "url2"]`)
				So(maps[0]["ext"], ShouldEqual, `{"num": 1}`)
				So(maps[1]["time"], ShouldEqual, now.Format("2006-01-02 15:04:05"))
			})

			Convey("query to struct slice", func() {
				var items []*Item
				err := db.QueryTo(&items, "SELECT * FROM test WHERE id < ?;", 3)
				So(err, ShouldBeNil)
				So(items[1].Name, ShouldEqual, "name2")
				So(items[0].Count, ShouldEqual, 1)
				So(items[1].Money, ShouldEqual, 1.2)
				So(items[0].Data, ShouldEqual, "binary")
				So(items[1].List, ShouldResemble, []map[string]interface{}{{"key": 2.0}})
				So(items[0].Images, ShouldResemble, []string{"url1", "url2"})
				So(items[1].Ext, ShouldResemble, map[string]interface{}{"num": 1.0})
				So(items[0].Time, ShouldEqual, now.Unix()*1e3)
			})
		})

		Convey("insert null value", func() {
			_, err := db.Exec("INSERT INTO test(id) values(10)")
			So(err, ShouldBeNil)

			Convey("query to map slice", func() {
				var maps []map[string]interface{}
				err := db.QueryTo(&maps, "SELECT * FROM test WHERE id = ?;", 10)
				So(err, ShouldBeNil)
				fmt.Println(maps)
				So(len(maps), ShouldEqual, 1)
				So(len(maps[0]), ShouldEqual, 1)
				So(maps[0]["id"], ShouldEqual, 10)
			})

			Convey("query to struct slice", func() {
				var items []*Item
				err := db.QueryTo(&items, "SELECT * FROM test WHERE id = ?;", 10)
				So(err, ShouldBeNil)
				So(items[0].Id, ShouldEqual, 10)
				So(items[0].Name, ShouldEqual, "")
				So(items[0].Count, ShouldEqual, 0)
				So(items[0].Money, ShouldEqual, 0)
				So(items[0].Data, ShouldEqual, "")
				So(items[0].List, ShouldResemble, []map[string]interface{}(nil))
				So(items[0].Images, ShouldResemble, []string(nil))
				So(items[0].Ext, ShouldResemble, map[string]interface{}(nil))
				So(items[0].Time, ShouldEqual, 0)
			})
		})

		Convey("insert null json value", func() {
			_, err := db.Exec("INSERT INTO test(id, ext) values(11, null)")
			So(err, ShouldBeNil)

			Convey("query to map slice", func() {
				var maps []map[string]interface{}
				err := db.QueryTo(&maps, "SELECT * FROM test WHERE id = ?;", 11)
				So(err, ShouldBeNil)
				fmt.Println(maps)
				So(len(maps), ShouldEqual, 1)
				So(len(maps[0]), ShouldEqual, 1)
				So(maps[0]["id"], ShouldEqual, 11)
			})

			Convey("query to struct slice", func() {
				var items []*Item
				err := db.QueryTo(&items, "SELECT * FROM test WHERE id = ?;", 11)
				So(err, ShouldBeNil)
				So(items[0].Id, ShouldEqual, 11)
				So(items[0].Ext, ShouldResemble, map[string]interface{}(nil))
			})
		})
	})
}

func mustInitTable(db *DB) {
	_, err := db.Exec("DROP TABLE IF EXISTS test;")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS test (
		id INT NOT NULL AUTO_INCREMENT,
		name VARCHAR(128),
		count INT,
		money DOUBLE,
		data VARBINARY(32),
		list JSON,
		images JSON,
		ext JSON,
		time DATETIME,
		PRIMARY KEY (id) 
	)`)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("INSERT INTO test(name,count,money,data,list,images,ext,time) values"+
		"(?,?,?,?,?,?,?,FROM_UNIXTIME(?)),(?,?,?,?,?,?,?,FROM_UNIXTIME(?))",
		"name1", 1, 1.1, []byte("binary"), `[{"key": "val"}]`, `["url1","url2"]`, `{"num": 1}`, time.Now().Unix(),
		"name2", 2, 1.2, []byte("binary2"), `[{"key": 2}]`, `["url1","url2"]`, `{"num": 1}`, time.Now().Unix(),
	)
	if err != nil {
		panic(err)
	}
}

func BenchmarkQueryMap(b *testing.B) {
	db := GetDb()
	defer db.Close()
	mustInitTable(db)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var maps []map[string]interface{}
		err := db.QueryTo(&maps, "SELECT * FROM test WHERE id < ?;", 3)
		if err != nil {
			b.Error(err)
			return
		}
	}
}

func BenchmarkQueryMapParallel(b *testing.B) {
	db := GetDb()
	defer db.Close()
	mustInitTable(db)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var maps []map[string]interface{}
			err := db.QueryTo(&maps, "SELECT * FROM test WHERE id < ?;", 3)
			if err != nil {
				b.Error(err)
				return
			}
		}
	})
}

func BenchmarkQueryStruct(b *testing.B) {
	db := GetDb()
	defer db.Close()
	mustInitTable(db)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var items []*Item
		err := db.QueryTo(&items, "SELECT * FROM test WHERE id < ?;", 3)
		if err != nil {
			b.Error(err)
			return
		}
	}
}

func BenchmarkQueryStructParallel(b *testing.B) {
	db := GetDb()
	defer db.Close()
	mustInitTable(db)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var items []*Item
			err := db.QueryTo(&items, "SELECT * FROM test WHERE id < ?;", 3)
			if err != nil {
				b.Error(err)
				return
			}
		}
	})
}
