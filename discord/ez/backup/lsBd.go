package backup

import (
	"fmt"
    "time"
	"path/filepath"
    "strings"
)

const (
    n = "\n"
    backupName = "worlds-"
    BdErr = "自定义查询路径非法!\n正确写法 <月/日> ，如: 12/10\n(无法查询无备份的日期时间!)"
)

var (
    bd string
)

func LsBD(lsdir string) string {
    //bd = Backup directory
    if lsdir == "lsbd" {
        _, month, day := time.Now().Date()
        bd = fmt.Sprintf("./backup/%d/%d/",int(month),day)
    }else{
        test := strings.Split(lsdir, " ")
        lsdir = test[1]
        bd = fmt.Sprintf("./backup/%s/",lsdir)
    }

    fmt.Println(bd)

    dir,err2 := filepath.Glob(filepath.Join(bd,"*"))
	if err2 != nil {
		fmt.Println(err2)
	}
    if len(dir) == 0 {
        return BdErr
    }

    fmt.Println(dir)

    backUpList := "Backup List\n"
	for i := range dir {
	    i2 := strings.Index(dir[i],backupName)
        if i2 == -1 {
            return BdErr
        }
        if i < len(dir) -1 {
            backUpList += dir[i][i2:] + n
        }else{
            backUpList += dir[i][i2:]
        }
	}
    return backUpList
}
