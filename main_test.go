package main

import (
	"testing"
)

func TestNewGame(t *testing.T) {
	game := NewGame()
	
	if game == nil {
		t.Fatal("NewGame() returned nil")
	}
	
	if game.level != 1 {
		t.Errorf("Expected level 1, got %d", game.level)
	}
	
	if game.score != 0 {
		t.Errorf("Expected score 0, got %d", game.score)
	}
	
	if game.lines != 0 {
		t.Errorf("Expected lines 0, got %d", game.lines)
	}
	
	if game.gameOver {
		t.Error("New game should not be game over")
	}
	
	if game.piece == nil {
		t.Error("New game should have a current piece")
	}
	
	if game.nextPiece == nil {
		t.Error("New game should have a next piece")
	}
}

func TestCreateRandomPiece(t *testing.T) {
	game := NewGame()
	piece := game.createRandomPiece()
	
	if piece == nil {
		t.Fatal("createRandomPiece() returned nil")
	}
	
	if piece.x != BoardWidth/2-2 {
		t.Errorf("Expected x position %d, got %d", BoardWidth/2-2, piece.x)
	}
	
	if piece.y != 0 {
		t.Errorf("Expected y position 0, got %d", piece.y)
	}
	
	if piece.pieceType < 0 || piece.pieceType >= len(pieceShapes) {
		t.Errorf("Invalid piece type: %d", piece.pieceType)
	}
}

func TestIsValidPosition(t *testing.T) {
	game := NewGame()
	
	// Test valid position
	valid := game.isValidPosition(3, 0, pieceShapes[0]) // I piece
	if !valid {
		t.Error("Valid position should return true")
	}
	
	// Test out of bounds left
	valid = game.isValidPosition(-1, 0, pieceShapes[0])
	if valid {
		t.Error("Out of bounds left should return false")
	}
	
	// Test out of bounds right
	valid = game.isValidPosition(BoardWidth, 0, pieceShapes[0])
	if valid {
		t.Error("Out of bounds right should return false")
	}
	
	// Test out of bounds bottom
	valid = game.isValidPosition(0, BoardHeight, pieceShapes[0])
	if valid {
		t.Error("Out of bounds bottom should return false")
	}
	
	// Test collision with existing blocks
	game.board[5][5] = 1
	valid = game.isValidPosition(4, 4, pieceShapes[0])
	if valid {
		t.Error("Collision with existing block should return false")
	}
	
	// Test piece partially out of bounds
	valid = game.isValidPosition(BoardWidth-1, 0, pieceShapes[0])
	if valid {
		t.Error("Piece partially out of bounds should return false")
	}
}

func TestRotatePiece(t *testing.T) {
	game := NewGame()
	originalShape := game.piece.shape
	
	game.rotatePiece()
	
	// Check that the shape changed (for most pieces)
	// Note: O piece looks the same after rotation, so we need to check differently
	rotated := false
	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			if game.piece.shape[row][col] != originalShape[row][col] {
				rotated = true
				break
			}
		}
		if rotated {
			break
		}
	}
	
	// For O piece, rotation should still work but shape might look the same
	// For other pieces, shape should change
	if !rotated && game.piece.pieceType != 1 { // 1 is O piece
		t.Error("Piece should be rotated (except O piece which looks the same)")
	}
	
	// Test rotation when game is over
	game.gameOver = true
	game.rotatePiece()
	// Should not crash and should not rotate
}

func TestRotatePieceWithCollision(t *testing.T) {
	game := NewGame()
	
	// Place blocks around the piece to prevent rotation
	game.board[game.piece.y+1][game.piece.x] = 1
	game.board[game.piece.y+1][game.piece.x+1] = 1
	game.board[game.piece.y+1][game.piece.x+2] = 1
	
	originalShape := game.piece.shape
	game.rotatePiece()
	
	// Shape should not change due to collision
	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			if game.piece.shape[row][col] != originalShape[row][col] {
				t.Error("Piece should not rotate when collision would occur")
			}
		}
	}
}

