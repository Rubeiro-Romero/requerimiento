package modelos

type RespuestaDTO[T any] struct {
	Data    Cancion
	Codigo  int32
	Mensaje string
}
