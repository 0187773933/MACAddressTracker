package utils

import (
	"os"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"path/filepath"
	"strconv"
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
	CronString string `json:"cron_string"`
	Latitude string "json:latitude"
	Longitude string "json:longitude"
	ServerPort string `json:"server_port"`
	SavedRecordTotal int `json:"saved_record_total"`
	NetworkHardWareInterfaceName string `json:"network_hardware_interface_name"`
	Redis RedisConnectionInfo `json:"redis"`
	Devices map[string]string `json:"devices"`
}

func ParseConfig( file_path string ) ( result ConfigFile ) {
	file_data , _ := ioutil.ReadFile( file_path )
	err := json.Unmarshal( file_data , &result )
	if err != nil { fmt.Println( err ) }
	cleaned_devices := map[string]string{}
	for mac , device_name := range result.Devices {
		cleaned_mac := strings.ToLower( strings.Join( strings.Split( mac , "-" ) , ":" ) )
		cleaned_devices[  cleaned_mac ] = device_name
	}
	result.Devices = cleaned_devices
	return
}
func ParseConfigENV() ( result ConfigFile ) {
	result.LocationName = os.Getenv( "MAC_LOCATION_NAME" )
	result.CronString = os.Getenv( "MAC_CRON_STRING" )
	result.Latitude = os.Getenv( "MAC_LATITUDE" )
	result.Longitude = os.Getenv( "MAC_LONGITUDE" )
	result.ServerPort = os.Getenv( "MAC_SERVER_PORT" )
	saved_record_total , _ := strconv.Atoi( os.Getenv( "MAC_SAVED_RECORD_TOTAL" ) )
	result.SavedRecordTotal = saved_record_total
	result.NetworkHardWareInterfaceName = os.Getenv( "MAC_NETWORK_HARDWARE_INTERFACE_NAME" )
	result.Redis.Host = os.Getenv( "MAC_REDIS_HOST" )
	result.Redis.Port = os.Getenv( "MAC_REDIS_PORT" )
	db , _ := strconv.Atoi( os.Getenv( "MAC_REDIS_DB" ) )
	result.Redis.DB = db
	result.Redis.Password = os.Getenv( "MAC_REDIS_PASSWORD" )
	result.Redis.Prefix = os.Getenv( "MAC_REDIS_PREFIX" )
	return
}

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
	fmt.Printf( "\nScanning Local Network\n" )
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
			if os.Getenv( "MAC_LOCATION_NAME" ) == "" {
				fmt.Println( "Pass filepath as argv1 or populate ~/.config/personal/mac_address_tracker.json" )
				panic( "Can't Locate Config File Anywhere" )
			}
			config_file_path = "ENV"
		}
	} else {
		config_file_path , _ = filepath.Abs( os.Args[ 1 ] )
	}
	if config_file_path == "ENV" {
		config = ParseConfigENV()
	} else {
		config = ParseConfig( config_file_path )
	}
	if config.CronString == "" { config.CronString = "@every 5m" }
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
	DeviceName string `json:"device_name"`
	CurrentMACAddress string `json:"current_mac_address"`
	CurrentIPAddress string `json:"current_ip_address"`
	CurrentTimeString string `json:"current_time_string"`
	CurrentLocation string `json:"current_location"`
	Records []string `json:"records"`
	Transitions []string `json:"transitions"`
	Departures []string `json:"departures"`
	Arrivals []string `json:"arrivals"`
}
func RedisGetSeenMACAddress( redis redis.Manager , db_item_key string ) ( db_item MacAddressRecord ) {
	record := MacAddressRecord{}
	db_item = record
	if RedisKeyExists( redis , db_item_key ) == true {
		existing_db_entry_json := redis.Get( db_item_key )
		json_unmarshal_error := json.Unmarshal( []byte( existing_db_entry_json ) , &record )
		if json_unmarshal_error != nil { return } else { db_item = record }
	}
	return
}
func TrackChanges( config ConfigFile , network_map [][2]string ) {

	fmt.Println( "Tracking Changes" )
	redis := GetRedisConnection( config.Redis )
	var ctx = context.Background()
	all_seen_at_time := time.Now()
	all_seen_at_time_string := GetFormattedTimeString( all_seen_at_time )
	var record_cutoff int
	if config.SavedRecordTotal == 0 { record_cutoff = 100 } else { record_cutoff = config.SavedRecordTotal }

	// 1.) Get State Variables / Reset All 'Snapshots' /
	// 1.A.) Build Keys
	network_latest_ip_set_key := fmt.Sprintf( "%sNETWORK.%s.LATEST.IPS" , config.Redis.Prefix , config.LocationName )
	network_latest_mac_set_key := fmt.Sprintf( "%sNETWORK.%s.LATEST.MACS" , config.Redis.Prefix , config.LocationName )
	// 1.B.) Cache any Needed
	network_latest_mac_set_cached , _ := redis.Redis.SMembersMap( ctx , network_latest_mac_set_key ).Result()
	// 1.C.) Delete Previous State
	RedisKeyDelete( redis , network_latest_ip_set_key )
	RedisKeyDelete( redis , network_latest_mac_set_key )

	// 2.) Detect Departures
	network_map_actual_map := map[string]string{}
	for index := range network_map { network_map_actual_map[ network_map[ index ][ 1 ] ] = network_map[ index ][ 0 ] }
	for mac , _ := range network_latest_mac_set_cached {
		_ , exists_in_latest := network_map_actual_map[ mac ]
		if exists_in_latest == false {
			db_item_key := fmt.Sprintf( "%sSEEN.%s" , config.Redis.Prefix , mac )
			db_item := RedisGetSeenMACAddress( redis , db_item_key )
			departure_string := fmt.Sprintf( "%s === %s === %s === %s === DEPARTED === FROM === %s" , all_seen_at_time_string , db_item.CurrentIPAddress , db_item.CurrentMACAddress , db_item.DeviceName , config.LocationName )
			fmt.Println( departure_string )
			// fmt.Printf( "B.) %s === %s === DEPARTED\n" , mac , network_map_actual_map[ mac ] )
			db_item.Departures = append( db_item.Departures , departure_string )
			json_string := JSONStringify( db_item )
			redis.Set( db_item_key , json_string )
		}
	}

	// 3.) 'Track Changes' For each Item in the Current Snapshot of the Local Network
	for index := range network_map {

		ip := network_map[ index ][ 0 ]
		mac := network_map[ index ][ 1 ]
		arrived := false

		// 3.A.) Grab Previous DB Entry If it Exists
		db_item_key := fmt.Sprintf( "%sSEEN.%s" , config.Redis.Prefix , mac )
		record := MacAddressRecord{}
		if RedisKeyExists( redis , db_item_key ) == true {
			existing_db_entry_json := redis.Get( db_item_key )
			json_unmarshal_error := json.Unmarshal( []byte( existing_db_entry_json ) , &record )
			if json_unmarshal_error != nil { panic( json_unmarshal_error ) }
			// 3.B) Detect Transition on DB Entry
			if record.CurrentLocation != "" {
				if record.CurrentLocation != config.LocationName {
					transition_string := fmt.Sprintf( "%s === %s === %s === %s === TRANSITIONED === FROM === %s === TO === %s" , all_seen_at_time_string , record.CurrentIPAddress , record.CurrentMACAddress , record.DeviceName , record.CurrentLocation , config.LocationName )
					fmt.Println( transition_string )
					record.Transitions = append( record.Transitions , transition_string )
					record.CurrentLocation = config.LocationName
				}
			}
		}

		// 3.C.) Manual Updates of DB Entry Values
		// 3.C.1.) Misc Values
		if config.Devices[ mac ] != "" { record.DeviceName = config.Devices[ mac ] }
		record.CurrentTimeString = all_seen_at_time_string
		record.CurrentLocation = config.LocationName
		record_string := fmt.Sprintf( "%s === %s === %s === %s === %s" , all_seen_at_time_string , config.LocationName , ip , mac , record.DeviceName )
		record.Records = append( record.Records , record_string )
		if len( record.Records ) > record_cutoff { record.Records = record.Records[1:] }
		record.CurrentMACAddress = mac
		record.CurrentIPAddress = ip

		// 3.C.2.) Detect Arrival
		_ , previously_existed := network_latest_mac_set_cached[ mac ]
		if previously_existed == false {
			arrived = true
			arrival_record_string := fmt.Sprintf( "%s === %s === %s === %s === ARRIVED AT === %s" , all_seen_at_time_string , record.CurrentIPAddress , record.CurrentMACAddress , record.DeviceName , config.LocationName )
			// fmt.Println( arrival_record_string )
			record.Arrivals = append( record.Arrivals , arrival_record_string )
			if len( record.Arrivals ) > record_cutoff { record.Arrivals = record.Arrivals[1:] }
		}

		// 3.C.3.) Save Updated DB Entry into Redis
		json_string := JSONStringify( record )
		redis.Set( db_item_key , json_string )

		// 3.D.) Add 'this' DB Entry's IP and MAC Address to Current Location's Sets
		// These are usefull for dicitonary-lookup-existance, instead of whole list
		redis.Redis.SAdd( ctx , network_latest_ip_set_key , ip )
		redis.Redis.SAdd( ctx , network_latest_mac_set_key , mac )

		// 3.E.) Print DB Item
		print_string := fmt.Sprintf( "%02d === %s === %s === %s === %s === %s" , index , all_seen_at_time_string , config.LocationName , ip , mac , record.DeviceName )
		if arrived == true { print_string = fmt.Sprintf( "%s \t\t=== ARRIVED" , print_string ) }
		fmt.Println( print_string )
	}

}