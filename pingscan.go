/*
1.输入被检查IP
 1.1 单个IP参数 -i
 1.2 一段IP参数 -r
 1.3 一个子网参数 -n
 1.4 一个文件参数 -f
 1.5 一个域名参数 -n
2.对输入的IP进行ping探测
3.探测成功，返回"IP可达"
4.探测失败，返回"IP不可达"
 */
package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"github.com/henrylee2cn/pholcus/common/ping"
	"log"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var ipsolo *string //挂载一个IP
var rangeip *string  //挂载一段IP
var netip *string //挂载一个IP子网
var file *string //挂载一个文件
var name *string //挂载一个域名

func init() {
	ipsolo = flag.String("i", "", "输入单个IP地址，如：192.168.1.1")
	rangeip = flag.String("r", "", "输入一段IP地址，如：192.168.1.1-192.168.1.2")
	netip = flag.String("n", "", "输入一个IP子网，如：192.168.1.0/24")
	file = flag.String("f","","附上一个包含IP地址的文本文件")
	name = flag.String("d","","输入一个域名")
}

//逐行读取文件
func Readfile(s string) []string {
	file, err := os.Open(s)
	if err != nil {
		fmt.Println("file is broken", err)
	}
	defer file.Close()
	data := make([]string, 0)
	reader := bufio.NewReader(file)
	for {
		linestr, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		linestr = strings.TrimSpace(linestr)
		if linestr == "" {
			continue
		}
		data = append(data, linestr)
	}
	return data
}
//ip到数字
func ip2Long(ip string) uint32 {
	var long uint32
	binary.Read(bytes.NewBuffer(net.ParseIP(ip).To4()), binary.BigEndian, &long)
	return long
}
//数字到IP
func backtoIP4(ipInt int64) string {
	// need to do two bit shifting and “0xff” masking
	b0 := strconv.FormatInt((ipInt>>24)&0xff, 10)
	b1 := strconv.FormatInt((ipInt>>16)&0xff, 10)
	b2 := strconv.FormatInt((ipInt>>8)&0xff, 10)
	b3 := strconv.FormatInt((ipInt & 0xff), 10)
	return b0 + "." + b1 + "." + b2 + "." + b3
}

//拆分ip段，形成IP数组
func Splitiprange(s string)[]string{
	var iprange []string
	pattern,_:=regexp.MatchString("-",s)
	if !pattern {
		fmt.Println("请按照格式输入一个IP段")
		os.Exit(1)
	}
	a:=strings.Split(s,"-")
	ip1:=ip2Long(strings.TrimSpace(a[0]))
	ip2:=ip2Long(strings.TrimSpace(a[1]))
	if net.ParseIP(a[0])==nil || net.ParseIP(a[1])==nil{
		fmt.Println("IP地址格式错误")
		os.Exit(1)
	}
	if ip1>ip2{
		fmt.Println("IP输反了")
		os.Exit(1)
	}
	for i:=ip1;i<=ip2;i++{
		i:=int64(i)
		iprange=append(iprange,backtoIP4(i))
	}
	return iprange
}

//拆分ip子网
func Splitipnet(s string)[]string{
	j:=func(char rune)bool{return !unicode.IsNumber(char)}
	numarray:=strings.FieldsFunc(s,j)
	return numarray
}

//形成IP子网的IP数组，支持超网
func IPnettoip(ipwithmask string)[]string{
	var ipnetrange []string
	a,_,err:=net.ParseCIDR(ipwithmask)
	if err != nil{
		log.Fatal(err)
		fmt.Println("IP子网输入格式错误")
		os.Exit(1)
	}
	startip:=ip2Long(a.String())
	mask,_:=strconv.Atoi(Splitipnet(ipwithmask)[4])
	yanma:=ip2Long("255.255.255.255")
	lastip:=startip|(yanma>>(mask))
	for i:=startip;i<=lastip;i++{
		i:=int64(i)
		ipnetrange=append(ipnetrange,backtoIP4(i))
	}
	return ipnetrange
}
//pingcheck
func pingcheck(ips []string){
	for _,k:=range ips{
		result,err,_:=ping.Ping(k,2)
		if err != nil{
			fmt.Printf("%s 不可达 \n",k)
		}
		if result {
			fmt.Printf("%s 可达 \n",k)
		}
	}
}



func main(){
	flag.Parse()
	if *ipsolo==""&&*rangeip==""&&*netip==""&&*file==""&&*name==""{
		fmt.Println("Please follow the below instructions:")
		flag.PrintDefaults()
		os.Exit(1)
	}
	if *ipsolo!=""{
		if net.ParseIP(*ipsolo)==nil{
			fmt.Println("IP地址输入错误")
			os.Exit(1)
		}
		result,err,_:=ping.Ping(*ipsolo,2)
		if err != nil{
			fmt.Printf("%v 不可达 \n",*ipsolo)
		}
		if result {
			fmt.Printf("%v 可达 \n",*ipsolo)
		}
	}
	if *rangeip!=""{
		a:=Splitiprange(*rangeip)
		pingcheck(a)
	}
	if *netip!=""{
		a:=IPnettoip(*netip)
		pingcheck(a)
	}
	if *file!=""{
		a:=Readfile(*file)
		pingcheck(a)
	}
	if *name!=""{
		result,err,_:=ping.Ping(*name,2)
		if err != nil{
			fmt.Printf("%v 不可达 \n",*ipsolo)
		}
		if result {
			fmt.Printf("%v 可达 \n",*ipsolo)
		}
	}
}
