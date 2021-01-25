package main

import (
    "strbot/ps"
    "strbot/bds"
    "strbot/tool"
    "strbot/config"
    "strbot/backup"
    "os"
    "log"
    "fmt"
    "time"
    "bufio"
    "strings"
    "syscall"
    "os/signal"
    "github.com/robfig/cron"
    "github.com/bwmarrin/discordgo"
)
const (
    Ls = "ls"
    Ps = "ps"
    Rd = "rd "
    Wl = "wl "
    LsBd = "lsbd"
    StartBds = "start"
    Backup = "backup"
    Help = "help"
    RetHelp = `Root Cmd:
    <ps> :View server status
    <start> :Start Bds
    <backup> :Backup Worlds (Bds Not available when not running)
    <lsbd> :View today's backup archive list
        You can view the backup archives of other dates with parameters, example：
        <lsbd 10/10> :View the backup archive of October 10th
    <rd name-time> :<lsbd> Get partial list, Choose a backup to roll back, example：
    <rd worlds-2020-10-10_10-10-10>
    ^ Roll back the worlds to October 10, 2020 1 ten ten ten seconds
        ps:Bds Unable to roll back while running
    <wl> whitelist Add/remove, example：
        <wl + tes  tID>  <wl - tes  tID> :Comes with "" by default
    </cmd> : example:
    	</say hello>, </stop>.
	`

    //服务器未运行
    bdsNo = "Bds is not running"
)

var (
    initErr  bool
    backupLock    bool
    bdsStartLock  bool

    S  *discordgo.Session
    Conf *config.Config

    backChan = make(chan string)
    rdChat   = make(chan string,1)
    logChan  = make(chan string,16)
    bdsChat  = make(chan string,32)
)

func init() {
    fmt.Print("strbot.yaml... ")
    var Nil bool
    Conf, Nil = config.YamlConfig()
    if Nil {
        fmt.Println(config.ConfNil)
        initErr = true
        return
    }else{
        bds.LogChan = logChan
        bds.BdsChat = bdsChat
        bds.BackChan = backChan
        bds.StartPath = Conf.Bds.StartPath
        bds.CrashStart = Conf.Bds.CrashStart
        bds.ZeroCrashStart = Conf.Bds.ZeroCrashStart
    }
    fmt.Println("ok")
    fmt.Print("init DiscordBot... ")
	dg, err := discordgo.New("Bot " + Conf.Bot.Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
        initErr = true
		return
	}
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)
    err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
        initErr = true
		return
	}
    if Conf.LogFile {
        tool.LogFile, err = os.OpenFile("./strbot.log",os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
        if err != nil {
            fmt.Println("open strbot.log err",err)
            initErr = true
            return
        }
        go tool.Log(true,logChan)
    }else{
        go tool.Log(false,logChan)
    }

    S = dg
    fmt.Println("ok")
}
func main() {
    if initErr {
        return
    }

    if Conf.Bot.ChannelTopic.Enabled {
        //性能监视，一分钟一次
        go RootChannelPsutil()
        //在线人数，累积登录
        go ChatChannelGameInfo()
    }

    if Conf.Bds.CronBackup != "* * * * * *" {
        //定时热备份
	    mCron(Conf.Bds.CronBackup)
    }

    //有信息时触发函数
	S.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
	    if m.Author.ID == s.State.User.ID {
	        return
        }

/*
        fmt.Println(">m.Id",m.Author.ID)
        fmt.Println(">m.ChannelID",m.ChannelID)
        fmt.Println(">m.GuildID",m.GuildID)
        fmt.Println(">S.StateUserID",s.State.User.ID)
*/
	    if m.GuildID != Conf.Bot.GuildID {
		    return
	    }

	    fmt.Printf("[%s] %s\n", m.Author, m.Content)
        logChan <- "disChat>>"+m.Content

	    if m.ChannelID == Conf.Bot.RootChannelID {
	        if m.Content == Ps {
		        MessageSend(true,ps.Get(Conf.Bot.ChannelTopic.ChatTopicNetCard))
                return
	        }
	        if m.Content == StartBds {
                mc := bds.Bds{}
	            go mc.Start(&bdsStartLock)
                return
	        }
	        if m.Content == Backup {
                mBackup()
                return
	        }
	        if strings.HasPrefix(m.Content, LsBd) {
                if len(m.Content) < 5 {
	                MessageSend(true,backup.LsBD(LsBd))
		            return
                }else{
		            i := strings.Index(m.Content," ")
                    if i == 4 {
	                    MessageSend(true,backup.LsBD(m.Content))
                        return
                    }
		            return
                }
            }
	        if strings.HasPrefix(m.Content, Rd) {
                if len(m.Content) > 23 {
		            go mRd(m.Content,rdChat)
                }
                return
            }
	        if m.Content == Help {
		        MessageSend(true,RetHelp)
                return
	        }
            if bdsStartLock {
                bds.RootW(true,m.Content)
            }
        }

        if m.ChannelID == Conf.Bot.ChatChannelID {
            if !bdsStartLock {
                return
            }
	        if m.Content == Ls {
		        bds.ChatW(Ls)
                return
	        }
            fmt.Println(m.Author.Username)
            bds.ChatW(fmt.Sprintf("%s@%s",tool.DuplicateName(m.Author.Username,m.Member.Nick), m.Content))
        }
    })

    go func(){
        for chat := range bdsChat {
           logChan <- "bdsChat>>"+chat
	        if strings.HasPrefix(chat,bds.ChatChannel) {
	            MessageSend(false,chat)
		        continue
            }
	        if strings.Contains(chat,bds.MtBackupFileChat) {
		        backup.BackUp(chat,backChan)
		        continue
	        }
            MessageSend(true,chat)
        }
    }()

    go func () {
        inputReader := bufio.NewReader(os.Stdin)
        fmt.Println("cmd ok! <start> run bds")
        var input string
	    for {
            if bdsStartLock {
                input, _ = inputReader.ReadString('\n')
                fmt.Println(len(input))
                bds.RootW(false,input)
                if strings.Contains(input,"stop") {
                    time.Sleep(time.Second*2)
                }
                continue
            }else{
                input, _ := inputReader.ReadString('\n')
                if strings.Contains(input,"start") {
                    if len(input) > 7 {
                        continue
                    }
                    mc := bds.Bds{}
                    go mc.Start(&bdsStartLock)
                    time.Sleep(time.Second*3)
                    continue
                }
            }
        }
    }()

    sc := make(chan os.Signal, 1)
    signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
    <-sc
    if bdsStartLock {
        bds.RootW(false,"stop\n")
        time.Sleep(time.Second*2)
    }
    tool.LogFile.Close()
    S.Close()
}

