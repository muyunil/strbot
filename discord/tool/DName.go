package tool

func DuplicateName(name, nick string) string {
    if nick != "" {
        return nick
    }
    return name
}
/*
import  "strings"
func DuplicateName(name, nick string,userList map[string]string) string {
    if nick != "" {
        return nick
    }
    tmp := strings.Split(name,"#")
    if v, ok := userList[tmp[0]]; ok {
        if name == v {
            return tmp[0]
        }else{
            return name
        }
    }else{
        userList[tmp[0]] = name
        return tmp[0]
    }
}
*/
