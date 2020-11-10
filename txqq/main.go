package main

import (
    "shibot/bot"
    "shibot/bds"
    "shibot/ps"
    "shibot/backup"
    "os"
    "log"
    "sync"
    "time"
    "bufio"
    "strings"
    "github.com/robfig/cron"
)

const (
    bdsNo = "服务器未运行"
    hdErr    = "储存空间不足!无法备份!"
)
var (
    qqChat   = make(chan string,10)
    bdsChat  = make(chan string,20)
    rdChat   = make(chan string,1)
    backOk   = make(chan bool)
    bdsOk  bool
    backupLock bool
)
func main() {

    var wg sync.WaitGroup
    wg.Add(2)
    go bot.Start(qqChat)

    go func() {
    defer wg.Done()
	for chat := range qqChat {
	    if strings.Contains(chat,"Cron") {
            mCron(chat)
		    continue
        }
        if chat == bot.StartBds {
	        mc := bds.Bds{}
	        go mc.Start(bdsChat,bdsOk)
		    continue
        }
	    if chat == bot.Ps {
		    bot.Printqq(ps.Get())
		    continue
	    }
	    if strings.HasPrefix(chat,bot.LsBd) {
		    bot.Printqq(backup.LsBD(chat))
		    continue
	    }
	    if strings.HasPrefix(chat,bot.Rd) {
		    go mRd(chat,rdChat)
		    continue
	    }
	    if chat == bot.Backup {
            mBackup()
		    continue
	    }
	    if bdsOk == true {
	        bds.W(chat,true)
	    }
    }
    }()

    go func(){
        defer wg.Done()
        for chat := range bdsChat {
	    if strings.Contains(chat,bds.BackupGo) {
		backupLock = true
		continue
	    }
	    if strings.Contains(chat,bds.BackupEnd) {
		backupLock = false
		bot.Printqq(bds.BackupEnd)
		continue
	    }
	    if strings.Contains(chat,bds.CmdStartOk) {
		bdsOk = true
		continue
	    }
	    if strings.Contains(chat,bds.CmdStopErr) {
		bdsOk = false
		bot.Printqq(bds.CmdStopErr)
		continue
	    }
	    if strings.Contains(chat,bds.CmdStopOk) {
		bdsOk = false
		bot.Printqq(bds.CmdStopOk)
		continue
	    }
	    if strings.Contains(chat,".ldb:") {
		backup.BackUp(chat,backOk)
		continue
	    }
	    bot.Printqq(chat)
        }
    }()

    go func () {
	inputReader := bufio.NewReader(os.Stdin)
	    for {
		if bdsOk == true {
		    input, _ := inputReader.ReadString('\n')
		    bds.W(input,false)
		    continue
	        }
		if bdsOk == false {
		        input, _ := inputReader.ReadString('\n')
		        if strings.Contains(input,"start") {
			    mc := bds.Bds{}
	                    go mc.Start(bdsChat,bdsOk)
			    time.Sleep(time.Duration(3) * time.Second)
			    continue
			}else{
			    log.Println("输入的命令：",input)
		            log.Println("服务器未启动控制台无法使用")
		            log.Println("可输入start启动")
			    continue
			}
		}
	}
    }()

    wg.Wait()
}

func mCron(chat string){
	i := strings.Index(chat,"_")
    setCron := chat[i+1:]
    log.Printf("setCron_%v\n",setCron)
    c := cron.New()
    c.AddFunc(setCron, func() {
        if bdsOk == true {
	        if _, f := ps.Hdv();f < 90.00 {
	            bds.Back(backOk,bdsChat,backupLock)
			    }else{
			        bot.Printqq(hdErr)
			    }
		    }else{
			    log.Println(bdsNo)
			    bot.Printqq(bdsNo)
		    }
		})
	go c.Start()
}
func mBackup() {
    if bdsOk == true {
        if _, f := ps.Hdv();f < 90.00 {
           go bds.Back(backOk,bdsChat,backupLock)
	    }else{
            bot.Printqq(hdErr)
	    }
	}else{
        log.Println(bdsNo)
        bot.Printqq(bdsNo)
    }
}

const RdErr = "正在运行 无法回滚备份!"

func mRd(rDir string,rdChat chan string) {
    if bdsOk == false {
        backup.Rd(rDir,rdChat)
        bot.Printqq(<-rdChat)
	}else{
        log.Println(RdErr)
        bot.Printqq(RdErr)
    }
}
