package Database

import (
	"GoFender/YamlConfig"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
	"time"
)

type MysqlPool struct {
	sync.Mutex
	db             *gorm.DB
	maxIdleConns   int
	maxOpenConns   int
	maxConnTimeout time.Duration
	connTimeout    time.Duration
	idleTimeout    time.Duration
	waitTimeout    time.Duration
	closed         bool
	idleConns      []*MysqlConn
}

type MysqlConn struct {
	db        *gorm.DB
	createdAt time.Time
}

// NewMysqlPool creates a new mysql connection pool
func NewMysqlPool() (*MysqlPool, error) {
	MysqlDns :=
		YamlConfig.Myconfig.MysqlUser + ":" +
			YamlConfig.Myconfig.MysqlPassword + "@tcp(" +
			YamlConfig.Myconfig.MysqlServer + ")/" +
			"?charset=utf8mb4&parseTime=True&loc=Local"

	db, _ := gorm.Open(mysql.New(mysql.Config{
		DSN:                       MysqlDns, //DSN data source name
		DefaultStringSize:         1500,
		SkipInitializeWithVersion: true, // autoconfigure based on currently MySQL version
	}), &gorm.Config{})

	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(1000)           //set max open connections
	sqlDB.SetMaxIdleConns(20)             //set max idle connections
	sqlDB.SetConnMaxLifetime(time.Minute) //set max lifetime of connection
	_, err := sqlDB.Exec("CREATE DATABASE IF NOT EXISTS " + YamlConfig.Myconfig.MysqlDatabase + " DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_general_ci;")
	if err != nil {
		panic("database init failed" + err.Error())
	}

	err = sqlDB.Close()
	if err != nil {
		panic("close db server failed" + err.Error())
	}
	pool := &MysqlPool{
		db:             db,
		maxIdleConns:   YamlConfig.Myconfig.MysqlMaxIdleConns,
		maxOpenConns:   YamlConfig.Myconfig.MysqlMaxOpenConns,
		maxConnTimeout: 10 * time.Second,
		connTimeout:    5 * time.Second,
		idleTimeout:    2*time.Minute + 30*time.Second,
		waitTimeout:    30 * time.Second,
		closed:         false,
		idleConns:      []*MysqlConn{},
	}

	// create initial connections
	for i := 0; i < YamlConfig.Myconfig.MysqlMaxIdleConns; i++ {
		conn, err := pool.newConn()
		if err != nil {
			return nil, err
		}
		pool.idleConns = append(pool.idleConns, conn)
	}
	return pool, nil
}

// Open creates a new mysql connection
func (pool *MysqlPool) Open() (*MysqlConn, error) {
	conn, err := pool.newConn()
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// Close closes the mysql connection pool
func (pool *MysqlPool) Close() error {
	if pool.closed {
		return nil
	}
	pool.closed = true
	for _, conn := range pool.idleConns {
		SqlDb, _ := conn.db.DB()
		SqlDb.Close()
	}
	return nil
}

// GetConn returns a mysql connection from the pool
func (pool *MysqlPool) GetConn() (*MysqlConn, error) {
	for {
		// check if the pool is closed
		if pool.closed {
			return nil, fmt.Errorf("mysql pool is closed")
		}

		// check if there are any idle connections
		n := len(pool.idleConns)
		if n > 0 {
			conn := pool.idleConns[n-1]
			pool.idleConns = pool.idleConns[:n-1]
			if time.Since(conn.createdAt) > pool.maxConnTimeout {
				SqlDb, _ := conn.db.DB()
				SqlDb.Close()
				continue
			}
			return conn, nil
		}

		// check if we can create a new connection
		if pool.maxOpenConns <= 0 || pool.maxOpenConns > len(pool.idleConns)+1 {
			conn, err := pool.newConn()
			if err != nil {
				return nil, err
			}
			return conn, nil
		}

		// wait for a connection to become available
		wait := pool.waitTimeout
		if wait <= 0 {
			continue
		}

		select {
		case <-time.After(wait):
			continue
		case <-time.After(time.Until(time.Now().Add(wait))):
			return nil, fmt.Errorf("mysql pool wait timeout")
		}
	}
}

// PutConn returns a mysql connection to the pool
func (pool *MysqlPool) PutConn(conn *MysqlConn) error {
	// check if the pool is closed
	if pool.closed {
		SqlDb, _ := conn.db.DB()
		return SqlDb.Close()
	}

	// check if the connection is still valid
	if time.Since(conn.createdAt) > pool.maxConnTimeout {
		SqlDb, _ := conn.db.DB()
		return SqlDb.Close()
	}

	// add the connection back to the pool
	pool.idleConns = append(pool.idleConns, conn)
	return nil
}

// newConn creates a new mysql connection
func (pool *MysqlPool) newConn() (*MysqlConn, error) {
	MysqlDns :=
		YamlConfig.Myconfig.MysqlUser + ":" +
			YamlConfig.Myconfig.MysqlPassword + "@tcp(" +
			YamlConfig.Myconfig.MysqlServer + ")/" +
			YamlConfig.Myconfig.MysqlDatabase +
			"?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       MysqlDns, //DSN data source name
		DefaultStringSize:         1500,
		SkipInitializeWithVersion: true, // autoconfigure based on currently MySQL version
	}), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	conn := &MysqlConn{db: db, createdAt: time.Now()}
	return conn, nil
}

func InitDatabase() *MysqlPool {
	pool, err := NewMysqlPool()
	if err != nil {
		panic(err)
	}
	return pool
}
