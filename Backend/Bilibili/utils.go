package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	jsoniter "github.com/json-iterator/go"
	//"github.com/ShigemoriHakura/BilibiliDanmu"
)

func getBUserPhoto(id int64) (string, error){
	client := &http.Client{Timeout: time.Second}
	var str =  strconv.Itoa(int(id))
	var url = "https://api.bilibili.com/x/space/acc/info?mid=" + str
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
			log.Println(err)
			return "", err
	}

	req.Header.Set("User-Agent", "Chrome/83.0.4103.61")

	resp, err := client.Do(req)
	if err != nil {
			log.Println(err)
			return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if(err != nil){
			return "", err
	}

	any := jsoniter.Get(body)
	var avatar = any.Get("data", "face").ToString()
	if(avatar != ""){
			log.Printf("B UserId(%v) match: %v", str, avatar)
			return avatar, nil
	}
	return "", nil
}

func Arrcmp(src []string, dest []string) ([]string, []string) {
	msrc := make(map[string]byte) //按源数组建索引
	mall := make(map[string]byte) //源+目所有元素建索引

	var set []string //交集

	//1.源数组建立map
	for _, v := range src {
		msrc[v] = 0
		mall[v] = 0
	}
	//2.目数组中，存不进去，即重复元素，所有存不进去的集合就是并集
	for _, v := range dest {
		l := len(mall)
		mall[v] = 1
		if l != len(mall) { //长度变化，即可以存
			l = len(mall)
		} else { //存不了，进并集
			set = append(set, v)
		}
	}
	//3.遍历交集，在并集中找，找到就从并集中删，删完后就是补集（即并-交=所有变化的元素）
	for _, v := range set {
		delete(mall, v)
	}
	//4.此时，mall是补集，所有元素去源中找，找到就是删除的，找不到的必定能在目数组中找到，即新加的
	var added, deleted []string
	for v := range mall {
		_, exist := msrc[v]
		if exist {
			deleted = append(deleted, v)
		} else {
			added = append(added, v)
		}
	}

	return added, deleted
}

func checkComments(comment string) bool {
	for _, word := range BanString {
		if strings.Contains(comment, word) {
			return true
		}
	}
	return false
}

func getUserMark(uid int64) string {
	uidString := strconv.FormatInt(uid, 10)
	userMark, ok := UserMarks[uidString]
	if ok {
		return userMark
	}
	return ""
}

func trimLastChar(s string) string {
	r, size := utf8.DecodeLastRuneInString(s)
	if r == utf8.RuneError && (size == 0 || size == 1) {
		size = 0
	}
	return s[:len(s)-size]
}
