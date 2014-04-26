package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

/*
https://openapi.baidu.com/oauth/2.0/authorize?client_id=XXXXXXXXXXXXXXXXXXXXXXXX&response_type=token&redirect_uri=oob&scope=netdisk
*/

type SyncFile struct {
	Path  string
	CTime int
	MTime int
	IsDir int
}

type SyncFS struct {
	AccessToken string
	Client      *http.Client
}

func (fs *SyncFS) Dir(path string) []SyncFile {
	resp, err := http.Get("https://pcs.baidu.com/rest/2.0/pcs/file?method=list&path=" + url.QueryEscape(path) + "&access_token=" + fs.AccessToken)
	if err != nil {
		log.Println("Dir Error:", err)
	}

	// body, _ := ioutil.ReadAll(resp.Body)
	// log.Println(string(body))

	type JsonFile struct {
		Path  string `json:"path"`
		CTime int    `json:"ctime"`
		MTime int    `json:"mtime"`
		IsDir int    `json:"isdir"`
	}
	var JsonFileList struct {
		List []JsonFile `json:"list"`
	}

	var result []SyncFile
	if err = json.NewDecoder(resp.Body).Decode(&JsonFileList); err != nil {
		log.Println("JsonDecode Failed: ", err)
		return result
	}
	// log.Println("succ", JsonFileList)
	for _, v := range JsonFileList.List {
		var file SyncFile
		file.Path = v.Path
		file.CTime = v.CTime
		file.MTime = v.MTime
		file.IsDir = v.IsDir
		result = append(result, file)
	}
	return result
}

func (fs *SyncFS) Space() {
	resp, err := http.Get("https://pcs.baidu.com/rest/2.0/pcs/quota?method=info&access_token=" + fs.AccessToken)
	if err != nil {
		log.Println("Space Error:", err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	log.Println(string(body))
}

func (fs *SyncFS) Move(FromPath string, ToPath string) {
	resp, err := fs.Client.Post("https://pcs.baidu.com/rest/2.0/pcs/file?method=move&access_token="+fs.AccessToken+
		"&from="+url.QueryEscape(FromPath)+"&to="+url.QueryEscape(ToPath), "", nil)

	defer resp.Body.Close()

	if err != nil {
		log.Println("Move Error:", err)
		return
	}

	var move_status struct {
		ErrCode int `json:"error_code"`
	}

	// body, _ := ioutil.ReadAll(resp.Body)
	// log.Println(string(body))
	if err = json.NewDecoder(resp.Body).Decode(&move_status); err != nil {
		log.Println("Move JsonDecode Failed: ", FromPath, err)
		return
	}
	if move_status.ErrCode != 0 {
		log.Println("MoveError: from", FromPath, move_status.ErrCode)
	}
}
