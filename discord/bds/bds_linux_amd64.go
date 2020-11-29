package bds

import (
    "os"
    "io"
    "fmt"
    "log"
    "time"
    "bufio"
    "os/exec"
    "runtime"
    "strings"
    "syscall"
)
//打包cmd  以做到反复重启
type Bds struct {
    StartPath string
    BdsChat chan <- string
    CrashStart  bool
    CrashStartAll  bool
    LogFile *os.File
    Cmd *exec.Cmd
}
const (
//监听bds输出的字符串
    //服务器启动关闭输出
    mtStarting = "Starting Server"
    mtVersion = "Version"
    mtIpv4 = "IPv4 supported"
    mtServerOk = "Server started."
    mtStop = "Quit correctly"
    //Ez输出
    mtEzChat = "[CHAT]"
    mtEzMod = "Loading_Loaded"
    //客户端连接断开输出
    mtPc  = "Player connected:"
    mtPd  = "Player disconnected:"
    //备份文件信息输出
    MtBackupFileChat = ".ldb:"
    //list命令输出
    mtPdList = "players online:"
    //chat Channel的信息前缀
    ChatChannel = "cchat"
    //传到main.bdsChat的字符串
    startBds = "Starting Bds..."
    startOk  = "Start Ok"
    listZero = "There is no one online at this moment."

    //服务器正在运行！无法再次启动
    CmdNo2Start = ":x:The server is running and cannot be started again"
    CmdStartOk = "Cmd.Start"
    CmdStopErr = ":x:The server shut down by mistake"
	CrashStart = "CrashStart..."

    //端口为0 如服务器立即非正常关闭 请通知管理员检查端口占用
    ProtZero = ":x:The port is zero, if the error is closed, please check the port occupation\n"

)
type sPlayer struct {
    uuid string
    pc int
}
type rcc struct {
    Out bool
    Time time.Time
}
var (
    PlayerNameList string
    ListOk bool
    ListTF bool // 判断是root还是chat频道的命令
    backupQuery bool

    Player = make(map[string]*sPlayer)
    PlayerAdd int
    stdin io.WriteCloser
    Rcc rcc
    StartTime  time.Time
    BdsStartTime time.Time
)

//    fmt.Fprintf(file, "abc123222\n")  // 向file对应文件中写入数据
//    fmt.Fprintf(os.Stdout, "abc123222")  // 向控制台(标准输出设备)输出数据

