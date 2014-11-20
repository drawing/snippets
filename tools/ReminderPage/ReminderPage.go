package main

import (
	"./sender"
	"code.google.com/p/goprotobuf/proto"
	"encoding/hex"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"tencent_im_reminder"
	"time"
)

func GetFormUint64(r *http.Request, key string) *uint64 {
	var iValue uint64
	var err error
	iValue, err = strconv.ParseUint(r.FormValue(key), 10, 64)
	if err != nil {
		iValue = 0
	}
	return proto.Uint64(iValue)
}

func GetFormUint64WithNil(r *http.Request, key string) *uint64 {
	var iValue uint64
	var err error
	iValue, err = strconv.ParseUint(r.FormValue(key), 10, 64)
	if err != nil {
		return nil
	}
	if iValue == 0 {
		return nil
	}
	return proto.Uint64(iValue)
}

func GetFormUint32(r *http.Request, key string) *uint32 {
	var iValue uint64
	var err error
	iValue, err = strconv.ParseUint(r.FormValue(key), 10, 32)
	if err != nil {
		iValue = 0
	}
	return proto.Uint32(uint32(iValue))
}

func GetFormUint32WithNil(r *http.Request, key string) *uint32 {
	var iValue uint64
	var err error
	iValue, err = strconv.ParseUint(r.FormValue(key), 10, 32)
	if err != nil {
		iValue = 0
	}
	if iValue == 0 {
		return nil
	}
	return proto.Uint32(uint32(iValue))
}

func GetFormString(r *http.Request, key string) *string {
	value := r.FormValue(key)
	return proto.String(value)
}

func GetFormStringWithNil(r *http.Request, key string) *string {
	value := r.FormValue(key)
	if value == "" {
		return nil
	}
	return proto.String(value)
}

func ConstructPackage(op_type uint64,
	head *tencent_im_reminder.PkgHead,
	reminder *tencent_im_reminder.Reminder) (*tencent_im_reminder.ReminderPackage, error) {
	pkg := &tencent_im_reminder.ReminderPackage{}

	switch *reminder.BussiType {
	case 0:
		head.Password = proto.String("7ee70c78beba9fdd119d8eb72806a5c1fccdf427")
	case 1:
		head.Password = proto.String("JD_df02e302689ccf03")
	case 2:
		head.Password = proto.String("GROUP_74656f64746fSt6f6d5f745f63")
	case 3:
		head.Password = proto.String("JD_312a2a23975959c712f28471a8b465d13")
	case 4:
		head.Password = proto.String("Call_142821b4b1f3077de76d766682519ef4")
	case 5:
		head.Password = proto.String("BIRTH_06b904887f6437d0c8b9fc37971f4014")
	}

	request := &tencent_im_reminder.Request{}
	request.Operation = new(tencent_im_reminder.Request_REQ_OP)
	*request.Operation = tencent_im_reminder.Request_REQ_OP(op_type)

	switch tencent_im_reminder.Request_REQ_OP(op_type) {
	case tencent_im_reminder.Request_OP_ADD:
		request.Add = &tencent_im_reminder.Request_Add{}
		request.Add.Reminder = reminder
	case tencent_im_reminder.Request_OP_REMOVE:
		request.Remove = &tencent_im_reminder.Request_Remove{}
		request.Remove.Reminder = reminder
	case tencent_im_reminder.Request_OP_UPDATE:

		/*err := proto.SetExtension(reminder, tencent_im_reminder.E_ForceReset, proto.Uint32(uint32(1)))
		if err != nil {
			log.Println("SetExtension err:", err)
		}*/
		request.Update = &tencent_im_reminder.Request_Update{}
		request.Update.From = reminder
		request.Update.To = reminder
	case tencent_im_reminder.Request_OP_GET:
		request.Get = &tencent_im_reminder.Request_Get{}
		request.Get.Reminder = reminder
	case tencent_im_reminder.Request_OP_DISABLE:
		request.Disable = &tencent_im_reminder.Request_Disable{}
		request.Disable.Reminder = reminder
	case tencent_im_reminder.Request_OP_ENABLE:
		request.Enable = &tencent_im_reminder.Request_Enable{}
		request.Enable.Reminder = reminder
	case tencent_im_reminder.Request_OP_CHECK_ENABLE:
		request.CheckEnable = &tencent_im_reminder.Request_CheckEnable{}
		request.CheckEnable.Reminder = reminder
	case tencent_im_reminder.Request_OP_PUSH_MESSAGE_TO_USER:
		request.Reminder = reminder
	}

	pkg.Head = head
	pkg.Request = append(pkg.Request, request)
	return pkg, nil
}

