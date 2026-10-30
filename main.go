package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

// функция для очистки экрана (работает и на Windows, и на macOS/Linux)
func clear_screen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// интерфейс shape
type shape interface {
	ReturnID() uint
	Area() float64
	Perimeter() float64
	String() string
	Change()
}

var CURR_ID uint = 0

// Круг
type Circle struct {
	ID     uint
	radius float64
}

// используется value receiver
func (c Circle) ReturnID() uint     { return c.ID }
func (c Circle) Area() float64      { return math.Pi * math.Pow(c.radius, 2) }
func (c Circle) Perimeter() float64 { return 2 * math.Pi * c.radius }
func (c Circle) String() string {
	return "ID: " + fmt.Sprint(c.ID) + ". Круг с радиусом = " + fmt.Sprint(c.radius)
}

// pointer receiver
func (c *Circle) Change() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Ввведите новый радиус:")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	radius, err := strconv.ParseFloat(input, 64)
	if err != nil || radius < 0 {
		fmt.Println("Некорректный ввод! Пожалуйста, введите положительное число.")
		return
	}
	c.radius = radius
}

// Прямоугольник
type Rectangle struct {
	ID            uint
	Width, Height float64
}

func (r Rectangle) ReturnID() uint     { return r.ID }
func (r Rectangle) Area() float64      { return r.Height * r.Width }
func (r Rectangle) Perimeter() float64 { return 2*r.Height + 2*r.Width }
func (r Rectangle) String() string {
	return "ID: " + fmt.Sprint(r.ID) + ". Прямоугольник со сторонами: " + fmt.Sprint(r.Width) + " и " + fmt.Sprint(r.Height)
}
func (r *Rectangle) Change() {}

// Треугольник
type Triangle struct {
	ID      uint
	A, B, C float64 // Длины сторон
}

func (t Triangle) ReturnID() uint     { return t.ID }
func (t Triangle) Perimeter() float64 { return t.A + t.B + t.C }

// Площадь по формуле Герона
func (t Triangle) Area() float64 {
	p := t.Perimeter() / 2
	return math.Sqrt(p * (p - t.A) * (p - t.B) * (p - t.C))
}
func (t Triangle) String() string {
	return fmt.Sprintf("ID: %d. Треугольник (стороны: %.2f, %.2f, %.2f)", t.ID, t.A, t.B, t.C)
}
func (r *Triangle) Change() {}

// Фабричные функции
func NewCircle(radius float64) shape {
	CURR_ID++
	return &Circle{ID: CURR_ID, radius: radius}
}

func NewRectangle(width, height float64) shape {
	CURR_ID++
	return &Rectangle{ID: CURR_ID, Width: width, Height: height}
}

func NewTriangle(a, b, c float64) shape {
	CURR_ID++
	return &Triangle{ID: CURR_ID, A: a, B: b, C: c}
}

func findShapeByID(shapes []shape, id int) (int, error) {
	if id >= int(CURR_ID+1) || id < 1 {
		return -1, fmt.Errorf("фигура с ID %d не найдена", id)
	}
	for i, s := range shapes {
		if int(s.ReturnID()) == id {
			return i, nil
		}
	}
	return -1, fmt.Errorf("товар с ID %d не найден", id)
}

func PrintAll(shapes []shape) {
	for _, s := range shapes {
		fmt.Println(s)
	}
}

// основная функция
func main() {
	// Создание фигур через фабричные функции
	all_shapes := []shape{NewCircle(5), NewRectangle(4, 6), NewTriangle(3, 4, 5)}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("Для начала работы пожалуйста, нажмите Enter...")
		reader.ReadString('\n')
		clear_screen()

		fmt.Println("Меню команд:")
		fmt.Println("1 - Показать все фигуры")
		fmt.Println("2 - Изменить фигуру")
		fmt.Println("Выберите пункт из меню управления: ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		a, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Некорректный ввод! Пожалуйста, введите число.")
			continue
		}

		switch a {
		case 0:
			return
		case 1:
			fmt.Println("Все фигуры:")
			PrintAll(all_shapes)
		case 2:
			fmt.Println("Список фигур")
			PrintAll(all_shapes)
			fmt.Println("\nВведите ID фигуры, которую хотите изменить:")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)

			shapeID, err := strconv.Atoi(input)
			if err != nil {
				fmt.Println("Некорректный ввод! Пожалуйста, введите число.")
				continue
			}
			shapeIndex, err := findShapeByID(all_shapes, shapeID)
			if err != nil {
				fmt.Println(err)
				continue
			}
			all_shapes[shapeIndex].Change()
			fmt.Println("\nФигура после изменения:")
			fmt.Println(all_shapes[shapeIndex])
		default:
			fmt.Println("Нет такого пункта меню!")
		}
	}
}
