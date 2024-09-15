package database

import "time"

type Config struct {
	MaxOpenConn     int `yaml:"max_open_conn"`
	MaxIdleConn     int `yaml:"max_idle_conn"`
	MaxIdleConnTime time.Duration `yaml:"max_idle_conn_time"`
	MaxConnLifeTime time.Duration `yaml:"max_conn_life_time"`
}
