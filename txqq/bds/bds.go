package bds

import (
    "log"
    "os/exec"
    "os"
    "io"
    "fmt"
    "time"
    "runtime"
    "bufio"
    "strings"
)
//打包cmd  以做到反复重启
type Bds struct {
    Cmd *exec.Cmd
}
const (
//监听bds输出的字符串
    //服务器启动关闭输出
    mtStarting = "Starting Server"
    mtVersion = "Version"
    mtLevelName = "Level Name"
    mtIpv4 = "IPv4 supported"
    mtServerOk = "Server started."
    mtStopOk = "Quit correctly"
    //Ez输出
    mtEzChat = "[CHAT]"
    mtEzMod = "Loading_Loaded"
    //客户端连接断开输出
    mtPc  = "Player connected:"
    mtPd  = "Player disconnected:"
    //备份文件信息输出
    mtHold = ".ldb:"
    //白名单输出
    mtWlAddOk = "Player added to whitelist"
    mtWlRmOk = "Player removed from whitelist"
    mtWlRmErr = "Player not in whitelist"
    //list命令输出
    mtPdList = "players online:"

    //传到main.bdsChat的字符串
    bdsNoStart = "服务器正在运行！无法再次启动"
    startBds = "启动服务器..."
    startOk  = "启动成功"
    stopYes   = "关闭成功"
    wlAddOk = "白名单添加成功!"
    wlRmOk  = "白名单移除成功!"
    wlRmErr = "白名单内不存在此ID!"
    listZero = "有没有人玩自己不知道?"

    CmdStartOk = "Cmd.Start"
    CmdStopOk = "服务器正常关闭"
    CmdStopErr = "服务器非正常关闭! 请通知管理员!"

    ProtZero = "端口为0 如服务器立即非正常关闭 请通知管理员检查端口占用!\n"
)

var (
    list string
    List int
    stdin io.WriteCloser
)

