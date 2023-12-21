package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func (interestApp *InterestApplication) readCurrentRole() {
	viper.SetConfigName("labels")
	viper.SetConfigType("properties") //Java properties style

	//Development mode
	viper.AddConfigPath(".")

	//This is injected from the Kubernetes downward API that maps
	// all labels as a file in the pod
	// See https://kubernetes.io/docs/concepts/workloads/pods/downward-api/
	viper.AddConfigPath("/etc/podinfo/")

	viper.SetDefault("role", "unknown")

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	//Reload configuration when file changes
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		fmt.Printf("Role is %s\n", viper.GetString("role"))
	})
	viper.WatchConfig()

	fmt.Printf("Role is set %t\n", viper.IsSet("role"))
	fmt.Printf("Role is %s\n", viper.GetString("role"))

	interestApp.CurrentRole = viper.GetString("role")
	if strings.HasPrefix(interestApp.CurrentRole, "\"") {
		interestApp.CurrentRole, _ = strconv.Unquote(interestApp.CurrentRole)
	}

}
