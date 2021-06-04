package scanner

import (
	"fmt"
	"net"
	//"io"
	//"bytes"
	//"os"
	"os/exec"
	//"strconv"
	"strings"
	"sort"
	"strconv"
	"runtime"
	default_gateway "github.com/jackpal/gateway"
	//"github.com/mdlayher/arp"
	//"github.com/mdlayher/ethernet"
)

type LocalNetwork struct {
	interfaces map[string]map[string]string
	default_gateway_ip string
	local_ip string
	public_ip string
}

type ArpResult map[string] string

func exec_process( bash_command string , arguments ...string ) ( result string ) {
	command := exec.Command( bash_command , arguments... )
	//command.Env = append( os.Environ() , "DISPLAY=:0.0" )
	out, err := command.Output()
	if err != nil {
		fmt.Println( bash_command )
		fmt.Println( arguments )
		fmt.Sprintf( "%s\n" , err )
	}
	result = string( out[:] )
	return
}

// Exec Function Style 2
// CombinedOutput???
func get_net_mask(deviceName string) string {
	switch runtime.GOOS {
	case "darwin":
		cmd := exec.Command("ipconfig", "getoption", deviceName, "subnet_mask")
		out, err := cmd.CombinedOutput()
		if err != nil {
			return ""
		}
		nm := strings.Replace(string(out), "\n", "", -1)
		fmt.Printf("netmask=%s OS=%s", nm, runtime.GOOS)
		return nm
	default:
		return ""
	}
	return ""
}

func probe_local_network() LocalNetwork {
	local_network := LocalNetwork{}
	local_network.interfaces = make( map[string]map[string]string )
	interfaces , _ := net.Interfaces()
	//address_info , _ := net.InterfaceAddrs()
	//fmt.Println( address_info )
	for _ , x_interface := range interfaces {
		local_network.interfaces[x_interface.Name] = make( map[string]string )
		addresses , error := x_interface.Addrs()
		if error != nil { continue }
		//fmt.Println( x_interface.HardwareAddr ) // aka the mac address of the interface?
		for _ , address := range addresses {
			var ip net.IP
			var mask net.IPMask
			switch v := address.(type) {
			case *net.IPNet:
				ip = v.IP
				local_network.interfaces[x_interface.Name]["our_ip"] = v.IP.String()
				mask = v.Mask
				//fmt.Println( net.ParseCIDR( local_network.interfaces[x_interface.Name]["our_ip"] ) )
			case *net.IPAddr:
				ip = v.IP
				local_network.interfaces[x_interface.Name]["our_ip"] = v.IP.String()
				mask = ip.DefaultMask()
				//fmt.Println( net.ParseCIDR( local_network.interfaces[x_interface.Name]["our_ip"] ) )
			}
			if ip == nil { continue }
			ip = ip.To4()
			if ip == nil { continue }
			cleanMask := fmt.Sprintf("%d.%d.%d.%d", mask[0], mask[1], mask[2], mask[3])
			fmt.Println( ip , cleanMask )
		}
		//ifi, err := net.InterfaceByName()
		//ifi := *x
	}
	return local_network
}

func nmap( gateway_ip string ) ( result string ) {
	result = "failed"
	switch runtime.GOOS {
		case "linux":
			nmap_command := fmt.Sprintf( "nmap -sn %s/24" , gateway_ip )
			result := exec_process( "/bin/bash" , "-c" , nmap_command )
			return result
		case "darwin":
			nmap_command := fmt.Sprintf( "nmap -sP %s/24" , gateway_ip )
			result := exec_process( "/bin/bash" , "-c" , nmap_command )
			return result
		case "windows":
			nmap_command := fmt.Sprintf( "nmap -sP %s/24" , gateway_ip )
			result := exec_process( `C:\Windows\System32\cmd.exe` , "/c" , nmap_command )
			return result
	}
	return result
}