func (bds Bds)Start(bdsStartLock *bool) {
    if *bdsStartLock == true {
	//服务器正在运行 无法再次启动
	    bds.BdsChat <- CmdNo2Start
        return
    }
    if runtime.GOOS == "windows" {
        bds.Cmd = exec.Command(bds.StartPath)
    }else{
        if bds.StartPath != "bedrock_server_mod.exe" {
            bds.Cmd = exec.Command(bds.StartPath)
        }else{
            bds.Cmd = exec.Command("wine64",bds.StartPath)
            bds.Cmd.Env = os.Environ()
            bds.Cmd.Env = append(bds.Cmd.Env, "WINEDLLOVERRIDES=vcruntime140_1,vcruntime140=n;mscoree,mshtml,explorer.exe,winemenubuilder.exe,services.exe,playplug.exe=d")
            bds.Cmd.Env = append(bds.Cmd.Env, "WINEDEBUG=-all")
            bds.Cmd.SysProcAttr = &syscall.SysProcAttr{
                Setpgid: true,
            }
        }
    }

    Stdout, _ := bds.Cmd.StdoutPipe()
    Stdin, _ := bds.Cmd.StdinPipe()
    stdin = Stdin
    bds.Cmd.Start()
    BdsStartTime = time.Now()
    *bdsStartLock = true
    //init player...
    Player = make(map[string]*sPlayer)
    PlayerAdd = 0

    go func(){
        var (
	        ezModLen int
	        ezModSize  int
	        startChat string
	    )
        reader := bufio.NewReader(Stdout)
        for {
            chat, err := reader.ReadString('\n')
	        if err != nil {
		        return
	        }
            fmt.Println(chat)
            if bds.LogFile != nil{
                // 向file中写入数据
                fmt.Fprintf(bds.LogFile, ("BdsOut>Bds>"+chat+"\n"))
            }
            if Rcc.Out {
                if time.Now().Before(Rcc.Time.Add(time.Second)) {
                    bds.BdsChat <- chat
                }else{
                    Rcc.Out = false
                }
            }
	        if ListOk && ListTF {
		        if PlayerNameList == "" {
		            bds.BdsChat <- ChatChannel + listZero
		            ListOk = false
		        }else{
		            bds.BdsChat <- ChatChannel + "Player list: " + PlayerNameList + "\n" + chat
		            ListOk = false
                }
                ListTF = false
                continue
	        }
	        if strings.Contains(chat,mtEzMod[:7]) {
	            ezModLen++
		        continue
	        }
	        if strings.Contains(chat,mtEzMod[7+1:]) {
	            ezModSize++
		        continue
	        }
	        if strings.Contains(chat,MtBackupFileChat) {
                if backupQuery {
                    bds.BdsChat <- chat
                    backupQuery = false
                }
		        continue
	        }
	        if strings.Contains(chat,mtStarting) {
		        //startBds = 启动服务器
	            bds.BdsChat <- startBds
		        continue
	        }
	        if strings.Contains(chat,mtStop) {
                bds.BdsChat <- mtStop
		        continue
            }
	        if strings.Contains(chat,mtVersion) {
	            i := strings.Index(chat,mtVersion)
	            i2 := strings.Index(chat,"with")
	            startChat += "BDS " + chat[i:i2] + "\n"
		        continue
	        }
	        if strings.Contains(chat,mtIpv4) {
		        chat = strings.TrimSpace(chat)
	            i := strings.Index(chat,mtIpv4)
		        if chat[i+22:] == "0" {
                    bds.BdsChat <- ProtZero + chat[i:]
		        }else{
		            startChat += "Port" + chat[i:i+4] + chat[i+20:] + "\n"
                    startChat += fmt.Sprintf("EzMod cap:%d len:%d", ezModSize,ezModLen)
		            bds.BdsChat <- startChat
		        }
		        continue
	        }
	        if strings.Contains(chat,mtServerOk) {
		        //startOk = 启动成功
	            bds.BdsChat <- startOk
		        continue
	        }
	        if strings.Contains(chat,mtEzChat) {
		        i := strings.Index(chat,")")
bds.BdsChat <- fmt.Sprintf("%s:arrow_right:`%s`",ChatChannel,chat[i+2:len(chat)-2])
		        continue
	        }
	        if strings.Contains(chat,mtPdList) {
		        i := strings.Index(chat,"are")
	            i2 := strings.Index(chat,"players")
		        if chat[i+4:i+5] == "0" {
		            PlayerNameList = ""
		            ListOk = true
		        }else{
	                PlayerNameList = chat[i+4:i2-1]
		            ListOk = true
		        }
		        continue
	        }
	        if strings.Contains(chat,mtPc) {
		        i := strings.Index(chat,mtPc)
		        i2 := strings.Index(chat,",")
                playerName :=  fmt.Sprintf(chat[i+18:i2])
                pcPlayer :=  fmt.Sprintf("Player %s joined the server",chat[i+18:i2])
	            bds.BdsChat <- ChatChannel + pcPlayer
                fmt.Println(playerName,len(playerName))
                if _, ok := Player[playerName]; !ok {
                    //No
                    PlayerAdd++
                    Player[playerName] = &sPlayer{}
                }else{
                    //Yes
//                    Player[playerName].pc++
                }
		        continue
	        }
	        if strings.Contains(chat,mtPd) {
		        i := strings.Index(chat,mtPd)
		        i2 := strings.Index(chat,",")
                playerName := fmt.Sprintf(chat[i+21:i2])
                pdPlayer := fmt.Sprintf("Player %s logged out of the server",chat[i+21:i2])
	            bds.BdsChat <- ChatChannel + pdPlayer
                fmt.Println(playerName,len(playerName))
                delete(Player,playerName)
		        continue
	        }
            //add code:
//	        bds.BdsChat <- chat
        }
    }()

    err1 := bds.Cmd.Wait()
    if err1 != nil {
        *bdsStartLock = false
	    log.Println(CmdStopErr,err1)
	    bds.BdsChat <- CmdStopErr
        if bds.CrashStart {
	        bds.BdsChat <- CrashStart
            mc := Bds{bds.StartPath,bds.BdsChat,bds.CrashStart,bds.CrashStartAll,bds.LogFile,nil}
            go mc.Start(bdsStartLock)
        }
    }else{
        *bdsStartLock = false
        if bds.CrashStartAll {
	        bds.BdsChat <- CrashStart
            mc := Bds{bds.StartPath,bds.BdsChat,bds.CrashStart,bds.CrashStartAll,bds.LogFile,nil}
            go mc.Start(bdsStartLock)
        }
    }
    Stdin.Close()
}

