package main

type MysqlConfig struct {
	dbHost     string
	dbPort     int
	dbNAME     string
	dbUser     string
	dbPassword string
}

type RedisConfig struct {
	redisHost     string
	redisPort     int
	redisPassword string
	redisDatabase int
}
