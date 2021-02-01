package ps

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"time"
	//    "github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

func Get(netCard int) string {
	v, _ := mem.VirtualMemory()
	c, _ := cpu.Info()
	cc, _ := cpu.Percent(time.Second, true)
	d, _ := disk.Usage("/")
	//    n, _ := host.Info()
	nv, _ := net.IOCounters(true)

	/*    var cpux string
	      var cpuzero string
	      if len(c) > 1 {*/

	var modelName string
	for _, sub_cpu := range c {
		modelName = sub_cpu.ModelName
		//        cores := sub_cpu.Cores
		//        cpux = fmt.Sprintf("        CPU       : %v   %v cores \n", modelname, cores)
	}
	/*
	       }else{
	           sub_cpu := c[0]
	           modelname := sub_cpu.ModelName
	           cores := sub_cpu.Cores

	   //        cpuzero = fmt.Sprintf("CPU: %v  %v cores \n", modelname, cores)
	       }*/
	//    os := fmt.Sprintf("OS: %v(%v)  %v\n", n.Platform, n.PlatformFamily, n.PlatformVersion)

	cpus := "CPU Used:"
	for i := 0; i < len(cc); i++ {
		cpus += fmt.Sprintf(" %vcore:%.0f%%", i, cc[i])
		if i+1 == len(cc) {
			cpus += fmt.Sprintf("\n%s\n", modelName)
		}
	}

	//    mem := fmt.Sprintf("Mem: %v MB  Free: %v MB Used:%v Usage:%.1f%%\n", v.Total/1024/1024, v.Available/1024/1024, v.Used/1024/1024, v.UsedPercent)
	mem := fmt.Sprintf("Mem: %v M Usage:%.0f%%\n", v.Total/1024/1024, v.UsedPercent)

	//    net := fmt.Sprintf("Network: %v / %v M\n", nv[0].BytesRecv/1024/1024, nv[0].BytesSent/1024/1024)
	net := fmt.Sprintf("Network: %v / %v M\n", nv[netCard].BytesSent/1024/1024, nv[netCard].BytesRecv/1024/1024)

	hd := fmt.Sprintf("HD: %v GB Usage:%.0f%%\n", d.Total/1024/1024/1024, d.UsedPercent)

	//    t := "============================\n"
	t := "-------------------------\n"

	return t + cpus + t + mem + t + net + t + hd + t[:len(t)-1]
	//	return os + mem + cpuzero + net  + cpus +  hd
}
