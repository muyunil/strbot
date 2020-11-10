package bot

import (
    "fmt"
    "github.com/catsworld/qq-bot-api"
    "shibot/config"
    "strings"
)
var (
    bot *qqbotapi.BotAPI
    conf *config.JsonConfig
)

const (
    Ls = "ls"
    Ps = "ps"
    Rd = "rd "
    LsBd = "lsbd"

    StartBds = "启动"
    StopBds = "关闭"
    Backup = "备份"
)

func Start(qqChat chan string) {
    fmt.Println("正在启动")

    if config.PathExists("shibot.json") {
        conf = config.Load("shibot.json")
    }
    fmt.Printf("confT %T\n",conf)
    if conf == nil {
	err := config.DefaultConfig().Save("shibot.json")
	if err != nil {
        fmt.Println("创建默认配置文件时出现错误: %v", err)
		return
	}
	fmt.Println("默认配置文件已生成, 请编辑 shibot.json 后重启程序.")
	return
    }
    if conf.Master == 0 {
	fmt.Println("请修改 shibot.json 以添加管理员qq")
	return
    }
    fmt.Println("加载配置完毕")

    url := fmt.Sprintf("ws://%v:%v", conf.Host, conf.Port)
    var err1 error
    bot, err1 = qqbotapi.NewBotAPI("", url, "")
    if err1 != nil {
	fmt.Println(err1)
    }
    bot.Debug = true

    u := qqbotapi.NewUpdate(0)
    u.PreloadUserInfo = true

    updates, err := bot.GetUpdatesChan(u)
    if err != nil {
	fmt.Println(err)
    }
    fmt.Println("启动完毕，正在运行")

//判断是否设置了Cron 如不是默认 则设置定时计划
    if conf.Cron != "* * * * * *" {
	qqChat <- ("Cron_" + conf.Cron)
    }

    for update := range updates {
	if update.Message == nil {
		continue
	}
	if update.GroupID != conf.Qqun {
	    continue
	}
	fmt.Printf("[%s] %s\n", update.Message.From.String(), update.Message.Text)

	if update.UserID != conf.Master {
	    if update.Message.Text == Ls {
		qqChat <- Ls
		continue
	    }
	    if update.Message.Text == Ps {
		qqChat <- Ps
		continue
	    }
	    if update.Message.Text != "" {
		qqChat <-  update.Message.From.String() + "@" +  update.Message.Text
		continue
	    }
	}else{
	    if update.Message.Text == Ps {
		    qqChat <- Ps
		    continue
	    }
	    if update.Message.Text == Ls {
		    qqChat <- Ls
		    continue
	    }
	    if update.Message.Text == StartBds {
		    qqChat <- StartBds
		    continue
        }
	    if update.Message.Text == StopBds {
	        qqChat <- StopBds
		    continue
        }
	    if update.Message.Text == Backup {
	        qqChat <- Backup
		    continue
        }
	    if strings.HasPrefix(update.Message.Text, LsBd) {
		    i := strings.Index(update.Message.Text," ")
            if i == 4 {
	            qqChat <- update.Message.Text
		        continue
            }
            if len(update.Message.Text) == 4 {
	            qqChat <- update.Message.Text
		        continue
            }
        }
	    if strings.HasPrefix(update.Message.Text, "su " ) {
	        qqChat <- update.Message.Text
		    continue
        }
	    if strings.HasPrefix(update.Message.Text, Rd) {
            if len(update.Message.Text) > 23 {
                qqChat <- update.Message.Text
            }
        }
	    if strings.HasPrefix(update.Message.Text, "su " ) {
	        qqChat <- update.Message.Text
		    continue
        }
	    if update.Message.Text != "" {
		    qqChat <- update.Message.From.String() + "@" +  update.Message.Text
	    }
	}
    }
}
func Printqq(pchat string) {
    chat := bot.NewMessage(conf.Qqun, "group").
//		At("1232332333").
		Text(pchat).
//		NewLine().
//		FaceByName("调皮").
//		Text("这是一个测试").
//		ImageBase64("img.jpg").
		Send()
    if chat.Err != nil {
	    bot.DeleteMessage(chat.Result.MessageID)
    }

}
