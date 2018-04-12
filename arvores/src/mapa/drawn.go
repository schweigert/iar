package mapa

import (
  "os"
  "fmt"
  "image"
  "image/png"
  "image/color"
)

func (m *Mapa) Imprimir() {
  for y := range(m.grid) {
    for x := range(m.grid) {
      fmt.Print(m.grid[y][x])
    }
    fmt.Println()
  }
}

func (m *Mapa) DesenharImagem(nome_arquivo string) {
  img := image.NewRGBA(image.Rect(0, 0, int(m.size_x), int(m.size_y)))

  for y := range(m.grid) {
    for x := range(m.grid[y]) {
      img.Set(x, y, m.PegarCorDaPosicao(x, y))
    }
  }

  f, _ := os.OpenFile(nome_arquivo + ".png", os.O_WRONLY|os.O_CREATE, 0600)
  defer f.Close()

  png.Encode(f, img)
}

func (m *Mapa) Custo(x, y int) uint {
  custos := [4]uint{1, 5, 10, 15}

  return custos[m.grid[y][x]]
}

func (m *Mapa) AtualizarCustos() {}

func (m *Mapa) PegarCorDaPosicao(x, y int) color.RGBA {
  cores := [4]color.RGBA{
    color.RGBA{55, 180, 55, 255},
    color.RGBA{100, 55, 10, 255},
    color.RGBA{0, 120, 255, 255},
    color.RGBA{255, 0, 0, 255},
  }

  return cores[m.grid[y][x]]
}
