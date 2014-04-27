package remote

import (
	"errors"
	"log"
	"net"
	"strings"
)

var PushSvcSet0 []string = []string{
	"10.130.10.228",
	"10.130.92.11",
	"10.134.10.98",
	"10.204.18.227",
	"10.204.18.228",
	"10.135.35.82",
	"10.135.35.83",
	"10.135.36.154",
	"10.135.36.155",
	"10.135.36.156",
}

var PushSvcSet1 []string = []string{
	"10.187.137.31",
	"10.128.69.12",
	"10.135.34.94",
	"10.204.23.154",
	"10.204.23.167",
	"10.135.36.157",
	"10.135.36.158",
	"10.135.36.159",
	"10.135.36.160",
	"10.135.36.161",
}

var PushSvcSet2 []string = []string{
	"10.204.18.91",
	"10.204.18.92",
	"10.128.15.230",
	"10.135.34.181",
	"10.136.170.35",
	"10.135.36.162",
	"10.135.36.163",
	"10.135.38.85",
	"10.135.39.42",
	"10.136.166.151",
}

var PushSvcSet3 []string = []string{
	"10.128.15.218",
	"10.128.36.181",
	"10.128.70.179",
	"172.27.31.24",
	"10.204.18.94",
	"10.130.139.16",
	"10.135.34.93",
	"10.136.166.152",
	"10.136.166.153",
	"10.136.166.154",
	"10.136.166.158",
	"10.137.142.29",
}

var PushSvcSet4 []string = []string{
	"10.153.132.85",
	"10.153.156.142",
	"10.153.156.145",
	"10.153.156.152",
	"10.153.156.249",
	"10.130.9.46",
	"10.130.9.54",
	"10.135.34.180",
	"10.177.147.83",
	"10.130.139.15",
	"10.187.137.30",
	"10.135.36.166",
}

var PushSvcSet5 []string = []string{
	"10.130.139.12",
	"10.130.139.14",
	"10.134.11.99",
	"10.135.34.179",
	"10.177.147.82",
	"10.204.17.239",
	"10.198.131.175",
	"10.198.131.176",
	"10.198.131.177",
	"10.198.132.105",
	"10.198.132.106",
	"10.128.69.10",
}

var PushSvcSet6 []string = []string{
	"10.128.70.115",
	"10.177.147.84",
	"10.204.23.152",
	"10.204.23.153",
	"10.187.151.222",
	"10.187.151.240",
	"10.198.131.174",
	"10.198.131.178",
	"10.198.132.45",
	"10.128.69.115",
}

var PushSvcSet7 []string = []string{
	"10.130.67.98",
	"10.128.15.210",
	"10.130.9.45",
	"10.130.9.62",
	"10.187.151.155",
	"10.187.151.156",
	"10.187.151.209",
	"10.187.151.210",
	"10.187.151.219",
	"10.128.69.17",
}

var PushSvcSet8 []string = []string{
	"10.204.18.93",
	"10.130.67.154",
	"10.130.10.230",
	"10.130.9.79",
	"10.134.11.95",
	"172.27.34.249",
	"10.128.70.178",
	"10.135.32.100",
	"10.135.32.53",
	"10.135.32.75",
	"10.135.35.43",
	"10.135.35.46",
}

var PushSvcSet9 []string = []string{
	"10.198.132.205",
	"10.130.9.44",
	"10.177.147.81",
	"10.135.35.76",
	"10.135.35.77",
	"10.135.35.78",
	"10.135.35.79",
	"10.135.35.81",
	"10.204.18.229",
	"10.204.18.230",
}

var PushSvcIP map[byte][]string = map[byte][]string{
	0: PushSvcSet0,
	1: PushSvcSet1,
	2: PushSvcSet2,
	3: PushSvcSet3,
	4: PushSvcSet4,
	5: PushSvcSet5,
	6: PushSvcSet6,
	7: PushSvcSet7,
	8: PushSvcSet8,
	9: PushSvcSet9,
}

func GetUserCurrentStatus(user string) (status string, ip string, err error) {
	if len(user) <= 3 {
		return "", "", errors.New("user error:" + user)
	}

	var set_id byte = user[len(user)-3] - '0'

	set := PushSvcIP[set_id]
	pkg := make([]byte, 1024)

	for _, v := range set {
		conn, err := net.Dial("tcp", v+":44953")
		if err != nil {
			log.Println("dial:", err)
			continue
		}

		conn.Write([]byte("show " + user))

		n, err := conn.Read(pkg)
		if err != nil {
			log.Println("Read:", err)
			continue
		}

		status = string(pkg[0:n])
		if strings.Contains(status, "can't find user,") {
			continue
		}
		return status, v, nil
		/*
			shell, err := NewRemote("mqq@"+v, "mqq2005")
			if err != nil {
				return "", "", "", err
			}

			status, err := shell.ExecTelnet("telnet "+v+" 44953", "show "+user)
			if err != nil {
				return "", "", "", err
			}

			log.Println(status)

			if strings.Contains(status, "can't find user") {
				log.Println("skip")
				continue
			}
			log.Println("after")

			command := fmt.Sprintf("grep -a %s /usr/local/app/taf/app_log/QQService/PushSvc%d/QQService.PushSvc%d_SendMsg_%s.log",
				user, set_id, set_id, access_time)
			log.Println(command)
			content, err := shell.RunCommand(command)
			if err == nil {
				return "", "", "", err
			}

			if strings.Contains(content, "No such file or directory") {
				command := fmt.Sprintf("zgrep %s /usr/local/app/taf/app_log/QQService/PushSvc%d/QQService.PushSvc%d_SendMsg_%s.log.gz",
					user, set_id, set_id, access_time)
				log.Println(command)
				content, err = shell.RunCommand(command)
				if err == nil {
					return "", "", "", err
				}
			}

			return status, v, content, err*/
	}
	return "", "", errors.New("can't find user")
}
