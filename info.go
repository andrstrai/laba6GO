package main

import "fmt"

// PrintAll печатает короткий список фигур. fmt.Println сам вызывает String().
func PrintAll(shapes []Shape) {
	for _, s := range shapes {
		fmt.Println(s)
	}
}

// PrintShapesInfo печатает подробную информацию о каждой фигуре.
// Полиморфизм: функция не знает, какие фигуры в срезе, и вызывает методы интерфейса.
func PrintShapesInfo(shapes []Shape) {
	if len(shapes) == 0 {
		fmt.Println("Фигур нет")
		return
	}
	for i, s := range shapes {
		fmt.Printf("%d) %v\n", i+1, s)
		fmt.Printf("   тип: %s, площадь: %.2f, периметр: %.2f\n", s.Type(), s.Area(), s.Perimeter())
	}
}

// TotalArea возвращает сумму площадей любых фигур (вариадический параметр)
func TotalArea(shapes ...Shape) float64 {
	total := 0.0
	for _, s := range shapes {
		total += s.Area()
	}
	return total
}

// TotalPerimeter возвращает сумму периметров
func TotalPerimeter(shapes ...Shape) float64 {
	total := 0.0
	for _, s := range shapes {
		total += s.Perimeter()
	}
	return total
}

// FilterShapes — функция высшего порядка: принимает функцию keep
// и возвращает только те фигуры, для которых keep вернула true
func FilterShapes(shapes []Shape, keep func(Shape) bool) []Shape {
	var result []Shape
	for _, s := range shapes {
		if keep(s) {
			result = append(result, s)
		}
	}
	return result
}
