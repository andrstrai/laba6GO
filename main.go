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
	Area() float64
	Perimeter() float64
	String() string
	Change()
}

// Круг
type Circle struct {
	radius float64
}

// используется value receiver
func (c Circle) Area() float64      { return math.Pi * math.Pow(c.radius, 2) }
func (c Circle) Perimeter() float64 { return 2 * math.Pi * c.radius }
func (c Circle) String() string     { return "Круг с радиусом = " + fmt.Sprint(c.radius) }

// pointer receiver
func (c *Circle) Change() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Ввведите новый радиус:")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	radius, err := strconv.ParseFloat(input, 64)
	if err != nil {
		fmt.Println("Некорректный ввод! Пожалуйста, введите положительное число.")
		return
	}
	c.radius = radius
}

// Прямоугольник
type Rectangle struct {
	Width, Height float64
}

func (r Rectangle) Area() float64      { return r.Height * r.Width }
func (r Rectangle) Perimeter() float64 { return 2*r.Height + 2*r.Width }
func (r Rectangle) String() string {
	return "Прямоугольник со сторонами: " + fmt.Sprint(r.Width) + " и " + fmt.Sprint(r.Height)
}
func (r *Rectangle) Change() {}

// Треугольник
type Triangle struct {
	A, B, C float64 // Длины сторон
}

func (t Triangle) Perimeter() float64 { return t.A + t.B + t.C }

// Площадь по формуле Герона
func (t Triangle) Area() float64 {
	p := t.Perimeter() / 2
	return math.Sqrt(p * (p - t.A) * (p - t.B) * (p - t.C))
}
func (t Triangle) String() string {
	return fmt.Sprintf("Треугольник (стороны: %.2f, %.2f, %.2f)", t.A, t.B, t.C)
}
func (r *Triangle) Change() {}

// Фабричные функции
func NewCircle(radius float64) shape {
	return &Circle{radius: radius}
}

func NewRectangle(width, height float64) shape {
	return &Rectangle{Width: width, Height: height}
}

func NewTriangle(a, b, c float64) shape {
	return &Triangle{A: a, B: b, C: c}
}

// основная функция
func main() {
	// Создание фигур через фабричные функции
	circle := NewCircle(5)
	rectangle := NewRectangle(4, 6)
	triangle := NewTriangle(3, 4, 5)

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("Для начала работы пожалуйста, нажмите Enter...")
		reader.ReadString('\n')
		clear_screen()

		fmt.Println("Меню команд:")
		fmt.Println("1 - Показать все фигуры")
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
			fmt.Println(circle)
			fmt.Println(rectangle)
			fmt.Println(triangle)
		case 2:
			circle.Change()
		default:
			fmt.Println("Нет такого пункта меню!")
		}
	}
}
