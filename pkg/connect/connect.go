package connect

import (
	"fmt"
	"sync"
	"time"

	"github.com/example/golang-api/pkg/env"
	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

var connections = new(sync.Map)
var onces = new(sync.Map)

func NewDB(dbname string) (*xorm.Engine, error) {
	once, _ := onces.LoadOrStore(dbname, new(sync.Once))
	once.(*sync.Once).Do(func() {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=UTC&charset=utf8mb4",
			env.Getenv(fmt.Sprintf("%s_DB_USER", dbname)),
			env.Getenv(fmt.Sprintf("%s_DB_PASSWORD", dbname)),
			env.Getenv(fmt.Sprintf("%s_DB_HOST", dbname)),
			env.Getenv(fmt.Sprintf("%s_DB_PORT", dbname), "3306"),
			env.Getenv(fmt.Sprintf("%s_DB_NAME", dbname)),
		)

		eng, err := xorm.NewEngine("mysql", dsn)
		if err != nil {
			panic(fmt.Sprintf("connect to %s failed: %v", dbname, err))
		}

		dbConn := env.GetenvInt(fmt.Sprintf("%s_DB_CONN", dbname), 10)
		eng.SetMaxIdleConns(dbConn)
		eng.SetMaxOpenConns(50)
		eng.SetConnMaxLifetime(30 * time.Minute)

		connections.Store(dbname, eng)
	})

	return Engine(dbname)
}

func Engine(dbname string) (*xorm.Engine, error) {
	eg, ok := connections.Load(dbname)
	if !ok {
		return nil, fmt.Errorf("db engine not found: %s", dbname)
	}
	return eg.(*xorm.Engine), nil
}

func Session(dbname string) (*xorm.Session, error) {
	eg, err := Engine(dbname)
	if err != nil {
		return nil, err
	}
	return eg.NewSession(), nil
}
