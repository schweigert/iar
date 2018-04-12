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

func (m *Mapa) PegarCorDaPosicao(x, y int) color.RGBA {
  cores_base := [4]color.RGBA{
    color.RGBA{55, 180, 55, 255},
    color.RGBA{100, 55, 10, 255},
    color.RGBA{0, 120, 255, 255},
    color.RGBA{255, 0, 0, 255},
  }

  cores_visited := [4]color.RGBA{
    color.RGBA{55/2, 180/2, 55/2, 255},
    color.RGBA{100/2, 55/2, 10/2, 255},
    color.RGBA{0, 120/2, 255/2, 255},
    color.RGBA{255/2, 0, 0, 255},
  }

  if m.visited[y][x] {
    return cores_visited[m.grid[y][x]]
  }

  return cores_base[m.grid[y][x]]
}
