package main

func main() {
	// example 01
	//var badMap1 = map[[]int]int{}
	//_ = badMap1

	// example 02
	//var badMap2 = map[interface{}]int{
	//	"1" : 1,
	//	[]int{2} : 2, // 这里会引发panic unhas of unhashable
	//	3: 3,
	//}
	//_  = badMap2

	// example 03
	//var badMap3 map[[1][]string]int // Compile error
	//_= badMap3

	// example 04
	//type BadKey1 struct {
	//	slice []string
	//}
	//var badMap4 map[BadKey1]int // 这里会引发编译错误
	//_ = badMap4

	// example 05
	//var badMap5 map[[1][2][3][]string]int // 这里会引发编译错误
	//_ = badMap5

	// example 06
	//type BadKey2Field1 struct {
	//	slice []string
	//}
	//
	//type BadKey2 struct {
	//	field BadKey2Field1
	//}
	//var badMap6 map[BadKey2]int // 这里会编译错误
	//_ = badMap6

}
