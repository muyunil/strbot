# discordbot

基本功能
游戏与频道消息互相转发

管理频道发送help/帮助  获取命令帮助
管理员基本维护命令

> 玩家频道显示当前在线人数，服务器加入人数，服务器在线时间
管理员频道显示服务器负载

-配置文件需要：  
  服务器id  
  服务器管理频道id  
  Chat信息转发频道id  
  Bot Token  

备份:   
. 灵活的Q群消息触发备份 还有Cron计划备份 修改shibot.yaml配置文件 设置Cron触发时间  
Cron *  *  *  *  *  *   

> 一共6个星意义和可选值分别是： 
秒（0-60）分（0-59）时（0-24） 天（0-31） 月（1-12） 周（0-6）

管理员命令:   
     ps :查看服务器状态  
     start :启动bds  
     stop :关闭bds  
     backup :备份worlds(bds未启动时无法使用)  
     lsbd :查看今天的备份存档列表  
        可选附加time查看其他日期的备份列表如：  
        lsbd 10/10 :即为查看10月10号的备份列表  
     rd name-time :lsbd 获取备份列表选一个执行>如：  
        rd worlds-2020-10-10_10-10-10  
        即为回滚到此name-time指向的备份存档  
        ps:bds运行时无法回滚  
    wl 白名单命令可选操作：  
        wl + tes  tID | wl - tes  tID :id自动附带""  
    /cmd :/ 开头的命令会发送给bds控制台如  
        /say HelloWorlds  
        /stop  

