package ps
import (
    "fmt"
    "github.com/shirou/gopsutil/disk"
)

func Hdv() (hd string, f float64) {
    d, _ := disk.Usage("/")

    hd = fmt.Sprintf("HD: %v GB  Free: %v GB Usage:%f%%", d.Total/1024/1024/1024, d.Free/1024/1024/1024, d.UsedPercent)
    fmt.Println(hd)
    fmt.Printf("%T,%v\n",d.UsedPercent,d.UsedPercent)

    return  hd, d.UsedPercent
}
