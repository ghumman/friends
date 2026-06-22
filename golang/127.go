package golang
func ladderLength(beginWord string, endWord string, wordList []string) int {
    wordSet := make(map[string]bool)
    for _, word := range wordList {
        wordSet[word] = true
    }
    if !wordSet[endWord] {
        return 0
    }
    queue := []string{beginWord}
    steps := 1
    for len(queue) > 0 {
        levelLength := len(queue)
        for i:=0; i<levelLength; i++ {
            currentWord := queue[0]
            queue = queue[1:]
            if currentWord == endWord {
                return steps
            }
            currentWordRunes := []rune(currentWord)
            for j:=0; j<len(currentWord);j++ {
                original := currentWordRunes[j]
                for k:='a'; k<='z'; k++ {
                    currentWordRunes[j] = k
                    newWord := string(currentWordRunes)
                    if newWord == currentWord {
                        continue
                    }
                    if wordSet[newWord] {
                        queue = append(queue, newWord)
                        delete(wordSet, newWord)
                    }
                }
                currentWordRunes[j] = original
            }
        }
        steps++
    }
    return 0
}