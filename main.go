package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell"
)

const SnakeSymbol = 0x2588 //0x0090
const AppleSymbol = 0x25CF
const GameFrameWidth = 30
const GameFrameHeight = 15
const GameFrameSymbol = '|' //0x01C1

type GameObject struct {
	row, col, width, height int
	velRow, velCol          int
	symbol                  rune
}

var screen tcell.Screen
var isGamePause bool
var debugLog string

var gameObjects []*GameObject

func main() {
	InitScreen()
	InitGameState()
	inputChan := InitUserInput()

	for {
		HandleUserInput(ReadInput(inputChan))
		UpdateState()
		DrawState()

		time.Sleep(75 * time.Millisecond)
	}

	screen.Fini()
}

func UpdateState() {
	if isGamePause {
		return
	}

	for i := range gameObjects {
		gameObjects[i].row += gameObjects[i].velRow
		gameObjects[i].col += gameObjects[i].velCol
	}
}

func DrawState() {
	if isGamePause {
		return
	}

	screen.Clear()
	PrintString(0, 0, debugLog)
	PrintGameFrame()
	for _, obj := range gameObjects {
		PrintFilledRect(obj.row, obj.col, obj.width, obj.height, obj.symbol)
	}

	screen.Show()
}

func InitScreen() {
	var err error
	screen, err := tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	if err := screen.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	defStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	screen.SetStyle(defStyle)
}

func InitGameState() {
	gameObjects = []*GameObject{}
}

func HandleUserInput(key string) {
	if key == "Rune[q]" {
		screen.Fini()
		os.Exit(0)
	}
}

func InitUserInput() chan string {
	inputChan := make(chan string)
	go func() {
		switch ev := screen.PollEvent().(type) {
		case *tcell.EventKey:
			inputChan <- ev.Name()
		}
	}()

	return inputChan
}

func ReadInput(inputChan chan string) string {
	var key string
	select {
	case key = <-inputChan:
	default:
		key = ""
	}
	return key
}

func PrintGameFrame() {
	//TODO
}

func PrinStringCentered(row, col int, str string) {
	col = col - len(str)/2
	PrintString(row, col, str)
}

func PrintString(row, col int, str string) {
	for _, c := range str {
		PrintFilledRect(row, col, 1, 1, c)
		col += 1
	}
}

func PrintFilledRect(row, col, width, height int, ch rune) {
	for r := 0; r < height; r++ {
		for c := 0; c < width; c++ {
			screen.SetContent(col+c, row+r, ch, nil, tcell.StyleDefault)
		}
	}
}
