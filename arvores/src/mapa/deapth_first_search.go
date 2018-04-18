package mapa

import (
  "strconv"
)

var frame_profundidade int

const (
  INITIAL_FRAME = 1000000
)

func (m *Mapa) BuscaEmProfundidade(start_x, start_y, end_x, end_y int) (int, *Estado) {
  if start_x >= int(m.size_x) || start_y >= int(m.size_y) || start_x < 0 || start_y < 0 {
    return -1, nil
  }

  if m.visited[start_y][start_x] {
    return -1, nil
  }

  m.visited[start_y][start_x] = true

  if start_x == end_x && start_y == end_y {
    return 0, NewEstado(uint(start_x), uint(start_y))
  }

  estado := NewEstado(uint(start_x), uint(start_y))

  if frame_profundidade % 100 == 0 {
    m.DesenharImagem(strconv.Itoa(INITIAL_FRAME + frame_profundidade))
  }
  frame_profundidade++

  c_norte, e_norte := m.BuscaEmProfundidade(start_x, start_y - 1, end_x, end_y)
  estado.norte = e_norte

  if c_norte != -1 {
    return m.Custo(start_x, start_y) + c_norte, estado
  }

  c_leste, e_leste := m.BuscaEmProfundidade(start_x + 1, start_y, end_x, end_y)
  estado.leste = e_leste

  if c_leste != -1 {
    return m.Custo(start_x, start_y) + c_leste, estado
  }

  c_sul, e_sul := m.BuscaEmProfundidade(start_x, start_y + 1, end_x, end_y)
  estado.sul = e_sul

  if c_sul != -1 {
    return m.Custo(start_x, start_y) + c_sul, estado
  }

  c_oeste, e_oeste := m.BuscaEmProfundidade(start_x - 1, start_y, end_x, end_y)
  estado.oeste = e_oeste

  return m.Custo(start_x, start_y) + c_oeste, estado
}
