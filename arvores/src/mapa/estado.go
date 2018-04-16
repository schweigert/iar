package mapa

type Estado struct {
  norte, sul, leste, oeste *Estado
  pos_x, pos_y uint

}

func NewEstado(pos_x, pos_y uint) *Estado {
  return &Estado{
    pos_x: pos_x,
    pos_y: pos_y,
  }
}

func (e *Estado) e_final(e_final *Estado) bool {
  return e.pos_x == e_final.pos_x && e.pos_y == e_final.pos_y
}
