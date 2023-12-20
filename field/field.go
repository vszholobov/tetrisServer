package field

import (
	"fmt"
	"math/big"
	"math/rand"
	"strconv"
	"strings"
)

const FieldWidth = 12
const FieldHeight = 21
const moveToTopASCII = "\033[22A"
const moveRightASCII = "\r\033[36C"
const moveDownOneLineASCII = "\r\033[1B"
const moveDownAllLinesASCII = "\r\033[17B"

var fullLine, _ = big.NewInt(0).SetString("111111111111", 2)
var emptyLine, _ = big.NewInt(0).SetString("100000000001", 2)

// ############
// #          #
// #          #
// #          #
// #          #
// ############

type Field struct {
	Val          *big.Int
	CurrentPiece *Piece
	NextPiece    *Piece
	Score        *int
	CleanCount   *int
}

func MakeDefaultField() Field {
	fieldVal, _ := big.NewInt(0).SetString(
		"111111111111"+
			"100000000001"+
			"100000000001"+
			"100000000001"+
			"100000000001"+
			"100000000001"+
			"100000000001"+
			"100000000001"+
			"100000000001"+
			"100000000001"+
			"100000000001"+
			"100000000001"+
			"100000000001"+
			"100000000001"+
			"100000000001"+
			"100000000001"+
			"100000000001"+
			"100000000001"+
			"100000000001"+
			"100000000001"+
			"100000000001", 2)
	return MakeField(fieldVal)
}

func MakeField(fieldVal *big.Int) Field {
	score := 0
	cleanCount := 0
	gameField := Field{Val: fieldVal, Score: &score, CleanCount: &cleanCount}
	gameField.CurrentPiece = SelectRandomPiece(&gameField)
	gameField.NextPiece = SelectRandomPiece(&gameField)
	return gameField
}

func (gameField *Field) SelectNextPiece() {
	gameField.CurrentPiece = gameField.NextPiece
	gameField.NextPiece = SelectRandomPiece(gameField)
}

func (gameField *Field) String() string {
	newField := big.NewInt(0).Set(gameField.Val)
	newShape := big.NewInt(0).Set(gameField.CurrentPiece.GetVal())
	newField.Or(newField, newShape)
	return fmt.Sprintf("%b", newField)
}

func (gameField *Field) Clean() {
	restField := big.NewInt(0)
	currentCleanCount := 0
	for i := 0; i < FieldHeight-1; i++ {
		curRange := uint(i * FieldWidth)
		lineMask := big.NewInt(0).Lsh(fullLine, curRange)
		lineIsFilled := big.NewInt(0).And(lineMask, gameField.Val).Cmp(lineMask) == 0

		if lineIsFilled {
			// add empy line to end of field
			restField.Lsh(restField, FieldWidth)
			restField.Or(restField, emptyLine)
			currentCleanCount += 1
		} else {
			// add current line to start of field
			lineMask.And(lineMask, gameField.Val)
			restField.Or(lineMask, restField)
		}
	}
	*gameField.CleanCount += currentCleanCount
	*gameField.Score += currentCleanCount * gameField.GetSpeed() * 10 / (5 - currentCleanCount)
	// 22 lines. One redundant line for correct or concatenation.
	// So shift to the right by the length of the field after concatenation to remove redundant empty line
	gameField.Val.SetString(
		"111111111111"+
			"000000000000"+
			"000000000000"+
			"000000000000"+
			"000000000000"+
			"000000000000"+
			"000000000000"+
			"000000000000"+
			"000000000000"+
			"000000000000"+
			"000000000000"+
			"000000000000"+
			"000000000000"+
			"000000000000"+
			"000000000000"+
			"000000000000"+
			"000000000000"+
			"000000000000"+
			"000000000000"+
			"000000000000"+
			"000000000000", 2)
	gameField.Val.Or(gameField.Val, restField)
}

func (gameField *Field) Intersects(pieceVal *big.Int) bool {
	newField := big.NewInt(0).Set(gameField.Val)
	newShape := big.NewInt(0).Set(pieceVal)
	return newField.And(newField, newShape).Cmp(big.NewInt(0)) != 0
}

func CopyBigInt(val *big.Int) *big.Int {
	return big.NewInt(0).Set(val)
}

var builder = strings.Builder{}

func PrintField(field *Field) {
	builder.Reset()
	builder.WriteString(moveToTopASCII)
	fieldStr := field.String()
	for i := 20; i >= 0; i-- {
		line := fieldStr[i*12 : i*12+12]
		line = strings.ReplaceAll(line, "1", " Ж ")
		line = strings.ReplaceAll(line, "0", "   ")
		builder.WriteString(line)
		builder.WriteString("\n")
	}
	builder.WriteString("Score: ")
	builder.WriteString(strconv.Itoa(*field.Score))
	builder.WriteString(" | Speed: ")
	builder.WriteString(strconv.Itoa(field.GetSpeed()))
	builder.WriteString(" | Cleaned: ")
	builder.WriteString(strconv.Itoa(*field.CleanCount))
	fmt.Println(builder.String())
	printNextPiece(field.NextPiece)
}

func printNextPiece(nextPiece *Piece) {
	//pieceVal := nextPiece.GetVal()

	fmt.Print(moveToTopASCII + moveRightASCII + " ##############")
	fmt.Printf(moveDownOneLineASCII + moveRightASCII + " #            #")
	pieceLines := RepresentationByType[nextPiece.pieceType]
	for i := 0; i < 2; i++ {
		curLine := "            "
		if i < len(pieceLines) {
			curLine = pieceLines[i]
		}
		fmt.Printf(moveDownOneLineASCII+moveRightASCII+" #%s#", curLine)
		//curLine := big.NewInt(0).Lsh(fullLine, uint(i)*FieldWidth)
		//checkCurrLine := big.NewInt(0).And(curLine, nextPiece.GetVal())
		//line := fmt.Sprintf("%10b", checkCurrLine)
		//line = strings.ReplaceAll(line, "1", "Ж")
		//line = strings.ReplaceAll(line, "0", "")
		//fmt.Print(moveDownOneLineASCII + moveRightASCII + " #          #")
	}
	fmt.Printf(moveDownOneLineASCII + moveRightASCII + " #            #")
	fmt.Print(moveDownOneLineASCII + moveRightASCII + " ##############")
	fmt.Print(moveDownAllLinesASCII)

	//fmt.Print(moveDownOneLineASCII + moveRightASCII + " #          #")
	//fmt.Print(moveDownOneLineASCII + moveRightASCII + " #          #")
	//fmt.Print(moveDownOneLineASCII + moveRightASCII + " #          #")
}

func (gameField *Field) GetSpeed() int {
	return *gameField.CleanCount/4 + 1
}

func SelectRandomPiece(gameField *Field) *Piece {
	pieceTypeRnd := rand.Intn(7)
	var pieceType PieceType
	if pieceTypeRnd == 0 {
		pieceType = IShape
	} else if pieceTypeRnd == 1 {
		pieceType = RightLShape
	} else if pieceTypeRnd == 2 {
		pieceType = TShape
	} else if pieceTypeRnd == 3 {
		pieceType = ZigZagRight
	} else if pieceTypeRnd == 4 {
		pieceType = ZigZagLeft
	} else if pieceTypeRnd == 5 {
		pieceType = SquareShape
	} else if pieceTypeRnd == 6 {
		pieceType = LeftLShape
	}
	piece := MakePiece(gameField, pieceType)
	return &piece
}
