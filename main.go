package main

import (
	"bufio"
	"errors"
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

// Shape — общий интерфейс всех фигур.
// Area, Perimeter, String и Type требует задание, ReturnID и Change нужны нашему меню.
type Shape interface {
	ReturnID() uint
	Area() float64
	Perimeter() float64
	String() string
	Type() string
	Change(reader *bufio.Reader) error
}

var CURR_ID uint = 0

// readPositive печатает prompt, читает строку и превращает её в положительное число.
// Возвращает ошибку, если ввод не удалось прочитать, это не число или число <= 0.
func readPositive(reader *bufio.Reader, prompt string) (float64, error) {
	fmt.Println(prompt)
	input, err := reader.ReadString('\n')
	if err != nil {
		return 0, fmt.Errorf("не удалось прочитать ввод: %w", err)
	}
	input = strings.TrimSpace(input)
	x, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return 0, fmt.Errorf("%q не число", input)
	}
	if x <= 0 || math.IsInf(x, 0) || math.IsNaN(x) {
		return 0, fmt.Errorf("%v: нужно положительное конечное число", x)
	}
	return x, nil
}

// Круг
type Circle struct {
	ID     uint
	Radius float64
}

// value receiver: эти методы только читают поля
func (c Circle) ReturnID() uint     { return c.ID }
func (c Circle) Type() string       { return "Круг" }
func (c Circle) Area() float64      { return math.Pi * c.Radius * c.Radius }
func (c Circle) Perimeter() float64 { return 2 * math.Pi * c.Radius }
func (c Circle) String() string {
	return fmt.Sprintf("ID: %d. Круг с радиусом %.2f", c.ID, c.Radius)
}

// pointer receiver: метод меняет поле настоящего круга
func (c *Circle) Change(reader *bufio.Reader) error {
	radius, err := readPositive(reader, "Введите новый радиус:")
	if err != nil {
		return err
	}
	c.Radius = radius
	return nil
}

// Прямоугольник
type Rectangle struct {
	ID            uint
	Width, Height float64
}

func (r Rectangle) ReturnID() uint     { return r.ID }
func (r Rectangle) Type() string       { return "Прямоугольник" }
func (r Rectangle) Area() float64      { return r.Height * r.Width }
func (r Rectangle) Perimeter() float64 { return 2 * (r.Height + r.Width) }
func (r Rectangle) String() string {
	return fmt.Sprintf("ID: %d. Прямоугольник %.2f x %.2f", r.ID, r.Width, r.Height)
}

func (r *Rectangle) Change(reader *bufio.Reader) error {
	width, err := readPositive(reader, "Введите ширину прямоугольника:")
	if err != nil {
		return err
	}
	height, err := readPositive(reader, "Введите высоту прямоугольника:")
	if err != nil {
		return err
	}
	r.Width = width
	r.Height = height
	return nil
}

// Треугольник
type Triangle struct {
	ID      uint
	A, B, C float64 // Длины сторон
}

var ErrNoTriangle = errors.New("треугольника с такими сторонами не существует")

func exists(A float64, B float64, C float64) bool {
	maxi := max(A, B, C)
	mini := min(A, B, C)
	mid := A + B + C - mini - maxi
	return maxi < mini+mid && !floatEqual(maxi, mini+mid)
}

func (t Triangle) ReturnID() uint     { return t.ID }
func (t Triangle) Type() string       { return "Треугольник" }
func (t Triangle) Perimeter() float64 { return t.A + t.B + t.C }

// Площадь по формуле Герона
func (t Triangle) Area() float64 {
	p := t.Perimeter() / 2
	return math.Sqrt(p * (p - t.A) * (p - t.B) * (p - t.C))
}
func (t Triangle) String() string {
	return fmt.Sprintf("ID: %d. Треугольник (стороны: %.2f, %.2f, %.2f)", t.ID, t.A, t.B, t.C)
}

func (t *Triangle) Change(reader *bufio.Reader) error {
	a, err := readPositive(reader, "Введите первую сторону треугольника:")
	if err != nil {
		return err
	}
	b, err := readPositive(reader, "Введите вторую сторону треугольника:")
	if err != nil {
		return err
	}
	c, err := readPositive(reader, "Введите третью сторону треугольника:")
	if err != nil {
		return err
	}
	if !exists(a, b, c) {
		return ErrNoTriangle
	}
	t.A, t.B, t.C = a, b, c
	return nil
}

// Фабричные функции (конструкторы). Проверяют данные, поэтому создать фигуру
// в неправильном состоянии нельзя. Внутри создаётся указатель &Circle{...}:
// у Change pointer receiver, и интерфейс Shape реализует только *Circle, а не Circle.
// Возвращаем интерфейс Shape, чтобы при ошибке вернуть настоящий nil.
func NewCircle(radius float64) (Shape, error) {
	if radius <= 0 {
		return nil, fmt.Errorf("радиус %v должен быть положительным", radius)
	}
	CURR_ID++
	return &Circle{ID: CURR_ID, Radius: radius}, nil
}

func NewRectangle(width, height float64) (Shape, error) {
	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("стороны %v и %v должны быть положительными", width, height)
	}
	CURR_ID++
	return &Rectangle{ID: CURR_ID, Width: width, Height: height}, nil
}

func NewTriangle(a, b, c float64) (Shape, error) {
	if a <= 0 || b <= 0 || c <= 0 {
		return nil, fmt.Errorf("стороны %v, %v, %v должны быть положительными", a, b, c)
	}
	if !exists(a, b, c) {
		return nil, ErrNoTriangle
	}
	CURR_ID++
	return &Triangle{ID: CURR_ID, A: a, B: b, C: c}, nil
}

