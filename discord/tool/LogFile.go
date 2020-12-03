package tool

import (
    "os"
    "fmt"
//    ""
)
var (
    StrBotLog bool
    LogFile *os.File
)
//const cs = "@@@"
func Log(logChan chan string){
    if !StrBotLog {
        for {
            <- logChan
        }
    }
    for chat := range logChan {
        fmt.Fprintf(LogFile, (chat+"\n"))
/*	    if strings.HasPrefix(chat,cs) {
            fmt.Fprintf(LogFile, (chat+"\n"))
        }else{
            fmt.Fprintf(LogFile, (chat+"\n"))
        }
*/
    }
}
