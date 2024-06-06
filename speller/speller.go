//go:build !solution

package speller

func AddDelimetr(n int64, delimetr string) string {
	if n == 0 {
		return ""
	}
	return delimetr + Spell(n)
}

func Spell(n int64) string {
	if n < 0 {
		return "minus " + Spell(n*(-1))
	}
	if n == 0 {
		return "zero"
	}
	numberToNineteen := map[int64]string{
		1: "one", 2: "two", 3: "three", 4: "four", 5: "five", 6: "six",
		7: "seven", 8: "eight", 9: "nine", 10: "ten", 11: "eleven", 12: "twelve",
		13: "thirteen", 14: "fourteen", 15: "fifteen", 16: "sixteen",
		17: "seventeen", 18: "eighteen", 19: "nineteen",
	}
	numberTens := map[int64]string{
		20: "twenty", 30: "thirty", 40: "forty", 50: "fifty", 60: "sixty", 70: "seventy", 80: "eighty", 90: "ninety",
	}
	more := map[int64]string{
		100: "hundred", 1_000: "thousand", 1_000_000: "million", 1_000_000_000: "billion",
	}
	number := []int64{1_000_000, 1_000_000_000, 1_000_000_000_000}
	if n < 20 {
		return numberToNineteen[n]
	}
	if n < 100 {
		return numberTens[10*(n/10)] + AddDelimetr(n%10, "-")
	}
	if n < 1000 {
		return numberToNineteen[n/100] + " " + more[100] + AddDelimetr(n%100, " ")
	}
	for _, num := range number {
		id := num / 1_000
		if n < num {

			return Spell(n/id) + " " + more[id] + AddDelimetr(n%id, " ")
		}
	}
	panic("error")

}