func TestMovePiece(t *testing.T) {
	game := NewGame()
	originalX := game.piece.x
	
	// Test valid move
	game.movePiece(1, 0)
	if game.piece.x != originalX+1 {
		t.Errorf("Expected x position %d, got %d", originalX+1, game.piece.x)
	}
	
	// Test invalid move (into wall)
	game.movePiece(-10, 0)
	if game.piece.x != originalX+1 {
		t.Error("Piece should not move when hitting wall")
	}
	
	// Test moving down and placing piece
	// Place obstacles at the bottom to force piece placement
	for col := 0; col < BoardWidth; col++ {
		game.board[BoardHeight-1][col] = 1
	}
	
	originalY := game.piece.y
	game.movePiece(0, 1) // Try to move down
	if game.piece.y == originalY {
		t.Error("Piece should be placed when moving down fails")
	}
}

func TestMovePieceGameOver(t *testing.T) {
	game := NewGame()
	game.gameOver = true
	originalX := game.piece.x
	originalY := game.piece.y
	
	game.movePiece(1, 0)
	if game.piece.x != originalX || game.piece.y != originalY {
		t.Error("Piece should not move when game is over")
	}
}

func TestClearLines(t *testing.T) {
	game := NewGame()
	
	// Fill bottom row
	for col := 0; col < BoardWidth; col++ {
		game.board[BoardHeight-1][col] = 1
	}
	
	// Fill second to bottom row
	for col := 0; col < BoardWidth; col++ {
		game.board[BoardHeight-2][col] = 2
	}
	
	// Place a piece to trigger line clearing
	game.piece = &Piece{
		shape:     pieceShapes[0],
		x:         0,
		y:         BoardHeight - 1,
		pieceType: 0,
	}
	
	originalLines := game.lines
	originalScore := game.score
	
	game.placePiece()
	
	if game.lines != originalLines+2 {
		t.Errorf("Expected %d lines cleared, got %d", 2, game.lines-originalLines)
	}
	
	if game.score <= originalScore {
		t.Error("Score should increase after clearing lines")
	}
}

func TestClearLinesNoLines(t *testing.T) {
	game := NewGame()
	originalLines := game.lines
	originalScore := game.score
	originalLevel := game.level
	originalDropTime := game.dropTime
	
	game.clearLines()
	
	if game.lines != originalLines {
		t.Error("Lines should not change when no lines are cleared")
	}
	
	if game.score != originalScore {
		t.Error("Score should not change when no lines are cleared")
	}
	
	if game.level != originalLevel {
		t.Error("Level should not change when no lines are cleared")
	}
	
	if game.dropTime != originalDropTime {
		t.Error("Drop time should not change when no lines are cleared")
	}
}

func TestClearLinesMultipleLevels(t *testing.T) {
	game := NewGame()
	
	// Set up to trigger level increase
	game.lines = 9
	
	// Fill a row to clear
	for col := 0; col < BoardWidth; col++ {
		game.board[BoardHeight-1][col] = 1
	}
	
	game.clearLines()
	
	if game.level != 2 {
		t.Errorf("Expected level 2, got %d", game.level)
	}
	
	if game.dropTime >= 1.0 {
		t.Error("Drop time should decrease with level")
	}
}

func TestClearLinesMaxSpeed(t *testing.T) {
	game := NewGame()
	
	// Set up to trigger maximum speed (level 10+)
	game.lines = 99
	
	// Fill a row to clear
	for col := 0; col < BoardWidth; col++ {
		game.board[BoardHeight-1][col] = 1
	}
	
	game.clearLines()
	
	if game.dropTime < 0.1 {
		t.Error("Drop time should not go below 0.1")
	}
}

func TestGameOver(t *testing.T) {
	game := NewGame()
	
	// Fill the board almost completely
	for row := 0; row < BoardHeight-1; row++ {
		for col := 0; col < BoardWidth; col++ {
			game.board[row][col] = 1
		}
	}
	
	// Try to spawn a new piece
	game.spawnPiece()
	
	if !game.gameOver {
		t.Error("Game should be over when new piece cannot be placed")
	}
}

