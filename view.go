package main

import (
	"fmt"
	"strings"

	"github.com/k0kubun/termbox-go"
)

const (
	background = `
		WWWWWWWWWWWW WWWWWW
		WkkkkkkkkkkW WkkkkW
		WkkkkkkkkkkW WkkkkW
		WkkkkkkkkkkW WkkkkW
		WkkkkkkkkkkW WkkkkW
		WkkkkkkkkkkW WWWWWW
		WkkkkkkkkkkW
		WkkkkkkkkkkW
		WkkkkkkkkkkW BBBBBB
		WkkkkkkkkkkW WWWWWW
		WkkkkkkkkkkW
		WkkkkkkkkkkW
		WkkkkkkkkkkW BBBBBB
		WkkkkkkkkkkW WWWWWW
		WkkkkkkkkkkW
		WkkkkkkkkkkW BBBBBB
		WkkkkkkkkkkW WWWWWW
		WkkkkkkkkkkW
		WkkkkkkkkkkW
		WWWWWWWWWWWW

		kkkkkkkkkkkkkkkkkkk
		WWWWWWWWWWWWWWWWWWW
	`
	boardXOffset, boardYOffset       = 3, 2
	nextMinoXOffset, nextMinoYOffset = 16, 2
)

var (
	colorMapping = map[rune]termbox.Attribute{
		'k': termbox.ColorBlack,
		'K': termbox.ColorBlack | termbox.AttrBold,
		'r': termbox.ColorRed,
		'R': termbox.ColorRed | termbox.AttrBold,
		'g': termbox.ColorGreen,
		'G': termbox.ColorGreen | termbox.AttrBold,
		'y': termbox.ColorYellow,
		'Y': termbox.ColorYellow | termbox.AttrBold,
		'b': termbox.ColorBlue,
		'B': termbox.ColorBlue | termbox.AttrBold,
		'm': termbox.ColorMagenta,
		'M': termbox.ColorMagenta | termbox.AttrBold,
		'c': termbox.ColorCyan,
		'C': termbox.ColorCyan | termbox.AttrBold,
		'w': termbox.ColorWhite,
		'W': termbox.ColorWhite | termbox.AttrBold,
	}
)

func refreshScreen() {
	rewriteScreen(func() {
		drawDropMarker()
		drawCurrentMino()
		if clock.gameover {
			for j := 0; j < boardHeight; j++ {
				board.colorizeLine(j, termbox.ColorBlack)
			}
			drawText(10, 4, "GAME OVER", termbox.ColorWhite, termbox.ColorBlack)
			drawText(7, 6, "<SPC> to continue", termbox.ColorWhite, termbox.ColorBlack)

			ranking := NewRanking()
			for idx, sc := range ranking.scores {
				drawText(9, 8+idx, fmt.Sprintf("%2d: %6d", idx+1, sc), termbox.ColorWhite, termbox.ColorBlack)
			}
		}
	})
}

func rewriteScreen(rewrite func()) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	drawBacks(background, 0, 0)
	drawCells(board.text(), boardXOffset, boardYOffset)
	drawNextMino()
	drawTexts()
	rewrite()

	termbox.Flush()
}

func drawTexts() {
	drawText(32, 9, "SCORE", termbox.ColorWhite, termbox.ColorBlue)
	drawText(32, 10, fmt.Sprintf("%7d", score), termbox.ColorBlack, termbox.ColorWhite)

	drawText(32, 13, "LEVEL", termbox.ColorWhite, termbox.ColorBlue)
	drawText(32, 14, fmt.Sprintf("%5d", level), termbox.ColorBlack, termbox.ColorWhite)

	drawText(32, 16, "LINES", termbox.ColorWhite, termbox.ColorBlue)
	drawText(32, 17, fmt.Sprintf("%5d", deleteLines), termbox.ColorBlack, termbox.ColorWhite)

	drawText(3, 22, "  ←     z     <SPC>    x,↑   →", termbox.ColorWhite, termbox.ColorBlack)
	drawText(3, 23, " left     ↺   drop      ↻  right", termbox.ColorBlack, termbox.ColorWhite)

	drawText(30, 19, " p: pause", termbox.ColorWhite, termbox.ColorDefault)
	drawText(30, 20, " q: quit", termbox.ColorWhite, termbox.ColorDefault)
	drawText(30, 21, " r: reset", termbox.ColorWhite, termbox.ColorDefault)
}

func drawText(x, y int, text string, fg, bg termbox.Attribute) {
	for index, ch := range text {
		termbox.SetCell(x+index, y, rune(ch), fg, bg)
	}
}

func drawCurrentMino() {
	drawMino(currentMino, boardXOffset, boardYOffset)
}

func drawDropMarker() {
	marker := *currentMino
	marker.putBottom()

	lines := strings.Split(marker.block, "\n")
	for y, line := range lines {
		for x, char := range line {
			if isOnBoard(x+marker.x, y+marker.y) && colorByChar(char) != blankColor &&
				colorByChar(char) != termbox.ColorDefault {
				drawCell(x+marker.x+boardXOffset, y+marker.y+boardYOffset, colorByChar('K'))
			}
		}
	}
}

func drawNextMino() {
	drawMino(nextMino, nextMinoXOffset-nextMino.x, nextMinoYOffset-nextMino.y)
}

func drawMino(mino *Mino, xOffset, yOffset int) {
	lines := strings.Split(mino.block, "\n")

	for y, line := range lines {
		for x, char := range line {
			if isOnBoard(x+mino.x, y+mino.y) {
				color := colorByChar(char)
				drawCell(x+mino.x+xOffset, y+mino.y+yOffset, color)
			}
		}
	}
}

func drawCells(text string, left, top int) {
	lines := strings.Split(text, "\n")

	for y, line := range lines {
		for x, char := range line {
			drawCell(left+x, top+y, colorByChar(char))
		}
	}
}

func drawCell(x, y int, color termbox.Attribute) {
	if color != termbox.ColorDefault && color != blankColor {
		if color == colorByChar('K') {
			termbox.SetCell(2*x-1, y, '▓', color, termbox.ColorWhite)
			termbox.SetCell(2*x, y, ' ', color, termbox.ColorWhite)
		} else {
			termbox.SetCell(2*x-1, y, '▓', color, color^termbox.AttrBold)
			termbox.SetCell(2*x, y, ' ', color, color^termbox.AttrBold)
		}
	}
}

func drawBacks(text string, left, top int) {
	lines := strings.Split(text, "\n")

	for y, line := range lines {
		for x, char := range line {
			drawBack(left+x, top+y, colorByChar(char))
		}
	}
}

func drawBack(x, y int, color termbox.Attribute) {
	termbox.SetCell(2*x-1, y, ' ', termbox.ColorDefault, color)
	termbox.SetCell(2*x, y, ' ', termbox.ColorDefault, color)
}

func colorByChar(ch rune) termbox.Attribute {
	return colorMapping[ch]
}

func charByColor(color termbox.Attribute) rune {
	for ch, currentColor := range colorMapping {
		if currentColor == color {
			return ch
		}
	}
	return '.'
}
