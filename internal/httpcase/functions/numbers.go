package functions

import "strconv"

type Numbers struct {
}

func (f *Numbers) Add(a, b string) int {
	x, y := f.parseInt(a, b)
	return x + y
}

func (f *Numbers) Sub(a, b string) int {
	x, y := f.parseInt(a, b)
	return x - y
}
func (f *Numbers) Multiply(a, b string) int {
	x, y := f.parseInt(a, b)
	return x * y
}

func (f *Numbers) Divide(a, b string) int {
	x, y := f.parseInt(a, b)
	return x / y
}

func (f *Numbers) Mod(a, b string) int {
	x, y := f.parseInt(a, b)
	return x % y
}

func (f *Numbers) parseInt(a, b string) (int, int) {
	int1, _ := strconv.Atoi(a)
	int2, _ := strconv.Atoi(b)
	return int1, int2
}
