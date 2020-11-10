package backup

import (
    "shibot/bot"

    "os"
    "fmt"
    "log"
    "time"
    "strings"
    "strconv"
    "github.com/otiai10/copy"
)
const backErr = "截断文件错误! 请通知管理员查看后台日志"

func BackUp(str string,backOk chan bool) {
    _, month, day := time.Now().Date()
    timeStr:=time.Now().Format("2006-01-02_15-04-05")

    backDir := fmt.Sprintf("./backup/%d/%d/worlds-%s/",int(month),day,timeStr)
    fmt.Println(backDir)
    err := os.MkdirAll(backDir, os.ModePerm);
    if err != nil {
        log.Fatal(err);
    }

    err2 := copy.Copy("./worlds/",backDir)
    if err2 != nil {
	fmt.Println(err2)
    }

    //传递复制成功信息
    backOk <- true

    sp := make([]string,0,10)

    tt := strings.SplitN(str,",",-1)

    for _, v := range tt {
	if v[0] == ' ' {
	    sp = append(sp,v[1:])
	}else{
	    sp = append(sp,v)
	}
    }
//	fmt.Println(sp)

	backOk <- tc(sp,backDir)
}

func tc(path []string,backDir string) bool {

    for _, v := range path {
	i := strings.Index(v,":")

	ldb := v[:i]
	si := v[i+1:]
//	size, _ := strconv.Atoi(si)
	size, _ := strconv.ParseInt(si,10,64)
	fmt.Printf("%T %s\n",ldb,ldb)
	fmt.Printf("%T %d\n",size,size)

	fmt.Println(backDir+ldb)
	err := os.Truncate(backDir+ldb, size)
	if err != nil {
	    bot.Printqq(backErr)
            log.Fatal(err)
        }
    }
    return true
    // 裁剪一个文件到size个字节。
    // 如果文件本来就少于size个字节，则文件中原始内容得以保留，剩余的字节以null字节填充。
    // 如果文件本来超过size个字节，则超过的字节会被抛弃。
    // 这样我们总是得到精确的size个字节的文件。
    // 传入0则会清空文件。

}