func findShapeByID(shapes []Shape, id int) (int, error) {
	for i, s := range shapes {
		if int(s.ReturnID()) == id {
			return i, nil
		}
	}
	return -1, fmt.Errorf("фигура с ID %d не найдена", id)
}

func floatEqual(a, b float64) bool {
	const eps = 1e-9
	diff := math.Abs(a - b)
	scale := math.Max(math.Abs(a), math.Abs(b))
	return diff <= eps*math.Max(1, scale)
}

// readInt читает целое число (пункт меню, ID)
func readInt(reader *bufio.Reader, prompt string) (int, error) {
	fmt.Println(prompt)
	input, err := reader.ReadString('\n')
	if err != nil {
		return 0, fmt.Errorf("не удалось прочитать ввод: %w", err)
	}
	n, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil {
		return 0, errors.New("нужно ввести целое число")
	}
	return n, nil
}

// readShape спрашивает пользователя, какую фигуру создать, и создаёт её
func readShape(reader *bufio.Reader) (Shape, error) {
	kind, err := readInt(reader, "Какую фигуру добавить? 1 - круг, 2 - прямоугольник, 3 - треугольник")
	if err != nil {
		return nil, err
	}
	switch kind {
	case 1:
		r, err := readPositive(reader, "Введите радиус:")
		if err != nil {
			return nil, err
		}
		return NewCircle(r)
	case 2:
		w, err := readPositive(reader, "Введите ширину:")
		if err != nil {
			return nil, err
		}
		h, err := readPositive(reader, "Введите высоту:")
		if err != nil {
			return nil, err
		}
		return NewRectangle(w, h)
	case 3:
		a, err := readPositive(reader, "Введите первую сторону:")
		if err != nil {
			return nil, err
		}
		b, err := readPositive(reader, "Введите вторую сторону:")
		if err != nil {
			return nil, err
		}
		c, err := readPositive(reader, "Введите третью сторону:")
		if err != nil {
			return nil, err
		}
		return NewTriangle(a, b, c)
	default:
		return nil, fmt.Errorf("нет фигуры с номером %d", kind)
	}
}

// startShapes создаёт фигуры, с которыми программа начинает работу
func startShapes() []Shape {
	circle, _ := NewCircle(5) // данные заведомо правильные, ошибки быть не может
	rect, _ := NewRectangle(4, 6)
	tri, _ := NewTriangle(3, 4, 5)
	red, _ := NewColoredShape(tri, "красный")
	return []Shape{circle, rect, red}
}

// основная функция
func main() {
	all_shapes := startShapes()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("\nДля продолжения нажмите Enter...")
		if _, err := reader.ReadString('\n'); err != nil {
			return // ввод закончился (Ctrl+D) — выходим, а не крутимся вечно
		}
		clear_screen()

		fmt.Println("Меню команд:")
		fmt.Println("1 - Показать все фигуры")
		fmt.Println("2 - Изменить фигуру")
		fmt.Println("3 - Добавить фигуру")
		fmt.Println("4 - Покрасить фигуру")
		fmt.Println("5 - Общая площадь и периметр")
		fmt.Println("6 - Фигуры с площадью больше заданной")
		fmt.Println("0 - Завершить работу")

		a, err := readInt(reader, "Выберите пункт из меню управления:")
		if err != nil {
			fmt.Println(err)
			continue
		}

		switch a {
		case 0:
			return
		case 1:
			PrintShapesInfo(all_shapes)
		case 2:
			PrintAll(all_shapes)
			shapeID, err := readInt(reader, "\nВведите ID фигуры, которую хотите изменить:")
			if err != nil {
				fmt.Println(err)
				continue
			}
			shapeIndex, err := findShapeByID(all_shapes, shapeID)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("Текущая фигура:", all_shapes[shapeIndex])
			if err := all_shapes[shapeIndex].Change(reader); err != nil {
				fmt.Println("Фигура не изменена:", err)
				continue
			}
			fmt.Println("Фигура после изменения:", all_shapes[shapeIndex])
		case 3:
			s, err := readShape(reader)
			if err != nil {
				fmt.Println("Фигура не добавлена:", err)
				continue
			}
			all_shapes = append(all_shapes, s)
			fmt.Println("Добавлена:", s)
		case 4:
			PrintAll(all_shapes)
			shapeID, err := readInt(reader, "\nВведите ID фигуры, которую хотите покрасить:")
			if err != nil {
				fmt.Println(err)
				continue
			}
			shapeIndex, err := findShapeByID(all_shapes, shapeID)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("Доступные цвета:", strings.Join(ColorNames(), ", "))
			fmt.Println("Введите цвет:")
			color, _ := reader.ReadString('\n')
			painted, err := Paint(all_shapes[shapeIndex], strings.TrimSpace(color))
			if err != nil {
				fmt.Println(err)
				continue
			}
			all_shapes[shapeIndex] = painted
			fmt.Println("Готово:", painted)
		case 5:
			fmt.Printf("Общая площадь всех фигур: %.2f\n", TotalArea(all_shapes...))
			fmt.Printf("Общий периметр всех фигур: %.2f\n", TotalPerimeter(all_shapes...))
		case 6:
			limit, err := readPositive(reader, "Показать фигуры с площадью больше:")
			if err != nil {
				fmt.Println(err)
				continue
			}
			big := FilterShapes(all_shapes, func(s Shape) bool { return s.Area() > limit })
			if len(big) == 0 {
				fmt.Println("Таких фигур нет")
				continue
			}
			PrintShapesInfo(big)
		default:
			fmt.Println("Нет такого пункта меню!")
		}
	}
}
