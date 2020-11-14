package main
import (
    "fmt"
    "os"
)

func main() {
    file, err := os.OpenFile("./strbot.log",os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
    if err != nil {
        fmt.Println("打开文件失败")
        return
    }
    fmt.Fprintf(file, "abc123222\n")  // 向file对应文件中写入数据
    file.Close()
}
