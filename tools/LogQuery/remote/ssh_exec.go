package remote

import (
	"bytes"
	"fmt"
	"os/exec"
	"time"
)

func RemoteExec(host string, user string, pass string, cmdline string) (string, error) {
	inst := exec.Command("./ssh2_exec", host, user, pass, cmdline)
	var out bytes.Buffer
	inst.Stderr = &out

	err := inst.Run()
	return out.String(), err
}

func GetTraceFromLogServer(user string, access_time time.Time) (string, error) {
	cmdline := fmt.Sprintf("grep -a \"%s\" /usr/local/app/LogServer/log/pushsvc/trace.log.%s00", user, access_time.Format("2006-01-02_15"))
	out, err := RemoteExec("10.187.147.94", "mqq", "mqq2005", cmdline)

	if err != nil {
		return out, err
	}

	cmdline = fmt.Sprintf("zgrep \"%s\" /usr/local/app/LogServer/log/pushsvc/trace.log.%s00.gz", user, access_time.Format("2006-01-02_15"))
	out2, err := RemoteExec("10.187.147.94", "mqq", "mqq2005", cmdline)

	if len(out2) > len(out) {
		return out2, err
	}

	return out, err
}

func GetTraceFromOnlineServer(user string, ip string, access_time time.Time) (string, error) {
	cmdline := fmt.Sprintf("grep -a \"|%s|\" /usr/local/app/taf/app_log/QQService/PushSvc%d/QQService.PushSvc%d_SendMsg_%s.log",
		user, user[len(user)-3]-'0', user[len(user)-3]-'0', access_time.Format("2006-01-02_15"))
	out, err := RemoteExec(ip, "mqq", "mqq2005", cmdline)

	if err != nil {
		return out, err
	}

	cmdline = fmt.Sprintf("zgrep -a \"|%s|\" /usr/local/app/taf/app_log/QQService/PushSvc%d/QQService.PushSvc%d_SendMsg_%s.log.gz",
		user, user[len(user)-3]-'0', user[len(user)-3]-'0', access_time.Format("2006-01-02_15"))
	out2, err := RemoteExec(ip, "mqq", "mqq2005", cmdline)

	if len(out2) > len(out) {
		return out2, err
	}

	return out, err
}
