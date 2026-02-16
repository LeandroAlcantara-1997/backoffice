package env

import "github.com/Netflix/go-env"

type environment struct {
	APIPort     string `env:"PORT,default=8080"`
	APIName     string `env:"API_NAME,default=backoffice"`
	Environment string `env:"ENVIRONMENT"`
	TasksDLQBrokerConfig
	TasksInBrokerConfig
	TasksOutBrokerConfig
	RedisCache
}

type TasksInBrokerConfig struct {
	QueueName string `env:"TASKS_IN_NAME"`
	URL       string `env:"TASKS_IN_URL"`
}

type TasksDLQBrokerConfig struct {
	QueueName string `env:"TASKS_DLQ_NAME"`
	URL       string `env:"TASKS_DLQ_URL"`
}

type TasksOutBrokerConfig struct {
	QueueName string `env:"TASKS_OUT_NAME"`
	URL       string `env:"TASKS_OUT_URL"`
}

type RedisCache struct {
	Pass         string `env:"REDIS_PASSWORD"`
	Host         string `env:"REDIS_HOST"`
	Port         string `env:"REDIS_PORT"`
	ReadTimeout  int    `env:"REDIS_READ_TIMEOUT"`
	WriteTimeout int    `env:"REDIS_WRITE_TIMEOUT"`
}

var Env environment

func LoadEnv() (err error) {
	_, err = env.UnmarshalFromEnviron(&Env)
	return
}
