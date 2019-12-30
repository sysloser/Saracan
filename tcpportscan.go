//中年人的第一款端口扫描器
//现在是单线程扫描,支持单个IP或者域名，以及挂载一个IP或域名的文件
//待改进为多线程扫描。
/*
Author: Phil
Version: 1.0
Date: 30 Dec, 2019 @ Nagoya
*/
package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)
var input *string //输入IP或者域名的参数
var file *string //挂载文件的参数

func init(){
	input=flag.String("n","","Please input IP address or domain name")
	file=flag.String("f","","Please attach a file within IP addresses or domain names")

}

func portscan(s string) []int {
	var PortOpen = make([]int, 0)
	ResolveIP,err:=net.ResolveIPAddr("ip",s)
	if err != nil{
		fmt.Println(s)
		fmt.Println(err)
	}else {
		var tcpport int
		for tcpport = 1; tcpport <= 100; tcpport++ {
			combine := fmt.Sprintf("%v:%v", ResolveIP, tcpport)
			conn, err := net.DialTimeout("tcp", combine, 1*time.Second)
			if err != nil {
				fmt.Printf("The port %d ......closed\n", tcpport)
			}
			if conn != nil {
				fmt.Printf("The port %d .....opened\n", tcpport)
				PortOpen = append(PortOpen, tcpport)
			}
		}
	}
	return PortOpen
}

//读取文件,每一行是一个数组元素
func Readfile(s string)[]string{
	file, err := os.Open(s)
	if err != nil{
		fmt.Println("file is broken",err)
	}
	defer file.Close()
	data := make([]string,0)
	reader := bufio.NewReader(file)
	for{
		linestr, err := reader.ReadString('\n')
		if err != nil{
			break
		}
		linestr = strings.TrimSpace(linestr)
		if linestr == ""{
			continue
		}
		data=append(data,linestr)
	}
	return data
}

//扫描文件中的域名端口开放
func main(){
	flag.Parse()
	if *input == "" && *file== "" {
		fmt.Println("Please follow the below instructions:")
		flag.PrintDefaults()
		os.Exit(1)
	}
	if *file != ""{
		a:=Readfile(*file)
		for _,k := range a{
				b:=portscan(k)
				fmt.Println(k)
				fmt.Println(b)
			}
		}
	if *input != ""{
		b:=portscan(*input)
		fmt.Println(*input)
		fmt.Println(b)
	}
}
