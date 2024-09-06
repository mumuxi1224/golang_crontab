package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
)

type Crontab struct {
	mysqlConf MysqlConfig
	redisConf RedisConfig
	listen    string
	db        *gorm.DB
}

/*
*

	工厂函数，在创建 Crontab 实例时自动执行初始化操作
*/
func getCrontab() *Crontab {
	c := &Crontab{}
	c.initConfig()
	return c
}

// 初始化配置
func (c *Crontab) initConfig() {
	// 获取当前文件所在的目录
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Println("Failed to get caller information")
		return
	}

	// 获取当前文件的完整路径
	currentFile := filepath.Clean(filename)

	// 获取当前文件所在的目录
	currentDir := filepath.Dir(currentFile)

	err := godotenv.Load(currentDir + "/.env") // 也可以指定文件路径如 godotenv.Load(".env")
	if err != nil {
		log.Fatal("No .env file found" + err.Error())
	}

	// mysql配置
	//c.mysqlConf :=MysqlConfig{}
	c.mysqlConf.dbHost = os.Getenv("DB_HOST")
	c.mysqlConf.dbPort, _ = strconv.Atoi(os.Getenv("DB_PORT"))
	c.mysqlConf.dbNAME = os.Getenv("DB_NAME")
	c.mysqlConf.dbUser = os.Getenv("DB_USER")
	c.mysqlConf.dbPassword = os.Getenv("DB_PASSWORD")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.mysqlConf.dbUser, c.mysqlConf.dbPassword, c.mysqlConf.dbHost, c.mysqlConf.dbPort, c.mysqlConf.dbNAME)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("failed to connect database", err)
	}
	c.db = db

	dd, err := c.db.DB()
	c.db.Table()
	fmt.Println(dd.Ping())
	fmt.Println(dd.Query("select  1 * 1"))

	fmt.Println(1111223213)
	// redis配置
	//c.redisConf :=RedisConfig{}
	c.redisConf.redisHost = os.Getenv("REDIS_HOST")
	c.redisConf.redisPort, _ = strconv.Atoi(os.Getenv("REDIS_PORT"))
	c.redisConf.redisPassword = os.Getenv("REDIS_PASSWORD")
	c.redisConf.redisDatabase, _ = strconv.Atoi(os.Getenv("REDIS_DATABASE"))

	// 监听地址
	c.listen = os.Getenv("CRONTAB_LISTEN")

	fmt.Println(c)
}

// 执行tcp服务
func (c *Crontab) Run() {
	fmt.Println("Crontab Run")
	listen, err := net.Listen("tcp", c.listen)
	if err != nil {
		fmt.Println("listen error,err:", err)
		return
	}
	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("accept error,err:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("accept success")
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read error,err:", err)
			return
		}
		fmt.Println("read success,data:", string(buf[:n]))
		conn.Write([]byte("hello world"))
	}
}
