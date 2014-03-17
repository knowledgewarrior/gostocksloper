package main
func main() {
	a := []float64{1, 3, 5}
	b := [...]float64 { 0.0, 0, 0 }	// compiles, runs
//	b := make([]float64, len(a)) 	// compiles, runs
	for i, _ := range a { b[i] = a[i] }
	for _,v  := range b { print(v," ") }

}
