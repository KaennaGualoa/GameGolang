package game

import (
	"fmt"
	"image/color"
	"my-game/assets"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

type Game struct {
	player           *Player
	lasers           []*Laser
	meteors          []*Meteor
	meteorSpawnTimer *Timer
	score            int
}

func NewGame() *Game {
	g := &Game{
		meteorSpawnTimer: NewTimer(24),
	}
	player := NewPlayer(g)
	g.player = player
	return g
}

// Rode em 60 FPS
// Responsavel por atualizar a logica do jogo
// 60 x por segundos
// um ticke é 1x rodando
func (g *Game) Update() error {
	g.player.Update()

	for _, l := range g.lasers {
		l.Update()
	}

	g.meteorSpawnTimer.Update()
	if g.meteorSpawnTimer.IsReady() {
		g.meteorSpawnTimer.Reset()

		m := NewMeteor()
		g.meteors = append(g.meteors, m)
	}

	for _, m := range g.meteors {
		m.Update()
	}

	for _, m := range g.meteors {
		if m.Collider().Intersects(g.player.Collider()) {
			fmt.Println("Game Over")
			g.Reset()
		}
	}

	for i, m := range g.meteors {

		for j, l := range g.lasers {

			if m.Collider().Intersects(l.Collider()) {
				g.meteors = append(g.meteors[:i], g.meteors[i+1:]...)

				g.lasers = append(g.lasers[:j], g.lasers[j+1:]...)
				g.score += 1

			}
		}
	}

	return nil
}

// Responsavel por desenhar objetos na tela
// lib tbm garante 60x por segundos
func (g *Game) Draw(screen *ebiten.Image) {
	g.player.Draw(screen)

	for _, l := range g.lasers {
		l.Draw(screen)
	}

	for _, m := range g.meteors {
		m.Draw(screen)
	}

	text.Draw(screen, fmt.Sprintf("Pontos: %d", g.score), assets.FontUi, 20, 100, color.White)

}

// Responsavel por retornar o tamanho da tela
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWight, screenHeight
}

func (g *Game) AddLasers(laser *Laser) {
	g.lasers = append(g.lasers, laser)
}

func (g *Game) Reset() {
	g.player = NewPlayer(g)
	g.meteors = nil
	g.lasers = nil
	g.meteorSpawnTimer.Reset()
	g.score = 0
}
