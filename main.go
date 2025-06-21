package main

import (
	"image/color"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	ScreenWidth  = 800
	ScreenHeight = 600
	BoardWidth   = 10
	BoardHeight  = 20
	BlockSize    = 30
	BoardOffsetX = (ScreenWidth - BoardWidth*BlockSize) / 2
	BoardOffsetY = (ScreenHeight - BoardHeight*BlockSize) / 2
)

// Game represents the main game state
type Game struct {
	board     [BoardHeight][BoardWidth]int
	piece     *Piece
	nextPiece *Piece
	score     int
	level     int
	lines     int
	gameOver  bool
	dropTime  float64
	lastDrop  float64
}

// Piece represents a Tetris piece
type Piece struct {
	shape     [4][4]int
	x, y      int
	pieceType int
}

// Tetris piece shapes (I, O, T, S, Z, J, L)
var pieceShapes = [][4][4]int{
	// I piece
	{
		{0, 0, 0, 0},
		{1, 1, 1, 1},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
	},
	// O piece
	{
		{0, 0, 0, 0},
		{0, 1, 1, 0},
		{0, 1, 1, 0},
		{0, 0, 0, 0},
	},
	// T piece
	{
		{0, 0, 0, 0},
		{0, 1, 0, 0},
		{1, 1, 1, 0},
		{0, 0, 0, 0},
	},
	// S piece
	{
		{0, 0, 0, 0},
		{0, 1, 1, 0},
		{1, 1, 0, 0},
		{0, 0, 0, 0},
	},
	// Z piece
	{
		{0, 0, 0, 0},
		{1, 1, 0, 0},
		{0, 1, 1, 0},
		{0, 0, 0, 0},
	},
	// J piece
	{
		{0, 0, 0, 0},
		{1, 0, 0, 0},
		{1, 1, 1, 0},
		{0, 0, 0, 0},
	},
	// L piece
	{
		{0, 0, 0, 0},
		{0, 0, 1, 0},
		{1, 1, 1, 0},
		{0, 0, 0, 0},
	},
}

// Piece colors
var pieceColors = []int{1, 2, 3, 4, 5, 6, 7}

func NewGame() *Game {
	rand.Seed(time.Now().UnixNano())
	game := &Game{
		level:    1,
		dropTime: 1.0, // 1 second per drop at level 1
	}
	game.spawnPiece()
	return game
}

func (g *Game) spawnPiece() {
	if g.nextPiece == nil {
		g.nextPiece = g.createRandomPiece()
	}
	g.piece = g.nextPiece
	g.nextPiece = g.createRandomPiece()
	
	// Check if the new piece can be placed
	if !g.isValidPosition(g.piece.x, g.piece.y, g.piece.shape) {
		g.gameOver = true
	}
}

func (g *Game) createRandomPiece() *Piece {
	pieceType := rand.Intn(len(pieceShapes))
	shape := pieceShapes[pieceType]
	return &Piece{
		shape:     shape,
		x:         BoardWidth/2 - 2,
		y:         0,
		pieceType: pieceType,
	}
}

func (g *Game) isValidPosition(x, y int, shape [4][4]int) bool {
	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			if shape[row][col] == 0 {
				continue
			}
			
			newX := x + col
			newY := y + row
			
			if newX < 0 || newX >= BoardWidth || newY >= BoardHeight {
				return false
			}
			
			if newY >= 0 && g.board[newY][newX] != 0 {
				return false
			}
		}
	}
	return true
}

func (g *Game) placePiece() {
	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			if g.piece.shape[row][col] == 0 {
				continue
			}
			boardX := g.piece.x + col
			boardY := g.piece.y + row
			if boardY >= 0 && boardY < BoardHeight && boardX >= 0 && boardX < BoardWidth {
				g.board[boardY][boardX] = g.piece.pieceType + 1
			}
		}
	}
	g.clearLines()
	g.spawnPiece()
}

func (g *Game) clearLines() {
	linesCleared := 0
	
	for row := BoardHeight - 1; row >= 0; row-- {
		fullLine := true
		for col := 0; col < BoardWidth; col++ {
			if g.board[row][col] == 0 {
				fullLine = false
				break
			}
		}
		
		if fullLine {
			// Move all lines above down
			for r := row; r > 0; r-- {
				g.board[r] = g.board[r-1]
			}
			// Clear top line
			for col := 0; col < BoardWidth; col++ {
				g.board[0][col] = 0
			}
			linesCleared++
			row++ // Check the same row again
		}
	}
	
	if linesCleared > 0 {
		g.lines += linesCleared
		g.score += linesCleared * 100 * g.level
		g.level = g.lines/10 + 1
		g.dropTime = 1.0 - float64(g.level-1)*0.1
		if g.dropTime < 0.1 {
			g.dropTime = 0.1
		}
	}
}

func (g *Game) rotatePiece() {
	if g.gameOver {
		return
	}
	
	// Create rotated shape
	rotated := [4][4]int{}
	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			rotated[col][3-row] = g.piece.shape[row][col]
		}
	}
	
	if g.isValidPosition(g.piece.x, g.piece.y, rotated) {
		g.piece.shape = rotated
	}
}

