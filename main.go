package main

import "CloudTusk/app"

func main() {
	distributor := app.LoadServers()

	go distributor.CheckLifeServers()

	app.New(distributor).Start()
}