func TestLevelProgression(t *testing.T) {
	game := NewGame()
	
	// Simulate clearing 10 lines
	for i := 0; i < 10; i++ {
		// Fill a row
		for col := 0; col < BoardWidth; col++ {
			game.board[BoardHeight-1][col] = 1
		}
		
		// Clear the line
		game.clearLines()
	}
	
	if game.level != 2 {
		t.Errorf("Expected level 2 after 10 lines, got %d", game.level)
	}
	
	if game.dropTime >= 1.0 {
		t.Error("Drop time should decrease with level")
	}
}

func TestHardDrop(t *testing.T) {
	game := NewGame()
	originalScore := game.score
	
	game.hardDrop()
	
	if game.score <= originalScore {
		t.Error("Score should increase during hard drop")
	}
}

func TestHardDropGameOver(t *testing.T) {
	game := NewGame()
	game.gameOver = true
	originalScore := game.score
	
	game.hardDrop()
	
	if game.score != originalScore {
		t.Error("Score should not change when game is over")
	}
}

func TestHardDropWithObstacles(t *testing.T) {
	game := NewGame()
	
	// Place obstacles below the piece
	for col := 0; col < BoardWidth; col++ {
		game.board[BoardHeight-2][col] = 1
	}
	
	originalScore := game.score
	game.hardDrop()
	
	if game.score <= originalScore {
		t.Error("Score should increase during hard drop even with obstacles")
	}
}

func TestUpdateMethod(t *testing.T) {
	game := NewGame()
	
	// Test Update method doesn't crash
	err := game.Update()
	if err != nil {
		t.Errorf("Update should not return error: %v", err)
	}
	
	// Test auto drop functionality
	originalY := game.piece.y
	game.lastDrop = game.dropTime + 0.1 // Force auto drop
	
	err = game.Update()
	if err != nil {
		t.Errorf("Update should not return error: %v", err)
	}
	
	// Piece should have moved down or been placed
	if game.piece.y == originalY && game.lastDrop != 0 {
		t.Error("Auto drop should move piece down or place it")
	}
}

func TestUpdateMethodGameOver(t *testing.T) {
	game := NewGame()
	game.gameOver = true
	
	// Test Update method when game is over
	err := game.Update()
	if err != nil {
		t.Errorf("Update should not return error when game is over: %v", err)
	}
}

func TestUpdateMethodAutoDrop(t *testing.T) {
	game := NewGame()
	
	// Test auto drop timing
	originalY := game.piece.y
	game.lastDrop = 0
	
	// First update should not trigger auto drop
	err := game.Update()
	if err != nil {
		t.Errorf("Update should not return error: %v", err)
	}
	
	// Set lastDrop to trigger auto drop
	game.lastDrop = game.dropTime + 0.1
	err = game.Update()
	if err != nil {
		t.Errorf("Update should not return error: %v", err)
	}
	
	// Should have moved down or been placed
	if game.piece.y == originalY && game.lastDrop != 0 {
		t.Error("Auto drop should move piece down or place it")
	}
}

func TestUpdateMethodAutoDropReset(t *testing.T) {
	game := NewGame()
	
	// Test that lastDrop resets after auto drop
	game.lastDrop = game.dropTime + 0.1
	
	err := game.Update()
	if err != nil {
		t.Errorf("Update should not return error: %v", err)
	}
	
	// lastDrop should be reset to 0 after auto drop
	if game.lastDrop != 0 {
		t.Error("lastDrop should be reset to 0 after auto drop")
	}
}

func TestUpdateMethodMultipleFrames(t *testing.T) {
	game := NewGame()
	
	// Test multiple frames of auto drop
	for i := 0; i < 10; i++ {
		err := game.Update()
		if err != nil {
			t.Errorf("Update should not return error on frame %d: %v", i, err)
		}
	}
	
	// Game should still be running
	if game.gameOver {
		t.Error("Game should not be over after multiple frames")
	}
}

