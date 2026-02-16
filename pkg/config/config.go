package config

import (
	"fmt"
	"os"
)

type Configuration struct {
	AppName string
	Port    string
	Env     string
	RDS     RDS
	Redis   Redis
}

type RDS struct {
	User     string
	Password string
	Addr     string
	DBName   string
}

type Redis struct {
	Addr     string
	Password string
	User     string
}

func Load() (*Configuration, error) {

	appName, err := getEnv("APP_NAME")
	if err != nil {
		return nil, err
	}

	port, err := getEnv("PORT")
	if err != nil {
		return nil, err
	}

	env, err := getEnv("ENV")
	if err != nil {
		return nil, err
	}

	redisAddr, err := getEnv("REDIS_ADDR")
	if err != nil {
		return nil, err
	}

	redisPassword, err := getEnv("REDIS_PASSWORD")
	if err != nil {
		return nil, err
	}

	redisUser, err := getEnv("REDIS_USER")
	if err != nil {
		return nil, err
	}

	rdsAddr, err := getEnv("RDS_ADDR")
	if err != nil {
		return nil, err
	}

	rdsPassword, err := getEnv("RDS_PASSWORD")
	if err != nil {
		return nil, err
	}

	rdsUsername, err := getEnv("RDS_USER")
	if err != nil {
		return nil, err
	}

	rdsDbName, err := getEnv("RDS_DBNAME")
	if err != nil {
		return nil, err
	}

	return &Configuration{
		AppName: appName,
		Port:    port,
		Env:     env,
		Redis: Redis{
			Addr:     redisAddr,
			Password: redisPassword,
			User:     redisUser,
		},
		RDS: RDS{
			User:     rdsUsername,
			Password: rdsPassword,
			Addr:     rdsAddr,
			DBName:   rdsDbName,
		},
	}, nil
}

func getEnv(key string) (string, error) {
	if v := os.Getenv(key); v != "" {
		return v, nil
	}
	return "", fmt.Errorf("key: %s is not present in the environment variables file", key)
}
