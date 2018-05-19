Common Mysql Functions
================

## Usage

```go
//init
con, err := sql.Open("mysql", "test:test@tcp(localhost:3306)/test")
if err != nil {
    panic(err)

}
db := &DB{con}

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
```

```go
//query to map slice
var maps []map[string]interface{}
err = db.QueryTo(&maps, "SELECT * FROM test WHERE id < ?;", 3)
if err != nil {
    panic(err)
}
```

```go
//qeury to struct slice
type Item struct {
	Id     int64                    `db:"id"`
	Name   string                   `db:"name"`
	Count  int                      `db:"count"`
	Money  float64                  `db:"money"`
	Data   string                   `db:"data"`
	List   []map[string]interface{} `db:"list,json"`
	Images []string                 `db:"images,json"`
	Ext    map[string]interface{}   `db:"ext,json"`
	Time   int64                    `db:"time,time"`    //millisecond timestamp
}

var items []*Item
err := db.QueryTo(&items, "SELECT * FROM test WHERE id < ?;", 3)
if err != nil {
    panic(err)
}
```