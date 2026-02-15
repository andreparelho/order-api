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
}

type Redis struct {
	Addr     string
	Password string
	User     string
	DBName   int
}

func Load() (*Configuration, error) {
	var appName, port, env string

	appName, errApp := getEnv("APP_NAME")
	if errApp != nil {
		return nil, errApp
	}

	port, errPort := getEnv("PORT")
	if errPort != nil {
		return nil, errPort
	}

	env, errEnv := getEnv("ENV")
	if errEnv != nil {
		return nil, errEnv
	}

	redisAddr, errAddr := getEnv("REDIS_ADDR")
	if errEnv != nil {
		return nil, errAddr
	}

	redisPassword, errPassw := getEnv("REDIS_PASSWORD")
	if errEnv != nil {
		return nil, errPassw
	}

	redisUser, errUser := getEnv("REDIS_USER")
	if errEnv != nil {
		return nil, errUser
	}

	return &Configuration{
		AppName: appName,
		Port:    port,
		Env:     env,
		Redis: Redis{
			Addr:     redisAddr,
			Password: redisPassword,
			User:     redisUser,
			DBName:   0,
		},
	}, nil
}

func getEnv(key string) (string, error) {
	if v := os.Getenv(key); v != "" {
		return v, nil
	}
	return "", fmt.Errorf("key: %s is not present in the environment variables file", key)
}