func RootChannelPsutil() {
    fmt.Println("RootChannel-SleepTime:",time.Duration(Conf.Bot.ChannelTopic.RootSleepTime*1e9))
    var i int
    for {
       fmt.Println("RootTopicUpdate...")
        _, err := S.ChannelEditComplex(Conf.Bot.RootChannelID,&discordgo.ChannelEdit{Topic:ps.TopicTop(Conf.Bot.ChannelTopic.ChatTopicNetCard)})
        if err != nil {
            fmt.Println(err)
        }
        i++
        fmt.Println(i,"RootTopicUpdateOk")
        time.Sleep(time.Duration(Conf.Bot.ChannelTopic.RootSleepTime*1e9))
    }
}
func ChatChannelGameInfo() {
    fmt.Println("ChatChannel-SleepTime:",time.Duration(Conf.Bot.ChannelTopic.ChatSleepTime*1e9))
    var i int
    for {
        fmt.Println("ChatTopicUpdate...")
        if bdsStartLock {
            S.ChannelEditComplex(Conf.Bot.ChatChannelID,&discordgo.ChannelEdit{Topic:fmt.Sprintf("Online players:%d|Accumulated players:%d|RuningTime:%v",len(bds.Player),bds.PlayerAdd,tool.StartTime(&bds.BdsStartTime))})
        }else{
            S.ChannelEditComplex(Conf.Bot.ChatChannelID,&discordgo.ChannelEdit{Topic:bdsNo})
        }
        i++
        fmt.Println(i,"ChatTopicUpdateOk")
        time.Sleep(time.Duration(Conf.Bot.ChannelTopic.ChatSleepTime*1e9))
    }
}

func mCron(setCron string){
    log.Printf("setCron_%v\n",setCron)
    c := cron.New()
    c.AddFunc(setCron,mBackup)
	go c.Start()
}

func mBackup() {
    if bdsStartLock {
        if ps.HdTF() {
           go bds.Back(&backupLock)
	    }else{
            MessageSend(true,ps.HdErr)
	    }
	}else{
        log.Println(bdsNo)
        MessageSend(true,bdsNo)
    }
}
//正在运行，无法回滚备份
const RdErr = "Running, unable to roll back backup"

func mRd(rDir string,rdChat chan string) {
    if !bdsStartLock {
        backup.Rd(rDir,rdChat)
        MessageSend(true,<-rdChat)
	}else{
        log.Println(RdErr)
        MessageSend(true,RdErr)
    }
}

func MessageSend(tf bool, chat string) {
    if tf {
        S.ChannelMessageSend(Conf.Bot.RootChannelID, chat)
    }else{
        S.ChannelMessageSend(Conf.Bot.ChatChannelID, (chat[5:]))
    }
}
