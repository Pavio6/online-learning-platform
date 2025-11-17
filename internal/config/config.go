package config

import (
	"github.com/spf13/viper"
)

// Config 应用配置
type Config struct {
	App      AppConfig      `mapstructure:"app"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Database DatabaseConfig `mapstructure:"database"`
	Branches []BranchConfig `mapstructure:"branches"`
	OSS      OSSConfig      `mapstructure:"oss"`
	Sync     SyncConfig     `mapstructure:"sync"`
}

// AppConfig 应用配置
type AppConfig struct {
	Name     string `mapstructure:"name"`
	Port     int    `mapstructure:"port"`
	Env      string `mapstructure:"env"`
	LogLevel string `mapstructure:"log_level"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret     string `mapstructure:"secret"`
	Expiration string `mapstructure:"expiration"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Central DBSettings `mapstructure:"central"`
}

// DBSettings 数据库连接设置
type DBSettings struct {
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	User            string `mapstructure:"user"`
	Password        string `mapstructure:"password"`
	DBName          string `mapstructure:"dbname"`
	SSLMode         string `mapstructure:"sslmode"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	ConnMaxLifetime string `mapstructure:"conn_max_lifetime"`
}

// BranchConfig 分支节点配置
type BranchConfig struct {
	BranchID uint       `mapstructure:"branch_id"`
	Name     string     `mapstructure:"name"`
	DB       DBSettings `mapstructure:",squash"`
}

// OSSConfig OSS配置
type OSSConfig struct {
	Endpoint        string `mapstructure:"endpoint"`
	AccessKeyID     string `mapstructure:"access_key_id"`
	AccessKeySecret string `mapstructure:"access_key_secret"`
	BucketName      string `mapstructure:"bucket_name"`
	Region          string `mapstructure:"region"`
}

// SyncConfig 同步配置
type SyncConfig struct {
	Replication  ReplicationConfig  `mapstructure:"replication"`
	Consolidation ConsolidationConfig `mapstructure:"consolidation"`
}

// ReplicationConfig РОК同步配置
type ReplicationConfig struct {
	Enabled  bool     `mapstructure:"enabled"`
	Schedule string   `mapstructure:"schedule"`
	Tables   []string `mapstructure:"tables"`
}

// ConsolidationConfig РОК+КД整合配置
type ConsolidationConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	Schedule string `mapstructure:"schedule"`
}

var globalConfig *Config

// LoadConfig 加载配置
func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigType("yaml")
	viper.SetConfigFile(configPath)

	// 设置环境变量
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	globalConfig = &config
	return &config, nil
}

// GetConfig 获取全局配置
func GetConfig() *Config {
	return globalConfig
}

// GetCentralDBConfig 获取中央服务器数据库配置
func GetCentralDBConfig() DBSettings {
	return globalConfig.Database.Central
}

// GetBranchDBConfig 获取分支节点数据库配置
func GetBranchDBConfig(branchID uint) *BranchConfig {
	for _, branch := range globalConfig.Branches {
		if branch.BranchID == branchID {
			return &branch
		}
	}
	return nil
}

// GetAllBranches 获取所有分支配置
func GetAllBranches() []BranchConfig {
	return globalConfig.Branches
}

// GetOSSConfig 获取OSS配置
func GetOSSConfig() OSSConfig {
	return globalConfig.OSS
}

// GetJWTConfig 获取JWT配置
func GetJWTConfig() JWTConfig {
	return globalConfig.JWT
}