func RootW(bdsShell bool, chat string) {
    if  bdsShell == true {
	    if strings.HasPrefix(chat,"/") {
            stdin.Write([]byte(chat[1:] + "\n"))
            goto End
        }
	    if strings.HasPrefix(chat,"wl ") {
            if len(chat) < 6 {
                return
            }
            if strings.HasPrefix(chat,"wl + ") {
                stdin.Write([]byte(fmt.Sprintf("whitelist add \"%s\"\n",chat[5:])))
                goto End
            }
            if strings.HasPrefix(chat,"wl - ") {
                stdin.Write([]byte(fmt.Sprintf("whitelist remove \"%s\"\n",chat[5:])))
            }
	    }

        End:
        if !strings.HasPrefix(chat,"/stop") {
            //暂时打开向root频道的log输出
            Rcc = rcc{true,time.Now()}
        }
    }else{
        stdin.Write([]byte(chat))
    }
}
func ChatW(chat string) {
    if chat == "ls" {
        stdin.Write([]byte("list\n"))
        ListTF = true
        return
    }
	if chat != "" {
	    log.Println("内容:",chat)
	    i := strings.Index(chat,"@")
        if i < 1 {
	        log.Println("转发错误！内容:",chat)
            return
        }
	    log.Println("下标(Index):",i)
        ch := fmt.Sprintf("tellraw @a {\"rawtext\":[{\"text\":\"§a[%s] §b%s\"}]}",chat[:i],chat[i+1:])
	    log.Println("转发信息：",ch)
        stdin.Write([]byte(ch + "\n"))
    }
}
const (
    startBackup = "say §6Startin Backup\n"
    backuping   = "say §6Backuping...\n"
    backupStop  = "say §aEnd of backup\n"


    BackupGo    = "BackupGo..."
    BackupEnd    = "BackupEnd!"
    //无法备份！当前已有备份进行中!
    BackupErr   = ":x:Could not backup!  Only one backup process can exist at the same time"

    cmdHold = "save hold\n"
    cmdQuery = "save query\n"
    cmdResume = "save resume\n"
)
func Back(backChat chan string,bdsChat chan string,backupLock *bool) {
    if *backupLock != true {
        *backupLock = true
	    bdsChat <- BackupGo
        stdin.Write([]byte(startBackup))
        stdin.Write([]byte(cmdHold))
        //等待3s
        time.Sleep(time.Duration(3)*time.Second)
        backupQuery = true
        stdin.Write([]byte(backuping))
        stdin.Write([]byte(cmdQuery))
        if str := <-backChat;str != "" {
            bdsChat <- str
        }
        stdin.Write([]byte(backupStop))
        stdin.Write([]byte(cmdResume))
        if str := <-backChat;str != "" {
            bdsChat <- str
        }
        *backupLock = false
        bdsChat <- BackupEnd
    }else{
	    bdsChat <- BackupErr
    }
}