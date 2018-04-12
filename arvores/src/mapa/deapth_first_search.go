package mapa

import (
  "math"
  "strconv"
)

var img_n int = 0


func (m *Mapa) BuscaEmProfundidade(start_x, start_y, end_x, end_y int) float64 {
  img_n ++

  if start_x == end_x && start_y == end_y {
    return float64(m.Custo(start_x, start_y))
  }

  if start_x >= int(m.size_x) || start_y >= int(m.size_y) || start_x < 0 || start_y < 0 {
    return math.NaN()
  }

  if m.visited[start_y][start_x] {
    return math.NaN()
  }

  m.visited[start_y][start_x] = true

  if img_n % 100 == 0 {
    m.DesenharImagem(strconv.Itoa(img_n))
  }

  custo_norte := m.BuscaEmProfundidade(start_x, start_y + 1, end_x, end_y)
  custo_leste := m.BuscaEmProfundidade(start_x - 1, start_y, end_x, end_y)
  custo_sul   := m.BuscaEmProfundidade(start_x, start_y - 1, end_x, end_y)
  custo_oeste := m.BuscaEmProfundidade(start_x + 1, start_y, end_x, end_y)

  if !math.IsNaN(custo_norte) {
    if custo_norte <= custo_leste && custo_norte <= custo_sul && custo_norte <= custo_oeste {
      return custo_norte
    }
  }

  if !math.IsNaN(custo_leste) {
    if custo_leste <= custo_norte && custo_leste <= custo_sul && custo_leste <= custo_oeste {
      return custo_leste
    }
  }

  if !math.IsNaN(custo_sul) {
    if custo_sul <= custo_norte && custo_sul <= custo_leste && custo_sul <= custo_oeste {
      return custo_sul
    }
  }

  return custo_oeste
}
