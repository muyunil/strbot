package tool

import (
    "os"
    "fmt"
)
var (
//    StrBotLog bool
    LogFile *os.File
)
func Log(on bool,logChan chan string){
    if on {
        for chat := range logChan {
            fmt.Fprintf(LogFile, (chat+"\n"))
        }
    }else{
        for {
            <- logChan
        }
    }
}
