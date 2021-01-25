package ps
import (
    "fmt"
    "github.com/shirou/gopsutil/disk"
    "github.com/shirou/gopsutil/mem"
    "github.com/shirou/gopsutil/net"
    "github.com/shirou/gopsutil/load"
)

func TopicTop(netCard int) string {
    v, _ := mem.VirtualMemory()
    d, _ := disk.Usage("/")
    nv, _ := net.IOCounters(true)
    info, _ := load.Avg()

    return fmt.Sprintf("Load:[%.0f-%.0f-%.0f]|Mem:%.0f%%|Net:%v/%vM|HD:%.0f%%",info.Load1,info.Load5,info.Load15,v.UsedPercent,nv[netCard].BytesSent/1024/1024,nv[netCard].BytesRecv/1024/1024,d.UsedPercent)
}
