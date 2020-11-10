package tool

import (
    "encoding/base64"
    "fmt"
    "io/ioutil"
    "net/http"
)

func ImageBase(name string) string {
    fmt.Println("Tool.ImageBase,Name:"+name)
    imgUrl := fmt.Sprintf("https://minotar.net/helm/%s/64.jpg",name)
    //获取远端图片
    res, err := http.Get(imgUrl)
    if err != nil {
        fmt.Println("Tool.ImageBase-getImgUrlErr!")
        return ""
    }
    defer res.Body.Close()
    // 读取获取的[]byte数据
    data, _ := ioutil.ReadAll(res.Body)
    fileType := http.DetectContentType(data)

    imageBase64 := base64.StdEncoding.EncodeToString(data)
    dataUrl := fmt.Sprintf("data:%s;base64,%s", fileType, imageBase64)
    return dataUrl
}
