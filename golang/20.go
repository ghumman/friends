package golang
func isValid(s string) bool {
    stack := []byte{}
    checkMap := map[byte]byte{
        '(' : ')',
        '{' : '}',
        '[' : ']',
    }

    for i:=0; i<len(s); i++ {
        if _, ok:= checkMap[s[i]]; ok {
            stack = append(stack, s[i])
        } else {
            if len(stack)<= 0{
                return false
            }
            element := stack[len(stack) - 1]
            stack = stack[:(len(stack) - 1)]
            if (s[i] != checkMap[element]) {
                return false
            }
        }
    }
    return len(stack) == 0
}