func (bds Bds)Start(bdsChat chan <- string, bdsOk bool) {
    if bdsOk == true {
	//服务器正在运行 无法再次启动
	bdsChat <- bdsNoStart
	return
    }
    if runtime.GOOS == "windows" {
	bds.Cmd = exec.Command("bedrock_server_mod.exe")

    }else{
	bds.Cmd = exec.Command("wine64","bedrock_server_mod.exe")
	bds.Cmd.Env = os.Environ()
	bds.Cmd.Env = append(bds.Cmd.Env, "WINEDLLOVERRIDES=vcruntime140_1,vcruntime140=n;mscoree,mshtml,explorer.exe,winemenubuilder.exe,services.exe,playplug.exe=d")
	bds.Cmd.Env = append(bds.Cmd.Env, "WINEDEBUG=-all")
    }

    Stdout, _ := bds.Cmd.StdoutPipe()
    Stdin, _ := bds.Cmd.StdinPipe()
    stdin = Stdin
    bds.Cmd.Start()
    bdsChat <- CmdStartOk
    bdsPid := bds.Cmd.Process.Pid
    fmt.Printf("Bds_Pid: %s ,%T\n",bdsPid,bdsPid)

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
	    if List != 0 {
		    if List == 1 {
		        bdsChat <- "在线人数: " + list + "\n" + chat
		        List = 0
		        continue
		    }
		    if List == 114514 {
		        bdsChat <- list
		        List = 0
		        continue
		    }
	    }
	    if strings.Contains(chat,mtEzMod[:7]) {
	        ezModLen++
		continue
	    }
	    if strings.Contains(chat,mtEzMod[7+1:]) {
	        ezModSize++
		continue
	    }
	    if strings.Contains(chat,mtHold) {
	        bdsChat <- chat
		continue
	    }
	    if strings.Contains(chat,mtStarting) {
		//startBds = 启动服务器
	        bdsChat <- startBds
		continue
	    }
	    if strings.Contains(chat,mtVersion) {
	        i := strings.Index(chat,mtVersion)
	        i2 := strings.Index(chat,"with")
	        startChat += "BDS " + chat[i:i2] + "\n"
		continue
	    }
	    if strings.Contains(chat,mtLevelName) {
                i := strings.Index(chat,mtLevelName)
	        startChat += "世界名称" + chat[i+10:]
		continue
	    }
	    if strings.Contains(chat,mtIpv4) {
		chat = strings.TrimSpace(chat)
	        i := strings.Index(chat,mtIpv4)
		if chat[i+22:] == "0" {
                    bdsChat <- ProtZero + chat[i:]
		}else{
		    startChat += "端口" + chat[i:i+4] + chat[i+20:] + "\n"
		    startChat += fmt.Sprintf("EzMod加载:%d/%d", ezModSize,ezModLen)
		    bdsChat <- startChat
		}
		continue
	    }
	    if strings.Contains(chat,mtServerOk) {
		//startOk = 启动成功
	        bdsChat <- startOk
		continue
	    }
	    if strings.Contains(chat,mtStopOk) {
		//stopOk = 关闭成功
	        bdsChat <- stopYes
		continue
	    }
	    if strings.Contains(chat,mtEzChat) {
		i := strings.Index(chat,")")
	        bdsChat <- chat[i+2:len(chat)-2]
		continue
	    }
	    if strings.Contains(chat,mtPdList) {
		i := strings.Index(chat,"are")
	        i2 := strings.Index(chat,"players")
		if chat[i+4:i+5] == "0" {
		    list = listZero
		    List = 114514
		}else{
	            list = chat[i+4:i2-1]
		    List = 1
		}
		continue
	    }
	    if strings.Contains(chat,mtPc) {
		    i := strings.Index(chat,mtPc)
		    i2 := strings.Index(chat,",")
	        bdsChat <- fmt.Sprintf("玩家：%s 加入了游戏",chat[i+17:i2])
		    continue
	    }
	    if strings.Contains(chat,mtPd) {
		    i := strings.Index(chat,mtPd)
		    i2 := strings.Index(chat,",")
	        bdsChat <- fmt.Sprintf("玩家：%s 退出了游戏",chat[i+20:i2])
		    continue
	    }
	    if strings.Contains(chat,mtWlAddOk) {
		//wlAdd = 白名单添加成功
		    bdsChat <- wlAddOk
		    continue
	    }
	    if strings.Contains(chat,mtWlRmOk) {
		//wlRm = 白名单移除成功
		bdsChat <- wlRmOk
		continue
	    }
	    if strings.Contains(chat,mtWlRmErr) {
		//wlRmErr = ID不存在白名单内 移除错误
		bdsChat <- wlRmErr
		continue
	    }
        }
    }()

    err1 := bds.Cmd.Wait()
    if err1 != nil {
	    log.Println(CmdStopErr,err1)
	    bdsChat <- CmdStopErr
    }else{
	    bdsChat <- CmdStopOk
    }
    Stdin.Close()
}
func W(chat string,qqShell bool) {
    if qqShell == true {
        if chat == "关闭" {
            stdin.Write([]byte("stop\n"))
	    return
        }
        if chat == "ls" {
            stdin.Write([]byte("list\n"))
	    return
        }
	if strings.HasPrefix(chat,"su ") {
            stdin.Write([]byte(chat[3:] + "\n"))
	    return
	}
	if chat != "" {
	    log.Println("内容:",chat)
	    i := strings.Index(chat,"@")
	    log.Println("下标:",i)
	    ch := fmt.Sprintf("tellraw @a {\"rawtext\":[{\"text\":\"§a[%s] §b%s\"}]}",chat[:i],chat[i+1:])
	    log.Println("转发信息：",ch)
            stdin.Write([]byte(ch + "\n"))
	}
    }else{
        stdin.Write([]byte(chat))
    }
}
const (
    startBackup = "say §6开始备份\n"
    backuping   = "say §6备份中...\n"
    backupStop  = "say §a备份结束\n"

    BackupGo    = "备份开始"
    BackupEnd    = "备份结束"
    BackupErr   = "无法备份！当前已有备份进行中!"

    cmdHold = "save hold\n"
    cmdQuery = "save query\n"
    cmdResume = "save resume\n"
)
func Back(backOk chan bool,bdsChat chan string,backupLock bool) {
    if backupLock != true {
	bdsChat <- BackupGo
    stdin.Write([]byte(startBackup))
    stdin.Write([]byte(cmdHold))
    //5s 应该可以适配99.9%的服务器了吧 :D
    time.Sleep(time.Duration(5)*time.Second)
    stdin.Write([]byte(backuping))
    stdin.Write([]byte(cmdQuery))
    <-backOk
    stdin.Write([]byte(backupStop))
    stdin.Write([]byte(cmdResume))
    <-backOk
        bdsChat <- BackupEnd
    }else{
	bdsChat <- BackupErr
    }
}
