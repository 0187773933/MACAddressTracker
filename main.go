package main

import (
	"fmt"
	utils "github.com/0187773933/MACAddressTracker/v1/utils"
)

func main() {
	config := utils.GetConfig()
	fmt.Println( config )
	local_network := utils.ScanLocalNetwork( config )
	utils.PrintLocalNetwork( config , local_network )
	utils.TrackChanges( config , local_network )
}