package utils

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v5/table"
	serial "github.com/tarm/goserial"
	"io"
	"log"
	"strings"
	"time"
)

const MAXRWLEN = 8000

func InitSerial() io.ReadWriteCloser {
	ComID := fmt.Sprintf("%v",GetConfig("serial.comid"))
	Baud := GetConfig("serial.buad").(int)
	ReadTimeout := 30 * time.Microsecond		/*毫秒*/

	SerialCfg := &serial.Config{Name: ComID, Baud: Baud, ReadTimeout: ReadTimeout}
	iorwc, err := serial.OpenPort(SerialCfg)
	if err != nil {
		fmt.Println("端口被占用或端口不存在！")
		log.Panic(err)
		return nil
	}
	return iorwc
}

func List()  {
	config := GetConfig("hospitals").([]interface{})		//转成[]interface{}格式
	t := table.NewWriter()
	t.AppendHeader(table.Row{"#", "名称"})

	for i,_ := range config {
		title := GetConfig("hospitals."+fmt.Sprintf("%v",i)+".title")
		t.AppendRow([]interface{}{i, title})
	}
	fmt.Println(t.Render())
}

func Now()  {
	iorwc := InitSerial()
	defer iorwc.Close()
	buffer := make([]byte, MAXRWLEN)
	//发命令之前清空缓冲区
	_, err := iorwc.Read(buffer)
	if err != nil {
		fmt.Println(err)
		return
	}

	//发命令数据类型为[]byte
	num, err := iorwc.Write([]byte("AT+MAC?\r\n"))
	if err != nil {
		fmt.Println(err)
		return
	}
	var tmpstr string = ""
	for i := 0; i < 3000; i++ {
		num, err = iorwc.Read(buffer)
		if num > 0 {
			tmpstr += fmt.Sprintf("%s", string(buffer[:num]))
		}
		//查找读到信息的结尾标志
		if strings.LastIndex(tmpstr, "\r\nOK\r\n") > 0 {
			break
		}
	}
	fmt.Println(tmpstr)

	var HospIndex int = -1
	var HospTitle string = ""
	// Mac地址存在的医院
	config := GetConfig("hospitals").([]interface{})		//转成[]interface{}格式
	for i,_ := range config {
		cmd := GetConfig("hospitals."+fmt.Sprintf("%v",i)+".cmd").([]interface{})
		MacAddress := cmd[2].(string)
		MacAddress = strings.Split(MacAddress,"=")[1]
		if strings.Contains(tmpstr,MacAddress) {
			HospIndex = i
			HospTitle = GetConfig("hospitals."+fmt.Sprintf("%v",i)+".title").(string)
			break
		}
	}
	if HospIndex < 0 {
		fmt.Println(tmpstr)
		fmt.Println("当前Mac地址没有匹配到任一医院")
		return
	}
	fmt.Println("hospital current:")
	t := table.NewWriter()
	t.AppendHeader(table.Row{"#", "名称"})
	t.AppendRow([]interface{}{HospIndex, HospTitle})
	fmt.Println(t.Render())
}

func Set(index int)  {
	iorwc := InitSerial()
	defer iorwc.Close()
	buffer := make([]byte, MAXRWLEN)

	//发命令之前清空缓冲区
	num, err := iorwc.Read(buffer)
	if err != nil {
		fmt.Println(err)
		return
	}

	cmd := GetConfig("hospitals."+fmt.Sprintf("%v",index)+".cmd").([]interface{})
	for i2,_ := range cmd {
		//按顺序执行命令
		num, err = iorwc.Write([]byte(cmd[i2].(string) + "\r\n"))
		if err != nil {
			log.Panic(err)
		}
		var tmpstr string = ""
		for i := 0; i < 3000; i++ {
			num, err = iorwc.Read(buffer)
			if num > 0 {
				tmpstr += fmt.Sprintf("%s", string(buffer[:num]))
			}
			//查找读到信息的结尾标志
			if strings.LastIndex(tmpstr, "\r\nOK\r\n") > 0 {
				break
			}
		}
		fmt.Println(tmpstr)
	}
	fmt.Println("设置成功，当前医院为：")
	t := table.NewWriter()
	t.AppendHeader(table.Row{"#", "名称"})
	t.AppendRow([]interface{}{index, GetConfig("hospitals."+fmt.Sprintf("%v",index)+".title").(string)})
	fmt.Println(t.Render())
}