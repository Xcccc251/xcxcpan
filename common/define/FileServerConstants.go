package define

var SERVER_1 = "server1"
var SERVER_2 = "server2"

var Server map[string]int = map[string]int{
	SERVER_1: 1,
	SERVER_2: 2,
}

func GetServerId(serverName string) int {
	return Server[serverName]
}
