package main

import (
	"fmt"
	"learn_testing/config"
	"learn_testing/routes"
)

func main() {
	config.Init()

	// start the server, and log if it fails
	e := routes.New()
	e.Logger.Fatal(e.Start(":8000"))
	fmt.Println(config.ViperEnvVariable("SECRET_KEY"))
}