func (g *Game) movePiece(dx, dy int) {
	if g.gameOver {
		return
	}
	
	if g.isValidPosition(g.piece.x+dx, g.piece.y+dy, g.piece.shape) {
		g.piece.x += dx
		g.piece.y += dy
	} else if dy > 0 {
		// If moving down fails, place the piece
		g.placePiece()
	}
}

func (g *Game) hardDrop() {
	if g.gameOver {
		return
	}
	
	for g.isValidPosition(g.piece.x, g.piece.y+1, g.piece.shape) {
		g.piece.y++
		g.score += 2 // Bonus points for hard drop
	}
	g.placePiece()
}

func (g *Game) Update() error {
	if g.gameOver {
		if inpututil.IsKeyJustPressed(ebiten.KeyR) {
			*g = *NewGame()
		}
		return nil
	}
	
	// Handle input
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		g.movePiece(-1, 0)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		g.movePiece(1, 0)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		g.movePiece(0, 1)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		g.rotatePiece()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.hardDrop()
	}
	
	// Auto drop
	g.lastDrop += 1.0 / 60.0 // Assuming 60 FPS
	if g.lastDrop >= g.dropTime {
		g.movePiece(0, 1)
		g.lastDrop = 0
	}
	
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw background
	ebitenutil.DrawRect(screen, 0, 0, ScreenWidth, ScreenHeight, color.Black)
	
	// Draw board background
	boardX := BoardOffsetX - 2
	boardY := BoardOffsetY - 2
	boardWidth := BoardWidth*BlockSize + 4
	boardHeight := BoardHeight*BlockSize + 4
	ebitenutil.DrawRect(screen, float64(boardX), float64(boardY), float64(boardWidth), float64(boardHeight), color.RGBA{R: 0x33, G: 0x33, B: 0x33, A: 0xFF})
	
	// Draw placed blocks
	for row := 0; row < BoardHeight; row++ {
		for col := 0; col < BoardWidth; col++ {
			if g.board[row][col] != 0 {
				x := BoardOffsetX + col*BlockSize
				y := BoardOffsetY + row*BlockSize
				colorVal := getColor(g.board[row][col] - 1)
				ebitenutil.DrawRect(screen, float64(x), float64(y), BlockSize, BlockSize, colorVal)
				ebitenutil.DrawRect(screen, float64(x), float64(y), BlockSize, BlockSize, color.White)
			}
		}
	}
	
	// Draw current piece
	if g.piece != nil {
		for row := 0; row < 4; row++ {
			for col := 0; col < 4; col++ {
				if g.piece.shape[row][col] != 0 {
					x := BoardOffsetX + (g.piece.x+col)*BlockSize
					y := BoardOffsetY + (g.piece.y+row)*BlockSize
					if y >= BoardOffsetY {
						colorVal := getColor(g.piece.pieceType)
						ebitenutil.DrawRect(screen, float64(x), float64(y), BlockSize, BlockSize, colorVal)
						ebitenutil.DrawRect(screen, float64(x), float64(y), BlockSize, BlockSize, color.White)
					}
				}
			}
		}
	}
	
	// Draw UI
	ebitenutil.DebugPrint(screen, "Score: "+strconv.Itoa(g.score))
	ebitenutil.DebugPrintAt(screen, "Level: "+strconv.Itoa(g.level), 0, 20)
	ebitenutil.DebugPrintAt(screen, "Lines: "+strconv.Itoa(g.lines), 0, 40)
	
	// Draw controls
	controls := []string{
		"Controls:",
		"Arrow Keys: Move",
		"Up: Rotate",
		"Space: Hard Drop",
		"R: Restart (when game over)",
	}
	
	for i, control := range controls {
		ebitenutil.DebugPrintAt(screen, control, ScreenWidth-200, 20+i*20)
	}
	
	if g.gameOver {
		ebitenutil.DebugPrintAt(screen, "GAME OVER", ScreenWidth/2-50, ScreenHeight/2)
		ebitenutil.DebugPrintAt(screen, "Press R to restart", ScreenWidth/2-70, ScreenHeight/2+20)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func getColor(pieceType int) color.Color {
	colors := []color.RGBA{
		{R: 0x00, G: 0xFF, B: 0xFF, A: 0xFF}, // I - Cyan
		{R: 0xFF, G: 0xFF, B: 0x00, A: 0xFF}, // O - Yellow
		{R: 0x80, G: 0x00, B: 0x80, A: 0xFF}, // T - Purple
		{R: 0x00, G: 0xFF, B: 0x00, A: 0xFF}, // S - Green
		{R: 0xFF, G: 0x00, B: 0x00, A: 0xFF}, // Z - Red
		{R: 0x00, G: 0x00, B: 0xFF, A: 0xFF}, // J - Blue
		{R: 0xFF, G: 0x80, B: 0x00, A: 0xFF}, // L - Orange
	}
	if pieceType >= 0 && pieceType < len(colors) {
		return colors[pieceType]
	}
	return color.White
}

func main() {
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Xtris Clone")
	
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
} 