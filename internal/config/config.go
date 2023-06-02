package config

import "github.com/spf13/viper"

// Функция инициализации конфига
func InitConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
