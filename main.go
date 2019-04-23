package main

import (
	"fmt"
)

type Item struct {
	title string
	t     int
}
type Game struct {
	emptyIndex int
	size       int
	items      []Item
}

func newGame(num int) *Game {
	g := &Game{
		emptyIndex: num,
		size:       num,
		items:      make([]Item, num*2+1, num*2+1),
	}
	for i := 0; i < num; i++ {
		g.items[i].t = 1
		g.items[i].title = fmt.Sprintf("%c", '1'+num-i-1)
	}
	for i := num + 1; i < num*2+1; i++ {
		g.items[i].t = -1
		g.items[i].title = fmt.Sprintf("%c", 'A'+i-num-1)
	}
	return g
}

func isOK(g *Game) bool {
	if g.emptyIndex != g.size {
		return false
	}
	for i := 0; i < g.size; i++ {
		if g.items[i].title != fmt.Sprintf("%c", 'A'+i) {
			return false
		}
	}
	for i := g.size + 1; i < g.size*2+1; i++ {
		if g.items[i].title != fmt.Sprintf("%c", '1'+g.size*2-i) {
			return false
		}

	}
	return true
}
func search(g *Game, ans []string, indexAns int) bool {
	if isOK(g) {
		return true
	}
	index := g.emptyIndex
	for i := index - 2; i < index; i++ {
		if i >= 0 && g.items[i].t == 1 {
			tmp := g.items[i]
			g.items[index] = g.items[i]
			g.emptyIndex = i
			ans[indexAns] = g.items[i].title
			if search(g, ans, indexAns+1) {
				return true
			}
			g.emptyIndex = index
			g.items[i] = tmp
		}
	}
	for i := index + 1; i < index+3; i++ {
		if i < g.size*2+1 && g.items[i].t == -1 {
			tmp := g.items[i]
			g.items[index] = g.items[i]
			g.emptyIndex = i
			ans[indexAns] = g.items[i].title
			if search(g, ans, indexAns+1) {
				return true
			}
			g.emptyIndex = index
			g.items[i] = tmp
		}
	}
	return false
}
func main() {
	var num int
	fmt.Printf("input num:")
	fmt.Scanf("%d", &num)
	g := newGame(num)
	ans := make([]string, 100)
	search(g, ans, 0)

	for i := 0; i < 100; i++ {
		if ans[i] != "" {
			fmt.Printf("%s ", ans[i])
		}
	}
	fmt.Println()
}
