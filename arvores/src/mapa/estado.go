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
