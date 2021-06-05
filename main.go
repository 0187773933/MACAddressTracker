package main

import (
	"fmt"
	utils "github.com/0187773933/MACAddressTracker/v1/utils"
	cron "github.com/robfig/cron/v3"
	fiber "github.com/gofiber/fiber/v2"
)

var config utils.ConfigFile
func run() {
	local_network := utils.ScanLocalNetwork( config )
	// utils.PrintLocalNetwork( config , local_network )
	utils.TrackChanges( config , local_network )
}

func main() {
	config = utils.GetConfig()
	fmt.Println( config.Devices )
	cron_runner := cron.New()
	cron_runner.AddFunc( config.CronString , run )
	cron_runner.Start()
	app := fiber.New()
	app.Get( "/update" , func( context *fiber.Ctx ) ( error ) {
		go run()
		context.Set( fiber.HeaderContentType , fiber.MIMETextHTML )
		return context.SendString( "<html><h1>ok</h1></html>" )
	})
	port := fmt.Sprintf( ":%s" , config.ServerPort )
	fmt.Printf( "Listening on %s\n" , port )
	result := app.Listen( port )
	fmt.Println( result )
}