func _remove_empty_strings( items []string ) ( results []string ) {
	for _ , item := range items {
		if item == "" { continue }
		results = append( results , item )
	}
	return
}
func arp_interface( interface_name string ) ( arp_result ArpResult ) {
	arp_result = ArpResult{}
	switch runtime.GOOS {
		case "linux":
			arp_string := exec_process( "/bin/bash" , "-c" , fmt.Sprintf( "arp -na -i %s | awk '{{print $2,$4}}'" , interface_name ) )
			lines := strings.Split( arp_string , "\n" )
			for _ , line := range lines {
				items := strings.Split( line , " " )
				if len( items ) < 2 { continue }
				if items[1] == "(incomplete)" { continue }
				if items[1] == "<incomplete>" { continue }
				if strings.Contains( items[0] , "(" ) == false { continue }
				ip_address := strings.Split( strings.Split( items[0] , "(" )[1] , ")" )[0]
				arp_result[ip_address] = items[1]
			}
		case "darwin":
			arp_string := exec_process( "/bin/bash" , "-c" , fmt.Sprintf( "arp -na -i %s | awk '{{print $2,$4}}'" , interface_name ) )
			lines := strings.Split( arp_string , "\n" )
			for _ , line := range lines {
				items := strings.Split( line , " " )
				if len( items ) < 2 { continue }
				if items[1] == "(incomplete)" { continue }
				if items[1] == "<incomplete>" { continue }
				if strings.Contains( items[0] , "(" ) == false { continue }
				ip_address := strings.Split( strings.Split( items[0] , "(" )[1] , ")" )[0]
				arp_result[ip_address] = items[1]
			}
		case "windows":
			arp_string := exec_process( `C:\Windows\System32\cmd.exe` , "/c" , "arp -a" )
			lines := strings.Split( arp_string , "\n" )
			// substitutions := []string{ "-" }
			for _ , line := range lines {
				if strings.Contains( line , "dynamic" ) == false { continue }
				items := strings.Split( line , " " )
				items = _remove_empty_strings( items )
				if len( items ) < 3 { continue }
				ip_address := items[ 0 ]
				mac_address := strings.Join( strings.Split( items[ 1 ] , "-" ) , ":" )
				arp_result[ip_address] = mac_address
			}
	}
	return
}

func GetIPAddressFromMacAddress( interface_name string , mac_address string ) ( ip_address string ) {
	default_gateway_ip , _ := default_gateway.DiscoverGateway()
	nmap( default_gateway_ip.String() )
	arp_result := arp_interface( interface_name )
	ip_address = arp_result[mac_address]
	return
}

func sort_local_network( arp_result ArpResult ) ( network_map [][2]string ) {
	//arp_result["192.168.1.52"] = "b8:27:eb:52:a7:6b"
	var ip_address_ends []int
	for key := range arp_result {
		ip_address_parts := strings.Split( key , "." )
		i , _ := strconv.Atoi( ip_address_parts[ len( ip_address_parts ) - 1 ] )
		ip_address_ends = append( ip_address_ends , i )
	}
	sort.Ints( ip_address_ends )
	for _ , key := range ip_address_ends {
		var new_item [2]string
		new_item[0] = fmt.Sprintf( "192.168.1.%d" , key )
		new_item[1] = arp_result[ fmt.Sprintf( "192.168.1.%d" , key ) ]
		//fmt.Println( new_item )
		network_map = append( network_map , new_item )
		//fmt.Println( fmt.Sprintf( "192.168.1.%d ===" , key ) , network_map[ fmt.Sprintf( "192.168.1.%d" , key ) ] )
	}
	return
}

func ScanLocalNetwork( interface_name string ) ( local_network [][2]string ) {
	default_gateway_ip , _ := default_gateway.DiscoverGateway()
	nmap( default_gateway_ip.String() )
	arp_result := arp_interface( interface_name )
	local_network = sort_local_network( arp_result )
	return
}