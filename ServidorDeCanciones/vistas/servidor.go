package main

import (
	"context"
	"log"
	"net"

	"ServidorDeCanciones.local/grpc-ServidorDeCanciones/modelos"
	"ServidorDeCanciones.local/grpc-ServidorDeCanciones/servicios"
	pb "ServidorDeCanciones.local/grpc-ServidorDeCanciones/serviciosCancion"
	"google.golang.org/grpc"
)

var vectorCanciones = make([]modelos.Cancion, 3)

// Estructura que implementa el servicio
type servidorCanciones struct {
	pb.UnimplementedServiciosCancionesServer
}

// Implementación del método buscarCancion
func (s *servidorCanciones) BuscarCancion(ctx context.Context, req *pb.PeticionDTO) (*pb.RespuestaCancionDTO, error) {

	titulo := req.GetTitulo()
	resp := servicios.BuscarCancion(titulo, vectorCanciones)

	var respuesta pb.RespuestaCancionDTO
	respuesta.Codigo = resp.Codigo
	respuesta.Mensaje = resp.Mensaje

	if resp.Codigo == 200 {
		respuesta.ObjCancion = new(pb.Cancion)
		respuesta.ObjCancion.Titulo = resp.Data.Titulo
		respuesta.ObjCancion.Tamanio = resp.Data.Tamanio
		respuesta.ObjCancion.Url = resp.Data.Url
		respuesta.ObjCancion.EsActivada = resp.Data.EsActivada
	}

	return &respuesta, nil
}

func main() {
	listener, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("Error al abrir el puerto: %v", err)
	}

	// Crear servidor gRPC
	grpcServer := grpc.NewServer()

	// Registrar el servicio
	pb.RegisterServiciosCancionesServer(grpcServer, &servidorCanciones{})

	// Cargar canciones en el vector
	servicios.CargarCanciones(vectorCanciones)

	// Iniciar el servidor
	log.Println("Servidor gRPC escuchando en puerto 50053...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