func SenderHandler(w http.ResponseWriter, r *http.Request) {
	reminder := &tencent_im_reminder.Reminder{}

	reminder.BussiType = GetFormUint32(r, "busi_type")
	reminder.Seq = GetFormUint64(r, "seq")
	reminder.AtTime = GetFormUint64(r, "at_time")

	reminder.Content = GetFormString(r, "content")
	reminder.Lang = GetFormStringWithNil(r, "lang")

	reminder.RedirectUrl = GetFormStringWithNil(r, "redirect_url")
	reminder.ReminderNick = GetFormStringWithNil(r, "reminder_nick")
	reminder.PcTipsTitle = GetFormStringWithNil(r, "pc_tips_title")
	reminder.PcTipsRedirectUrl = GetFormStringWithNil(r, "pc_tips_redirect_url")

	reminder.FromUser = &tencent_im_reminder.Reminder_User{}
	reminder.ToUser = &tencent_im_reminder.Reminder_User{}

	reminder.FromUser.Uin = GetFormUint64(r, "from_uin")

	reminder.AssociationSeq = GetFormUint64WithNil(r, "association_seq")

	user_type := r.FormValue("to_type")
	switch user_type {
	case "C2C":
		reminder.ToUser.Uin = GetFormUint64(r, "to_uin")
	case "GROUP":
		reminder.ToUser.GroupCode = GetFormUint64(r, "to_uin")
	case "DISCUSS":
		reminder.ToUser.DiscussUin = GetFormUint64(r, "to_uin")
	}

	op_type := GetFormUint64(r, "op_type")
	server_type := GetFormUint32(r, "server_type")

	head := &tencent_im_reminder.PkgHead{}
	head.Uin = GetFormUint64(r, "head_uin")
	head.BussiType = GetFormUint32(r, "busi_type")

	head.AuthMethod = GetFormUint32(r, "login_type")
	switch *head.AuthMethod {
	case 1:
		head.Sid = GetFormString(r, "login_string")
	case 2:
		head.Skey = GetFormString(r, "login_string")
	case 3:
		svalue := r.FormValue("login_string")
		bvalue, err := hex.DecodeString(svalue)
		if err != nil {
			log.Println("hex.Decode", err)
		} else {
			head.AuthA2 = proto.String(string(bvalue))
		}
	}
	head.ClientIp = GetFormStringWithNil(r, "client_ip")
	head.ClientPort = GetFormUint32WithNil(r, "client_port")
	head.ClientAppid = GetFormUint32WithNil(r, "client_appid")

	resp := &tencent_im_reminder.ReminderPackage{}
	var cost int64 = 0

	pkg, err := ConstructPackage(*op_type, head, reminder)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	log.Println("REQ:", pkg)

	resp, cost = sender.SendPkg(*server_type, pkg)

	log.Println("RESP:", resp)

	tmpl, err := template.ParseFiles("template/result.html")
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	data := map[string]interface{}{}

	data["RawPackage"] = proto.MarshalTextString(resp)
	data["Cost"] = cost / 1000 / 1000

	data["Content"] = ""
	data["AtTime"] = ""

	if len(resp.Response) >= 1 && len(resp.Response[0].Reminder) >= 1 {
		data["Content"] = resp.Response[0].Reminder[0].Content

		if resp.Response[0].Reminder[0].AtTime != nil {
			data["AtTime"] = time.Unix(int64(*resp.Response[0].Reminder[0].AtTime), 0)
		}
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}

func main() {
	http.HandleFunc("/sender", SenderHandler)

	http.Handle("/page/", http.StripPrefix("/page/", http.FileServer(http.Dir("page"))))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
