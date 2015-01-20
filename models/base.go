package models

import (
	"crypto/md5"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"net/url"
	"strings"
)

func Init() {
	/*
		dbhost := beego.AppConfig.String("dbhost")
		dbport := beego.AppConfig.String("dbport")
		dbuser := beego.AppConfig.String("dbuser")
		dbpassword := beego.AppConfig.String("dbpassword")
		dbname := beego.AppConfig.String("dbname")
		if dbport == "" {
			dbport = "3306"
		}
		dsn := dbuser + ":" + dbpassword + "@tcp(" + dbhost + ":" + dbport + ")/" + dbname + "?charset=utf8"
		orm.RegisterDataBase("default", "mysql", dsn)
	*/
	orm.Debug = true
	orm.RegisterDataBase("default", "sqlite3", "data.db", 30)
	orm.RegisterModel(new(User), new(Post), new(Tag), new(Option), new(TagPost))
	/* 自动建表*/
	name := "default"                          //数据库别名
	force := false                             //不强制建数据库
	verbose := true                            //打印建表过程
	err := orm.RunSyncdb(name, force, verbose) //建表
	if err != nil {
		beego.Error(err)
	}

}

func Md5(buf []byte) string {
	hash := md5.New()
	hash.Write(buf)
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func Rawurlencode(str string) string {
	return strings.Replace(url.QueryEscape(str), "+", "%20", -1)
}

func GetOptions() map[string]string {
	if !Cache.IsExist("options") {
		var result []*Option
		o := orm.NewOrm()
		o.QueryTable(&Option{}).All(&result)
		options := make(map[string]string)
		for _, v := range result {
			options[v.Name] = v.Value
		}
		Cache.Put("options", options, 0)
	}
	v := Cache.Get("options")
	return v.(map[string]string)
}

//返回带前缀的表名
func TableName(str string) string {
	return fmt.Sprintf("%s%s", beego.AppConfig.String("dbprefix"), str)
}
