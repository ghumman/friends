package golang
func groupAnagrams(strs []string) [][]string {
    ana := make(map[[26]int][]string)
    result := [][]string{}
    for _, str := range strs {
        common := [26]int{}
        for _, ch := range str {
            common[ch - 'a']++
        }
        ana[common] = append(ana[common], str)
    }

    for _, value := range ana {
        result = append(result, value)
    }

    return result
}