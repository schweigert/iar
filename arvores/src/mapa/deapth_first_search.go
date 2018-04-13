package mapa

import (
  "fmt"
)

func (m *Mapa) BuscaEmProfundidade(start_x, start_y, end_x, end_y int) (uint, *Estado) {
  fmt.Println(start_x, start_y)
  if start_x >= int(m.size_x) || start_y >= int(m.size_y) || start_x < 0 || start_y < 0 {
    return 100, nil
  }

  if m.visited[start_y][start_x] {
    return 100, nil
  }

  m.visited[start_y][start_x] = true

  if start_x == end_x && start_y == end_y {
    return m.Custo(start_x, start_y), NewEstado(uint(start_x), uint(start_y))
  }

  estado := NewEstado(uint(start_x), uint(start_y))

  c_norte, e_norte := m.BuscaEmProfundidade(start_x, start_y + 1, end_x, end_y)
  c_leste, e_leste := m.BuscaEmProfundidade(start_x - 1, start_y, end_x, end_y)
  c_sul, e_sul     := m.BuscaEmProfundidade(start_x, start_y - 1, end_x, end_y)
  c_oeste, e_oeste := m.BuscaEmProfundidade(start_x + 1, start_y, end_x, end_y)

  custos := []uint{
    c_norte,
    c_leste,
    c_sul,
    c_oeste,
  }

  estado.norte = e_norte
  estado.leste = e_leste
  estado.oeste = e_oeste
  estado.sul = e_sul

  return custo_minimo(custos), estado
}

func custo_minimo(values []uint) uint {
  if len(values) == 1 {
    return values[0]
  }

  if len(values) == 2 {
    if values[0] < values[1] {
      return values[0]
    }
    return values[1]
  }

  h, hs := values[0], custo_minimo(values[1:len(values)])

  if h < hs {
    return h
  }

  return hs
}
