package utils

import (
	"os"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"path/filepath"
	"context"
	"time"
	"strings"
	// redis_lib "github.com/go-redis/redis/v7"
	redis "github.com/0187773933/RedisManagerUtils/manager"
	scanner "github.com/0187773933/MACAddressTracker/v1/scanner"
)

type RedisConnectionInfo struct {
	Host string `json:"host"`
	Port string `json:"port"`
	DB int `json:"db"`
	Password string `json:"password"`
	Prefix string `json:"prefix"`
}
type ConfigFile struct {
	LocationName string `json:"location_name"`
	Latitude string "json:latitude"
	Longitude string "json:longitude"
	ServerPort string `json:"server_port"`
	NetworkHardWareInterfaceName string `json:"network_hardware_interface_name"`
	Redis RedisConnectionInfo `json:"redis"`
}
// https://github.com/0187773933/VizioController/blob/master/controller/viziocontroller.go#L133
// func ParseConfig( file_path string ) ( result interface{} ) {
// 	file_data , _ := ioutil.ReadFile( file_path )
// 	err := json.Unmarshal( file_data , &result )
// 	if err != nil { fmt.Println( err ) }
// 	return
// }
func ParseConfig( file_path string ) ( result ConfigFile ) {
	file_data , _ := ioutil.ReadFile( file_path )
	err := json.Unmarshal( file_data , &result )
	if err != nil { fmt.Println( err ) }
	return
}

// func GetRedisConnection( info RedisConnectionInfo ) ( redis_client *redis_lib.Client ) {
// 	redis_client = redis_lib.NewClient( &redis_lib.Options{
// 		Addr: fmt.Sprintf( "%s:%s" , info.Host , info.Port ) ,
// 		DB: info.DB ,
// 		Password: info.Password ,
// 	})
// 	return
// }
func GetRedisConnection( info RedisConnectionInfo ) ( redis_client redis.Manager ) {
	redis_client.Connect( fmt.Sprintf( "%s:%s" , info.Host , info.Port ) , info.DB , info.Password )
	return
}

func PrintLocalNetwork( config ConfigFile , network_map [][2]string ) {
	for index := range network_map {
		mac_hostname_key := fmt.Sprintf( "%sNETWORK.%s.%s" , config.Redis.Prefix , config.LocationName , network_map[index][1] )
		fmt.Printf( "%d === %s === %s === %s\n" , index , network_map[index][0] , network_map[index][1] , mac_hostname_key )
	}
}

func GetFormattedTimeString( time_object time.Time ) ( result string ) {
	// https://stackoverflow.com/a/51915792
	// month_name := strings.ToUpper( time_object.Format( "Feb" ) ) // monkaHmm
	month_name := strings.ToUpper( time_object.Format( "Jan" ) )
	milliseconds := time_object.Format( ".000" )
	date_part := fmt.Sprintf( "%02d%s%d" , time_object.Day() , month_name , time_object.Year() )
	time_part := fmt.Sprintf( "%02d:%02d:%02d%s" , time_object.Hour() , time_object.Minute() , time_object.Second() , milliseconds )
	result = fmt.Sprintf( "%s === %s" , date_part , time_part )
	return
}

func ScanLocalNetwork( config ConfigFile ) ( network_map [][2]string ) {
	fmt.Println( "Scanning Local Network" )
	// https://stackoverflow.com/questions/40260599/difference-between-two-time-time-objects
	start_time := time.Now()
	start_time_string := GetFormattedTimeString( start_time )
	fmt.Println( start_time_string )
	fmt.Println( start_time.Location )
	network_map = scanner.ScanLocalNetwork( config.NetworkHardWareInterfaceName )
	end_time := time.Now()
	delta_time := end_time.Sub( start_time )
	fmt.Println( delta_time )
	end_time_string := GetFormattedTimeString( end_time )
	fmt.Println( end_time_string )
	return
}

func FileExists( filename string ) bool {
	info , err := os.Stat( filename )
	if os.IsNotExist( err ) {
		return false
	}
	return !info.IsDir()
}

