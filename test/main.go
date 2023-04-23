package main

import (
	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/filter/cors"
	_ "github.com/lib/pq"
	_ "test/routers"
)

const (
	host   = "localhost"
	port   = "5432"
	user   = "postgres"
	pass   = "Devife4567"
	db     = "docs"
	schema = "public"
)

func main() {
	orm.RegisterDataBase(
		"default", "postgres",
		"postgres://"+user+":"+pass+"@"+host+":"+port+"/"+db+"?sslmode=disable&search_path="+schema,
	)
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.InsertFilter(
		"*", beego.BeforeRouter, cors.Allow(
			&cors.Options{
				AllowOrigins: []string{"*"},
				AllowMethods: []string{"PUT", "PATCH", "GET", "POST", "OPTIONS", "DELETE"},
				AllowHeaders: []string{
					"Origin", "x-requested-with",
					"content-type",
					"accept",
					"origin",
					"authorization",
					"x-csrftoken",
				},
				ExposeHeaders:    []string{"Content-Length"},
				AllowCredentials: true,
			},
		),
	)
	beego.Run()
}
