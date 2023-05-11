package Database

import (
	"GoFender/Utils"
	"gorm.io/gorm"
	"log"
	"time"
)

type Commpacket struct {
	gorm.Model
	ComTime       time.Time
	ComDesIp      string
	ComSrcIp      string
	ComDesPort    string
	ComSrcPort    string
	ComProtocol   string
	Type          string
	ComPacketData []byte `gorm:"type:blob"`
}

type Evilpacket struct {
	gorm.Model
	ComTime       time.Time
	ComDesIp      string
	ComSrcIp      string
	ComDesPort    string
	ComSrcPort    string
	ComProtocol   string
	Type          string
	AttackType    string
	ComPacketData []byte `gorm:"type:blob"`
}

var DataPool *MysqlPool

func Insert(SqlStruct interface{}) {
	conn, err := DataPool.GetConn()
	if err != nil {
		panic(err)
	}
	defer func(pool *MysqlPool, conn *MysqlConn) {
		err := pool.PutConn(conn)
		if err != nil {
			panic(err)
		}
	}(DataPool, conn)

	switch SqlStruct.(type) {
	case Utils.NomPacket:
		cp := Commpacket{
			ComTime:       SqlStruct.(Utils.NomPacket).CommInfo.ComTime,
			ComDesIp:      SqlStruct.(Utils.NomPacket).CommInfo.ComDesIp,
			ComSrcIp:      SqlStruct.(Utils.NomPacket).CommInfo.ComSrcIp,
			ComDesPort:    SqlStruct.(Utils.NomPacket).CommInfo.ComDesPort,
			ComSrcPort:    SqlStruct.(Utils.NomPacket).CommInfo.ComSrcPort,
			ComProtocol:   SqlStruct.(Utils.NomPacket).CommInfo.ComProtocol,
			Type:          SqlStruct.(Utils.NomPacket).Type,
			ComPacketData: SqlStruct.(Utils.NomPacket).CommInfo.ComPacketData,
		}
		if !conn.db.Migrator().HasTable(&cp) {
			err := conn.db.Migrator().CreateTable(&cp)
			if err != nil {
				log.Fatal(err)
			}
		}
		if err := conn.db.Create(&cp).Error; err != nil {
			panic("Insert Data Error" + err.Error())
		}
	case Utils.EvilPacket:
		ep := Evilpacket{
			ComTime:       SqlStruct.(Utils.EvilPacket).CommInfo.ComTime,
			ComDesIp:      SqlStruct.(Utils.EvilPacket).CommInfo.ComDesIp,
			ComSrcIp:      SqlStruct.(Utils.EvilPacket).CommInfo.ComSrcIp,
			ComDesPort:    SqlStruct.(Utils.EvilPacket).CommInfo.ComDesPort,
			ComSrcPort:    SqlStruct.(Utils.EvilPacket).CommInfo.ComSrcPort,
			ComProtocol:   SqlStruct.(Utils.EvilPacket).CommInfo.ComProtocol,
			Type:          SqlStruct.(Utils.EvilPacket).Type,
			AttackType:    SqlStruct.(Utils.EvilPacket).AttackType,
			ComPacketData: SqlStruct.(Utils.EvilPacket).CommInfo.ComPacketData,
		}
		if !conn.db.Migrator().HasTable(&ep) {
			err := conn.db.Migrator().CreateTable(&ep)
			if err != nil {
				log.Fatal(err)
			}
		}
		if err := conn.db.Create(&ep).Error; err != nil {
			log.Fatal("Insert Data Error" + err.Error())
		}
	}

}

func QueryEvalPacket() []Evilpacket {
	conn, err := DataPool.GetConn()
	if err != nil {
		panic(err)
	}
	defer func(pool *MysqlPool, conn *MysqlConn) {
		err := pool.PutConn(conn)
		if err != nil {
			panic(err)
		}
	}(DataPool, conn)
	var resultArr []Evilpacket
	conn.db.Raw("SELECT * FROM evilpackets WHERE com_time >= DATE_SUB(NOW(),INTERVAL 15 MINUTE) LIMIT 15").Scan(&resultArr)
	return resultArr
}