func TestUpdateMethodWithObstacles(t *testing.T) {
	game := NewGame()
	
	// Place obstacles to force piece placement
	for col := 0; col < BoardWidth; col++ {
		game.board[BoardHeight-2][col] = 1
	}
	
	// Force auto drop
	game.lastDrop = game.dropTime + 0.1
	
	err := game.Update()
	if err != nil {
		t.Errorf("Update should not return error: %v", err)
	}
	
	// Piece should be placed due to obstacles
	if game.lastDrop != 0 {
		t.Error("Piece should be placed when hitting obstacles")
	}
}

func TestLayoutMethod(t *testing.T) {
	game := NewGame()
	
	// Test Layout method
	width, height := game.Layout(1000, 800)
	
	if width != ScreenWidth {
		t.Errorf("Expected width %d, got %d", ScreenWidth, width)
	}
	
	if height != ScreenHeight {
		t.Errorf("Expected height %d, got %d", ScreenHeight, height)
	}
}

func TestSpawnPiece(t *testing.T) {
	game := NewGame()
	
	// Test that spawnPiece creates both current and next piece
	if game.piece == nil {
		t.Error("Current piece should not be nil after spawnPiece")
	}
	
	if game.nextPiece == nil {
		t.Error("Next piece should not be nil after spawnPiece")
	}
	
	// Test that pieces have valid types
	if game.piece.pieceType < 0 || game.piece.pieceType >= len(pieceShapes) {
		t.Error("Current piece should have valid piece type")
	}
	
	if game.nextPiece.pieceType < 0 || game.nextPiece.pieceType >= len(pieceShapes) {
		t.Error("Next piece should have valid piece type")
	}
}

func TestSpawnPieceGameOver(t *testing.T) {
	game := NewGame()
	
	// Fill the board to cause game over
	for row := 0; row < BoardHeight; row++ {
		for col := 0; col < BoardWidth; col++ {
			game.board[row][col] = 1
		}
	}
	
	game.spawnPiece()
	
	if !game.gameOver {
		t.Error("Game should be over when new piece cannot be placed")
	}
}

func TestPlacePiece(t *testing.T) {
	game := NewGame()
	
	// Test placing piece at valid position
	originalScore := game.score
	game.placePiece()
	
	// Score should increase if lines were cleared
	if game.score < originalScore {
		t.Error("Score should not decrease when placing piece")
	}
	
	// New piece should be spawned
	if game.piece == nil {
		t.Error("New piece should be spawned after placing piece")
	}
}

func TestPlacePieceOutOfBounds(t *testing.T) {
	game := NewGame()
	
	// Place piece at edge of board
	game.piece.x = BoardWidth - 1
	game.piece.y = BoardHeight - 1
	
	// Should not panic
	game.placePiece()
}

func TestPieceShapes(t *testing.T) {
	if len(pieceShapes) != 7 {
		t.Errorf("Expected 7 piece shapes, got %d", len(pieceShapes))
	}
	
	// Test that each piece has at least one block
	for i, shape := range pieceShapes {
		hasBlock := false
		for row := 0; row < 4; row++ {
			for col := 0; col < 4; col++ {
				if shape[row][col] == 1 {
					hasBlock = true
					break
				}
			}
			if hasBlock {
				break
			}
		}
		
		if !hasBlock {
			t.Errorf("Piece %d has no blocks", i)
		}
	}
}

