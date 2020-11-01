package main

import "fmt"

// example 01
// AnimalCategory 代表动物分类学中的基本分类法。
type AnimalCategory struct {
	kingdom string // 界。
	phylum  string // 门。
	class   string // 纲。
	order   string // 目。
	family  string // 种。
	genus   string // 属。
	species string // 种。
}

func (ac AnimalCategory) String() string {
	return fmt.Sprintf("%s%s%s%s%s%s%s",
		ac.kingdom, ac.phylum, ac.class, ac.order,
		ac.family, ac.genus, ac.species)
}

// example 02
type Animal struct {
	scientificName string // 学名
	AnimalCategory        // 动物基本分类。
}

// 该方法会"屏蔽"掉入字段中的同名方法。
func (a Animal) String() string {
	return fmt.Sprintf("%s (category: %s)",
		a.scientificName, a.AnimalCategory)
}

// example 03
type Cat struct {
	name string
	Animal
}

// 该方法会"屏蔽"掉嵌入字段中的同名方法
func (cat Cat) String() string {
	return fmt.Sprintf("%s (category: %s, name: %q)",
		cat.scientificName, cat.AnimalCategory, cat.name)
}

func main() {
	// example 01
	category := AnimalCategory{species: "cat"}
	fmt.Printf("The animal category: %s\n", category)

	// example 02
	animal := Animal{
		scientificName: "American Shorthair",
		AnimalCategory: category,
	}
	fmt.Printf("The animal: %s\n", animal)

	// example 03
	cat := Cat{
		name:   "title pig",
		Animal: animal,
	}
	fmt.Printf("The cat :%s\n", cat)

}
