package config

import "os"

type MongoDb struct {
	Uri    string
	DbName string
}

type RedisDb struct {
	Addr string
}

func LoadMongo() *MongoDb {
	return &MongoDb{
		Uri:    os.Getenv("MONGO_URI"),
		DbName: os.Getenv("DB_NAME"),
	}
}

func LoadRedis() *RedisDb {
	return &RedisDb{
		Addr: os.Getenv("ADDR"),
	}
}