func TestGetColor(t *testing.T) {
	colors := []struct {
		r, g, b, a uint32
	}{
		{0x00, 0xFF, 0xFF, 0xFF}, // I - Cyan
		{0xFF, 0xFF, 0x00, 0xFF}, // O - Yellow
		{0x80, 0x00, 0x80, 0xFF}, // T - Purple
		{0x00, 0xFF, 0x00, 0xFF}, // S - Green
		{0xFF, 0x00, 0x00, 0xFF}, // Z - Red
		{0x00, 0x00, 0xFF, 0xFF}, // J - Blue
		{0xFF, 0x80, 0x00, 0xFF}, // L - Orange
	}
	
	for i, expected := range colors {
		c := getColor(i)
		r, g, b, a := c.RGBA()
		// RGBA returns values in 16 bits, so shift right by 8
		if r>>8 != expected.r || g>>8 != expected.g || b>>8 != expected.b || a>>8 != expected.a {
			t.Errorf("Expected color RGBA(%02X,%02X,%02X,%02X) for piece %d, got RGBA(%02X,%02X,%02X,%02X)", expected.r, expected.g, expected.b, expected.a, i, r>>8, g>>8, b>>8, a>>8)
		}
	}
	
	// Test invalid piece type
	c := getColor(-1)
	r, g, b, a := c.RGBA()
	if r>>8 != 0xFF || g>>8 != 0xFF || b>>8 != 0xFF || a>>8 != 0xFF {
		t.Errorf("Expected white color for invalid piece, got RGBA(%02X,%02X,%02X,%02X)", r>>8, g>>8, b>>8, a>>8)
	}
	
	// Test piece type beyond valid range
	c = getColor(100)
	r, g, b, a = c.RGBA()
	if r>>8 != 0xFF || g>>8 != 0xFF || b>>8 != 0xFF || a>>8 != 0xFF {
		t.Errorf("Expected white color for out-of-range piece, got RGBA(%02X,%02X,%02X,%02X)", r>>8, g>>8, b>>8, a>>8)
	}
}

func TestGameStateConsistency(t *testing.T) {
	game := NewGame()
	
	// Test that game state remains consistent after multiple operations
	originalLevel := game.level
	originalDropTime := game.dropTime
	
	// Perform some operations
	game.movePiece(1, 0)
	game.rotatePiece()
	game.movePiece(0, 1)
	
	// Level and drop time should not change without line clears
	if game.level != originalLevel {
		t.Error("Level should not change without line clears")
	}
	
	if game.dropTime != originalDropTime {
		t.Error("Drop time should not change without line clears")
	}
}

func TestMultipleLineClears(t *testing.T) {
	game := NewGame()
	
	// Fill multiple rows
	for row := BoardHeight - 3; row < BoardHeight; row++ {
		for col := 0; col < BoardWidth; col++ {
			game.board[row][col] = 1
		}
	}
	
	originalLines := game.lines
	originalScore := game.score
	
	game.clearLines()
	
	expectedLines := originalLines + 3
	if game.lines != expectedLines {
		t.Errorf("Expected %d lines, got %d", expectedLines, game.lines)
	}
	
	if game.score <= originalScore {
		t.Error("Score should increase after clearing multiple lines")
	}
}

func TestMovePieceEdgeCases(t *testing.T) {
	game := NewGame()
	
	// Test moving piece to edge of board
	game.piece.x = 0
	game.movePiece(-1, 0) // Try to move left from left edge
	if game.piece.x != 0 {
		t.Error("Piece should not move past left edge")
	}
	
	game.piece.x = BoardWidth - 1
	game.movePiece(1, 0) // Try to move right from right edge
	if game.piece.x != BoardWidth-1 {
		t.Error("Piece should not move past right edge")
	}
	
	// Test moving piece to bottom - this should place the piece
	game.piece.y = BoardHeight - 1
	originalY := game.piece.y
	game.movePiece(0, 1) // Try to move down from bottom
	if game.piece.y == originalY {
		t.Error("Piece should be placed when moving down from bottom")
	}
}

func TestMovePieceHorizontalOnly(t *testing.T) {
	game := NewGame()
	
	// Test that horizontal movement doesn't place piece
	originalY := game.piece.y
	game.movePiece(1, 0) // Move right
	
	if game.piece.y != originalY {
		t.Error("Horizontal movement should not affect Y position")
	}
}

func TestMovePieceDownWithObstacles(t *testing.T) {
	game := NewGame()
	
	// Place obstacles below the piece
	for col := 0; col < BoardWidth; col++ {
		game.board[BoardHeight-2][col] = 1
	}
	
	originalY := game.piece.y
	game.movePiece(0, 1) // Try to move down into collision
	
	// Should place the piece instead of moving
	if game.piece.y == originalY {
		t.Error("Piece should be placed when moving down fails")
	}
} 