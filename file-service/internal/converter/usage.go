package converter

import "strconv"

const gigaByteMetric = "G"
const gigaByteMultiplier = 1073741824
const megaByteMetric = "M"
const megaByteMultiplier = 1048576

func ToIntBytes(strUsage string) int {
	metric := strUsage[len(strUsage)-1:]
	amount := strUsage[0 : len(strUsage)-1]

	amountInt, _ := strconv.Atoi(amount)

	if metric == gigaByteMetric {
		return amountInt * gigaByteMultiplier
	}

	if metric == megaByteMetric {
		return amountInt * megaByteMultiplier
	}

	return amountInt
}
