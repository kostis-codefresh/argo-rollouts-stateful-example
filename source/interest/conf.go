package main

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func (interestApp *InterestApplication) readCurrentRole() {
	viper.SetConfigName("labels")
	viper.SetConfigType("properties") //Java properties style

	//Development mode
	viper.AddConfigPath(".")

	//This is inject from the Kubernetes downward API
	viper.AddConfigPath("/etc/podinfo/")

	// viper.SetDefault("role", "demo")

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	// str, err = strconv.Unquote("'\u2639\u2639'")

	//Reload configuration when file changes
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		fmt.Printf("Role is %s\n", viper.GetString("role"))
	})
	viper.WatchConfig()

	fmt.Printf("Role is set %t\n", viper.IsSet("role"))
	fmt.Printf("Role is %s\n", viper.GetString("role"))
}
