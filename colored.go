package main

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

// Color — цвет фигуры. Эта структура встраивается в ColoredShape.
type Color struct {
	Name string // название цвета, например "красный"
	code string // ANSI-код: терминал понимает его как «дальше пиши этим цветом»
}

// colors — какие цвета умеет программа и их ANSI-коды
var colors = map[string]string{
	"красный":    "\033[31m",
	"зеленый":    "\033[32m",
	"желтый":     "\033[33m",
	"синий":      "\033[34m",
	"фиолетовый": "\033[35m",
}

// resetColor возвращает обычный цвет текста
const resetColor = "\033[0m"

// NewColor создаёт цвет по названию или возвращает ошибку, если такого цвета нет
func NewColor(name string) (Color, error) {
	name = strings.ToLower(strings.TrimSpace(name))
	code, ok := colors[name]
	if !ok {
		return Color{}, fmt.Errorf("цвет %q не поддерживается", name)
	}
	return Color{Name: name, code: code}, nil
}

// Paint «раскрашивает» текст: ставит код цвета в начало и сброс в конец
func (c Color) Paint(text string) string {
	return c.code + text + resetColor
}

// ColorNames возвращает названия всех цветов по алфавиту (для подсказки в меню)
func ColorNames() []string {
	names := make([]string, 0, len(colors))
	for name := range colors {
		names = append(names, name)
	}
	sort.Strings(names) // у map нет порядка, поэтому сортируем
	return names
}

// ColoredShape — цветная фигура. Композиция вместо наследования:
//   - встроенный интерфейс Shape: внутри лежит любая фигура (круг, прямоугольник...);
//   - встроенная структура Color: цвет со своим методом Paint.
//
// Методы ReturnID, Area, Perimeter, Change поднимаются от встроенной фигуры,
// поэтому ColoredShape сама реализует Shape. Type и String мы переопределяем.
type ColoredShape struct {
	Shape
	Color
}

// NewColoredShape оборачивает фигуру s в цветную фигуру
func NewColoredShape(s Shape, colorName string) (Shape, error) {
	if s == nil {
		return nil, errors.New("нет фигуры для покраски")
	}
	color, err := NewColor(colorName)
	if err != nil {
		return nil, err
	}
	return &ColoredShape{Shape: s, Color: color}, nil
}

// Type переопределяет Type встроенной фигуры и добавляет к нему цвет
func (cs ColoredShape) Type() string {
	return cs.Shape.Type() + " (" + cs.Name + ")"
}

// String переопределяет String встроенной фигуры: тот же текст, но цветом
func (cs ColoredShape) String() string {
	return cs.Paint(cs.Shape.String() + ", цвет: " + cs.Name)
}

// Recolor меняет цвет. Pointer receiver, потому что метод меняет поле Color.
func (cs *ColoredShape) Recolor(colorName string) error {
	color, err := NewColor(colorName)
	if err != nil {
		return err
	}
	cs.Color = color
	return nil
}

// Paint красит фигуру: уже цветную перекрашивает, обычную оборачивает в ColoredShape
func Paint(s Shape, colorName string) (Shape, error) {
	if cs, ok := s.(*ColoredShape); ok { // это уже цветная фигура?
		if err := cs.Recolor(colorName); err != nil {
			return nil, err
		}
		return cs, nil
	}
	return NewColoredShape(s, colorName)
}
