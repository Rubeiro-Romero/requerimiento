package main

import (
	"context"
	"log"
	"time"

	"bufio"
	"fmt"
	"os"
	"strings"

	pbs "ServidorDeStreaming.local/grpc-ServidorDeStreaming/serviciosStreaming"
	menu "cliente.local/grpc-cliente/vistas"
	"google.golang.org/grpc"

	pb "ServidorDeCanciones.local/grpc-ServidorDeCanciones/serviciosCancion"
)

// ruta generada por protoc

func main() {
	// Conectar al servidor gRPC (localhost:50053)
	conn, err := grpc.Dial("localhost:50053", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("No se pudo conectar: %v", err)
	}
	defer conn.Close()

	// Crear cliente
	c := pb.NewServiciosCancionesClient(conn)

	// Contexto con timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	// Se captura el título de la canción a buscar
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("\nIngrese el título de la canción a buscar: ")
	tituloLeido, _ := reader.ReadString('\n')
	tituloLeido = strings.TrimSpace(tituloLeido)

	// Se crea un objeto de tipo DTO que contiene el título de la canción
	objPeticion := &pb.PeticionDTO{Titulo: tituloLeido}

	// Llamada al procedimiento remoto buscarCancion
	res, err := c.BuscarCancion(ctx, objPeticion)
	if err != nil {
		fmt.Printf("Error en la llamada gRPC: %v", err)
	}

	// Impresión de la respuesta
	fmt.Printf("Mensaje: %s\n", res.Mensaje)
	fmt.Printf("Código: %d\n", res.Codigo)
	if res.Codigo == 200 {
		fmt.Printf("Canción: %s, Tamaño: %d, URL: %s, Activada: %v\n",
			res.ObjCancion.Titulo,
			res.ObjCancion.Tamanio,
			res.ObjCancion.Url,
			res.ObjCancion.EsActivada)

		conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		client := pbs.NewAudioServiceClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
		defer cancel()
		menu.MostrarMenuPrincipal(client, ctx, tituloLeido)
	}

}
