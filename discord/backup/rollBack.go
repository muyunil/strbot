package backup

import (
    "os"
    "fmt"
    "strings"
    "github.com/otiai10/copy"
)

const (
    RdDirErr = ":x: 回滚存档路径错误!\n检查参数路径是否存在及格式正确!\n格式参考:\n<rd worlds-2020-12-10_16-44-44>\n使用lsbd查看可回滚存档列表"
    RdOk  = "rollBack Ok!\n回滚完成！"
    RdCpErr = ":x: rollBack Err!\ncopyErr！请管理员查看日志！"
)

/*func main(){
//    fmt.Println(Rd("rb worlds-2020-10-08_20-14-45"))
    fmt.Println(Rd("rb worlds-2020-10-98_20-14-45"))
}*/

func Rd(dir string,rdChat chan <- string) {
    var rDir string
    tmp := strings.Split(dir, " ")
    if len(tmp[1]) < 21 {
        rdChat <- RdDirErr + "1"
        return
    }
	if strings.HasPrefix(tmp[1],"worlds-") {
        tmp2 := strings.Split(tmp[1],"_")
        if len(tmp2) < 2 {
            rdChat <- RdDirErr + "2"
            return
        }
        tmp3 := strings.Split(tmp2[0],"-")
        if len(tmp3) < 4 {
            rdChat <- RdDirErr + "3"
            return
        }
        if tmp3[3][0] == '0' {
            rDir += fmt.Sprintf("./backup/%s/%s/%s/",tmp3[2],string(tmp3[3][1]),tmp[1])
        }else{
            rDir += fmt.Sprintf("./backup/%s/%s/%s/",tmp3[2],tmp3[3],tmp[1])
        }
        fmt.Println(rDir)
    }else{
        rdChat <- RdDirErr + "4"
        return
    }

    if Exist(rDir) == true {
        err := copy.Copy(rDir,"./worlds/")
        if err != nil {
            rdChat <-RdCpErr
	        fmt.Println(err)
        }
    }else{
        rdChat <- RdDirErr + "5"
        return
    }
    rdChat <- RdOk
}
func Exist(filename string) bool {
     _, err := os.Stat(filename)
    return err == nil || os.IsExist(err)
}
