package utils

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

func GetConfig(key string) interface{} {
	//configs.yaml
	corrPath, err := os.Getwd() //获取项目的执行路径
	if err != nil {
		fmt.Println(err)
	}
	config := viper.New()
	config.AddConfigPath(corrPath)  //设置读取的文件路径
	config.SetConfigName("configs") //设置读取的文件名
	config.SetConfigType("yaml")    //设置文件的类型
	err = config.ReadInConfig()     //尝试进行配置读取
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	return config.Get(key)
}
