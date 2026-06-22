package golang
func maxProfit(prices []int) int {
    if len(prices) < 1 {
        return 0
    }
    maxProfit := 0
    minBuy := prices[0]
    for _, val := range prices {
        if val < minBuy {
            minBuy = val
        } else {
            if val - minBuy > maxProfit {
                maxProfit = val - minBuy
            }
        }
    }
    return maxProfit
}