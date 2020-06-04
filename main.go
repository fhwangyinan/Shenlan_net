package main

import (
	"bufio"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

)

var user string
var password string
func main(){
	for range time.Tick(time.Millisecond*3000){
		NetWorkStatus := NetWorkStatus()
		if NetWorkStatus!=true {
			err := handleText(filepath.Dir(os.Args[0]) + "/user.txt")
			if err != nil {
				panic(err)
			}
			debuff("logout")
			time.Sleep(1000)
			debuff("login")
			//if "windows" == runtime.GOOS {
			//	fmt.Println("请手动关闭窗口")
			//	select {}
			//}
		}

	}
		select {}

}

func NetWorkStatus() bool {
	cmd := exec.Command("ping","-n","1","baidu.com")
	fmt.Println("NetWorkStatus Start:", time.Now().Unix())
	err := cmd.Run()
	fmt.Println("NetWorkStatus End  :", time.Now().Unix())
	if err != nil {
		fmt.Println(err.Error())
		return false
	} else {
		fmt.Println("Net Status , OK")
	}
	return true
}

func debuff(action string){
	TimeStamp:=fmt.Sprintf("%v",time.Now().Unix())
	url := "http://10.152.250.2/cgi-bin/get_challenge?callback=jsonp" + TimeStamp + "000&username=" + strings.Replace(user,"@","%40",-1)
	res := req(url)
	if strings.Index(res, "\"error\":\"ok\"") == -1 {
		PrintRes(res, action, "failed")
	}else{
		token := strings.Split(strings.Split(res, "lenge\":\"")[1], "\",\"cli")[0]
		ip := strings.Split(strings.Split(res, "_ip\":\"")[1], "\",\"ecode")[0]
		xEncodeStr := "{\"username\":\"" + user + "\",\"ip\":\"" + ip + "\",\"password\":\"" + password + "\",\"acid\":\"1\",\"enc_ver\":\"srun_bx1\"}"
		info := encode(xEncodeStr, token)
		hmd5 := encodeMD5("", token)
		ChkSumStr:=chksum(strings.Join([]string{user, hmd5[5:], "1", ip, "200", "1", info}, token), token)
		info=strings.Replace(strings.Replace(info,"=","%3D",-1),"/","%2F",-1)
		url:=fmt.Sprintf("http://10.152.250.2/cgi-bin/srun_portal?callback=jsonp%v&username=%s&info=%s&chksum=%s&action=%s&ip=%s&password=%s&type=1&ac_id=1&n=200", time.Now().UnixNano()/1000000,user,info,ChkSumStr, action, ip,hmd5)
		if action=="logout"{
			url=fmt.Sprintf("http://10.152.250.2/cgi-bin/srun_portal?callback=jsonp%v&username=%s&info=%s&chksum=%s&action=%s&ip=%s&type=1&ac_id=1&n=200", time.Now().UnixNano()/1000000,user,info,ChkSumStr, action, ip)
		}
		url=strings.Replace(strings.Replace(strings.Replace(strings.Replace(url,"+","%2B",-1),"@","%40",-1),"{","%7B",-1),"}","%7D",-1)
		res = req(url)
		if strings.Index(res, "\"error\":\"ok\"") == -1 {
			PrintRes(res, action, "failed")
		}else{
			PrintRes("IP: " + ip, action, "success")
		}
	}
}

func PrintRes(res string, action string, status string){
	fmt.Println()
	fmt.Println("---------------------------------")
	fmt.Println(res)
	fmt.Println("---------------------------------")
	fmt.Println(action, status)
}
func req(url string) string {
	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err!=nil{
		return err.Error()
	}
	response, _ := client.Do(request)
	if response.StatusCode ==200 {
		body, _ := ioutil.ReadAll(response.Body)
		strs := string(body)
		return strs
	}
	return "failed"
}
func s(a string, b bool) []int {
	c:=len(a)
	var v []int
	for i:= 0;i < c; i = i + 4	{
		if c-i == 1{
			v = append(v, int(a[i]))
		}else if c-i ==2{
			v = append(v, int(a[i])|int(a[i+1])<<8)
		}else if c-i ==3{
			v = append(v, int(a[i])|int(a[i+1])<<8|int(a[i+2])<<16)
		}else{
			v = append(v, int(a[i])|int(a[i+1])<<8|int(a[i+2])<<16|int(a[i+3])<<24)
		}
	}
	if b {
		v = append(v, c)
	}
	return v
}
func l(a []int, b bool) string {
	d:=len(a)
	var bytes []byte
	for i:= 0;i < d; i++	{
		bytes = append(bytes, byte(a[i]&0xff))
		bytes = append(bytes, byte(a[i] >> 8&0xff))
		bytes = append(bytes, byte(a[i] >> 16&0xff))
		bytes = append(bytes, byte(a[i] >> 24&0xff))
	}
	return encodeBase64(bytes)
}
func encode(a string, b string) string{
	v := s(a, true)
	k := s(b, false)
	n := uint(len(v) - 1)
	z := uint(v[n])
	y := uint(v[0])
	c := uint(0x86014019 | 0x183639A0)
	m := uint(0)
	e := uint(0)
	p := uint(0)
	q := uint(6 + 52 / (n + 1))
	d := uint(0)
	for {
		q -= 1
		d = (d + c) & (0x8CE0D9BF | 0x731F2640)
		e = d >> uint(2) & uint(3)
		for p = 0;p < n; p++{
			y = uint(v[p+1])
			m = z >> 5 ^ y << 2
			m += (y>>3 ^ z<<4) ^ (d ^ y)
			m += uint(k[(p&3)^e]) ^ z
			z = (uint(v[p]) + m) & (0xEFB8D130|0x10472ECF)
			v[p] = int(z)
		}
		y = uint(v[0])
		m = z >> 5 ^ y << 2
		m += (y >> 3 ^ z << 4) ^ (d ^ y)
		m += uint(k[(n & 3) ^ e]) ^ z
		v[n] = int((uint(v[n]) + m) & uint(0xBB390742 | 0x44C6F8BD))
		z = uint(v[n])
		if 0 >= q{
			break
		}
	}
	return l(v, false)
}
func encodeBase64(bytes []byte) string{
	const CodeList = "LVoJPiCN2R8G90yg+hmFHuacZ1OWMnrsSTXkYpUq/3dlbfKwv6xztjI7DeBE45QA"
	src := bytes
	encoder := base64.NewEncoding(CodeList)
	out := encoder.EncodeToString(src)
	return "{SRBX1}" + out
}
func encodeMD5(data, key string) string {
	mac := hmac.New(md5.New, []byte(key))
	mac.Write([]byte(data))
	return "{MD5}" + hex.EncodeToString(mac.Sum(nil))
}
func Sha1(data []byte) string {
	sha := sha1.New()
	sha.Write(data)
	return hex.EncodeToString(sha.Sum([]byte(nil)))
}
func chksum(data string, token string) string {
	str:=token+data
	return Sha1([]byte(str))
}
func handleText(fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		log.Printf("Cannot open text file: %s, err: [%v]", fileName, err)
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	user = scanner.Text()
	scanner.Scan()
	password = scanner.Text()
	if err := scanner.Err(); err != nil {
		log.Printf("Cannot scanner text file: %s, err: [%v]", fileName, err)
		return err
	}
	return nil
}
