package main

import (
	"./remote"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

func GetTraceFromLogServer(user string, access_time string) (string, error) {
	shell, err := remote.NewRemote("mqq@10.187.147.94", "mqq2005")
	if err != nil {
		return "", err
	}

	command := fmt.Sprintf("grep -a \"%s\" /usr/local/app/LogServer/log/pushsvc/trace.log.%s00", user, access_time)
	log.Println(command)
	color, err := shell.RunCommand(command)
	if err != nil {
		return "", err
	}

	if strings.Contains(color, "No such file or directory") {
		command := fmt.Sprintf("zgrep \"%s\" /usr/local/app/LogServer/log/pushsvc/trace.log.%s00.gz", user, access_time)
		log.Println(command)
		color, err = shell.RunCommand(command)
	}

	shell.Close()

	return color, err
}

func GetTraceFromOnlineServer(user string, ip string, access_time string) (string, error) {
	shell, err := remote.NewRemote("mqq@"+ip, "mqq2005")
	if err != nil {
		return "", err
	}

	command := fmt.Sprintf("grep -a \"|%s|\" /usr/local/app/taf/app_log/QQService/PushSvc%d/QQService.PushSvc%d_SendMsg_%s.log",
		user, user[len(user)-3]-'0', user[len(user)-3]-'0', access_time)
	log.Println(command)
	color, err := shell.RunCommand(command)
	if err != nil {
		return "", err
	}

	if strings.Contains(color, "No such file or directory") {
		command := fmt.Sprintf("zgrep -a \"|%s|\" /usr/local/app/taf/app_log/QQService/PushSvc%d/QQService.PushSvc%d_SendMsg_%s.log.gz",
			user, user[len(user)-3]-'0', user[len(user)-3]-'0', access_time)
		log.Println(command)
		color, err = shell.RunCommand(command)
	}

	shell.Close()

	return color, err
}

func logHandler(w http.ResponseWriter, r *http.Request) {
	uin := r.FormValue("uin")
	access_time := r.FormValue("time")
	op_trace := r.FormValue("trace")
	op_sendmsg := r.FormValue("sendmsg")

	log.Println("uin=", uin, ", time=", access_time)

	var color string
	var err error
	var sendmsg string
	var status string
	var ip string

	if op_trace == "on" {
		color, err = GetTraceFromLogServer(uin, access_time)
		if err != nil {
			color = err.Error()
		}
	}

	if op_sendmsg == "on" {
		status, ip, err = remote.GetUserCurrentStatus(uin)
		log.Println("Status", status, ip, err)

		if err == nil {
			sendmsg, err = GetTraceFromOnlineServer(uin, ip, access_time)
			if err != nil {
				sendmsg = err.Error()
			}
		} else {
			sendmsg = err.Error()
		}
	}

	tmpl, err := template.ParseFiles("log.html")
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	data := map[string]interface{}{}

	data["TraceLogs"] = strings.Split(color, "\n")
	data["Status"] = status
	data["IP"] = ip
	data["SendMsg"] = strings.Split(sendmsg, "\n")

	err = tmpl.Execute(w, data)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	data := map[string]interface{}{}
	err = tmpl.Execute(w, data)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}

func main() {
	http.HandleFunc("/log", logHandler)
	http.HandleFunc("/index.html", mainHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
