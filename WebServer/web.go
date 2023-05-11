package WebServer

import (
	"GoFender/Database"
	"GoFender/Utils"
	"GoFender/YamlConfig"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func PaserPacket(c *gin.Context) {
	eps := Database.QueryEvalPacket()
	var myinfos []MyInfo
	for _, ep := range eps {
		if len(ep.ComSrcIp) > 1 || len(ep.ComDesIp) > 1 {
			phl := Utils.IpToLocation{}
			phl.GetIpInfo(ep.ComSrcIp, ep.ComDesIp)
			myinfo := MyInfo{
				Time:     ep.ComTime.String(),
				Type:     ep.AttackType,
				DestIp:   ep.ComDesIp,
				DestName: phl.DestPhyLocation,
				DestLocX: fmt.Sprintf("%f", phl.DestLocX),
				DestLocY: fmt.Sprintf("%f", phl.DestLocY),
				SrcIp:    ep.ComSrcIp,
				SrcName:  phl.SrcPhyLocation,
				SrcLocX:  fmt.Sprintf("%f", phl.SrcLocX),
				SrcLocY:  fmt.Sprintf("%f", phl.SrcLocY),
			}
			myinfos = append(myinfos, myinfo)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"data": myinfos,
	})
}

func WebStart() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(Cors())
	r.Static("/index", "./WebUI")
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/index/index.html")
	})

	ap := r.Group("/api")
	{
		ap.GET("/atkinfo", PaserPacket)
	}

	err := r.Run(YamlConfig.Myconfig.WebAddr)
	if err != nil {
		log.Fatal("WebServer Run Error: ", err)
	}
}

// Cors cross domain
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		if method != "" {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}
