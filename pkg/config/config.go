package config

import (
	"fmt"
	"os"
)

type Configuration struct {
	AppName  string
	Port     string
	Env      string
	RDS      RDS
	Redis    Redis
	SQS      SQS
	DynamoDB DynamoDB
	AWS      AWS
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

type SQS struct {
	OrdersQueue   string
	PaymentsQueue string
}

type DynamoDB struct {
	TableName string
}

type AWS struct {
	Key      string
	Secret   string
	Session  string
	Endpoint string
	Region   string
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

	sqsOrdersQueue, err := getEnv("AWS_SQS_ORDERS_QUEUE_NAME")
	if err != nil {
		return nil, err
	}

	sqsPaymentsQueue, err := getEnv("AWS_SQS_PAYMENTS_QUEUE_NAME")
	if err != nil {
		return nil, err
	}

	dynamoTableName, err := getEnv("DYNAMO_TABLE_NAME")
	if err != nil {
		return nil, err
	}

	awsRegion, err := getEnv("AWS_REGION")
	if err != nil {
		return nil, err
	}

	awsEndpoint, err := getEnv("AWS_ENDPOINT")
	if err != nil {
		return nil, err
	}

	awsAccessKey, err := getEnv("AWS_ACCESS_KEY_ID")
	if err != nil {
		return nil, err
	}

	awsSecretKey, err := getEnv("AWS_SECRET_ACCESS_KEY")
	if err != nil {
		return nil, err
	}

	awsSession, err := getEnv("AWS_SESSSION")
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
		SQS: SQS{
			OrdersQueue:   sqsOrdersQueue,
			PaymentsQueue: sqsPaymentsQueue,
		},
		DynamoDB: DynamoDB{
			TableName: dynamoTableName,
		},
		AWS: AWS{
			Key:      awsAccessKey,
			Secret:   awsSecretKey,
			Session:  awsSession,
			Endpoint: awsEndpoint,
			Region:   awsRegion,
		},
	}, nil
}

func getEnv(key string) (string, error) {
	if v := os.Getenv(key); v != "" {
		return v, nil
	}
	return "", fmt.Errorf("key: %s is not present in the environment variables file", key)
}
