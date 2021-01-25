package tool

func DuplicateName(name, nick string) string {
    if nick != "" {
        return nick
    }
    return name
}
