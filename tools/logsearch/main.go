package main

import (
	"./remote"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
)

func crawlHandler(w http.ResponseWriter, r *http.Request) {
	uin := r.FormValue("uin")
	input_time := r.FormValue("time")
	op_trace := r.FormValue("trace")
	op_sendmsg := r.FormValue("sendmsg")

	log.Println("uin=", uin, ", time=", input_time)
	const layout = "2006-1-2_15"
	access_time, err := time.Parse(layout, input_time)

	var color string
	var sendmsg string
	var status string
	var ip string

	if op_trace == "on" {
		color, err = remote.GetTraceFromLogServer(uin, access_time)
		if err != nil {
			color = err.Error()
		}
	}

	if op_sendmsg == "on" {
		status, ip, err = remote.GetUserCurrentStatus(uin)
		log.Println("Status", status, ip, err)

		if err == nil {
			sendmsg, err = remote.GetTraceFromOnlineServer(uin, ip, access_time)
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

func main() {
	http.HandleFunc("/crawl", crawlHandler)
	http.Handle("/page/", http.StripPrefix("/page/", http.FileServer(http.Dir("page"))))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