func GetConfig() ( config ConfigFile ) {
	fmt.Println( "Finding Config File" )
	var config_file_path string
	if len( os.Args ) < 2 {
		home_directory , _ := os.UserHomeDir()
		config_file_path = filepath.Join( home_directory , ".config" , "personal" , "mac_address_tracker.json" )
		if FileExists( config_file_path ) == false {
			fmt.Println( "Pass filepath as argv1 or populate ~/.config/personal/mac_address_tracker.json" )
			panic( "Can't Locate Config File Anywhere" )
		}
	} else {
		config_file_path , _ = filepath.Abs( os.Args[ 1 ] )
	}
	config = ParseConfig( config_file_path )
	fmt.Println( config_file_path )
	return
}

func RedisKeyExists( redis redis.Manager , redis_key string ) ( result bool ) {
	result = false
	var ctx = context.Background()
	exists_result , exists_error := redis.Redis.Exists( ctx , redis_key ).Result()
	if exists_error != nil { fmt.Println( exists_error ); }
	if exists_result == 1 { result = true }
	return
}
func RedisKeyDelete( redis redis.Manager , redis_key string ) {
	if RedisKeyExists( redis , redis_key ) == true {
		var ctx = context.Background()
		_ , delete_error := redis.Redis.Del( ctx , redis_key ).Result()
		if delete_error != nil { fmt.Println( delete_error ); }
		// fmt.Println( delete_result )
	}
}

func JSONStringify( object interface{} ) ( json_string string ) {
	json_marshal_result , json_marshal_error := json.Marshal( object )
	if json_marshal_error != nil { panic( json_marshal_error ) }
	json_string = string( json_marshal_result )
	return
}

type MacAddressRecord struct {
	CurrentTimeString string `json:"current_time_string"`
	Records []string `json:"records"`
	Transitions []string `json:"transitions"`
}
func TrackChanges( config ConfigFile , network_map [][2]string ) {
	fmt.Println( "Tracking Changes" )
	redis := GetRedisConnection( config.Redis )
	var ctx = context.Background()

	// 0.) Reset All 'Snapshots'
	network_latest_set_key := fmt.Sprintf( "%sNETWORK.%s.IPS.LATEST" , config.Redis.Prefix , config.LocationName )
	RedisKeyDelete( redis , network_latest_set_key )

	all_seen_at_time := time.Now()
	all_seen_at_time_string := GetFormattedTimeString( all_seen_at_time )
	for index := range network_map {
		ip := network_map[index][0]
		mac := network_map[index][1]
		mac_hostname_key := fmt.Sprintf( "%sNETWORK.%s.%s" , config.Redis.Prefix , config.LocationName , mac )
		fmt.Printf( "%d === %s === %s === %s\n" , index , ip , mac , mac_hostname_key )

		// Build A DB Item with MetaData and Stuff , Or Update Previously Existing Entry
		db_item_key := fmt.Sprintf( "%sSEEN.%s" , config.Redis.Prefix , mac )
		var record MacAddressRecord
		if RedisKeyExists( redis , db_item_key ) == true {
			existing_db_entry_json := redis.Get( db_item_key )
			json_unmarshal_error := json.Unmarshal( []byte( existing_db_entry_json ) , &record )
			if json_unmarshal_error != nil { panic( json_unmarshal_error ) }
			// Perform Update

		} else {
			record = MacAddressRecord{}
		}
		record.CurrentTimeString = all_seen_at_time_string
		fmt.Println( db_item_key )
		// json_marshal_result , json_marshal_error := json.Marshal( test_data )
		// fmt.Println( reflect.TypeOf( json_marshal_result ) )
		// if json_marshal_error != nil { panic( json_marshal_error ) }
		// json_string := string( json_marshal_result )
		json_string := JSONStringify( record )
		redis.Set( db_item_key , json_string )

		// 1.) Store Snapshot as Set of 'Latest' IP's
		redis.Redis.SAdd( ctx , network_latest_set_key , ip )

		// 2.) Store Snapshot of S
		// redis.ListPushRight(  )
	}
}