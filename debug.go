//中年人的第一款端口扫描器
package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)
var input *string //输入IP或者域名的参数
var file *string //挂载文件的参数

func init(){
	input=flag.String("n","","Please input IP address or domain name")
	file=flag.String("f","","Please attach a file within IP addresses or domain names")
	runtime.GOMAXPROCS(runtime.NumCPU())

}

func portscan(s string,tcpport int,wg *sync.WaitGroup){
	//var PortOpen = make([]int, 0)
	ResolveIP,err:=net.ResolveIPAddr("ip",s)
	if err != nil{
		fmt.Println(s)
		fmt.Println(err)
	}else {
			combine := fmt.Sprintf("%v:%v", ResolveIP, tcpport)
			conn, _ := net.DialTimeout("tcp", combine, 3*time.Second)
			/*if err != nil {
				fmt.Printf("The port %d ......closed\n", tcpport)
			}*/
			if conn != nil {
				fmt.Printf("The port %d .....opened\n", tcpport)
				//PortOpen = append(PortOpen, tcpport)
			}
		}
		wg.Done()
	//return PortOpen
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

func Scan(){

}

func main(){
	flag.Parse()
	wg:=sync.WaitGroup{}
	wg.Add(21)
	if *input == "" && *file== "" {
		fmt.Println("Please follow the below instructions:")
		flag.PrintDefaults()
		os.Exit(1)
	}
	if *file != ""{
		a:=Readfile(*file)
		for _,k := range a{
			for i:=0;i<=65535;i++{
				go portscan(k,i,&wg)
				fmt.Println(k)
				//fmt.Println(b)//10k ports lasts 1min
			}
		}
	}
	//var b []int
	if *input != ""{
		//var a []int
		for i:=8070;i<=8090;i++{
			go portscan(*input,i,&wg)
			//a=portscan(*input,i,&wg)
		}
		fmt.Println(*input)
		//fmt.Println(a)
	}
	wg.Wait()

}

