package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type Board struct {
	b    [][]int
	i, j int
}

func NewBoard() *Board {
	b := make([][]int, 4)
	for i := range b {
		b[i] = make([]int, 4)
		for j := range b[i] {
			b[i][j] = i*4 + j
		}
	}
	return &Board{b, 3, 3}
}

func (b *Board) move(moves string) {
	for _, m := range moves {
		i, j := b.i, b.j
		switch m {
		case 'L':
			b.b[i][j] = b.b[i][j-1]
			b.b[i][j-1] = 15
			b.i, b.j = i, j-1
		case 'R':
			b.b[i][j] = b.b[i][j+1]
			b.b[i][j+1] = 15
			b.i, b.j = i, j+1
		case 'U':
			b.b[i][j] = b.b[i-1][j]
			b.b[i-1][j] = 15
			b.i, b.j = i-1, j
		case 'D':
			b.b[i][j] = b.b[i+1][j]
			b.b[i+1][j] = 15
			b.i, b.j = i+1, j
		}
	}
}
func (b *Board) isValid() bool {
	for i, row := range b.b {
		for j, num := range row {
			if num != i*4+j {
				return true
			}
		}
	}
	return false
}

type Process []interface{}

func main() {
	var (
		i         uint64
		processes = make([]Process, 0)
	)

	data, err := ioutil.ReadFile("chals/b16_obf.json")
	if err != nil {
		panic(err)
	}
	json.Unmarshal(data, &processes)

	for i = 0; i <= 1<<16-1; i++ {
		b := NewBoard()
		bits := strconv.FormatUint(i, 2)
		if len(bits) < 16 {
			bits = strings.Repeat("0", 16-len(bits)) + bits
		}
		var builder strings.Builder
		for _, process := range processes {
			bit := bits[int(process[0].(float64))]
			if bit == '1' {
				builder.WriteString(process[1].(string))
			} else {
				builder.WriteString(process[2].(string))
			}
		}
		move := builder.String()
		b.move(move)
		if b.isValid() {
			fmt.Println(i, bits)
		}
	}
}
