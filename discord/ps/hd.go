package ps

import (
	"github.com/shirou/gopsutil/disk"
)

const HdErr = ":x:储存容量达到或超过90%警戒线，备份禁止执行！"

func HdTF() bool {
	d, _ := disk.Usage("/")

	if d.UsedPercent >= 90.0 {
		return false
	} else {
		return true
	}
}
