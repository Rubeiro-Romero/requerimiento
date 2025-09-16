package utilidades

import (
	"fmt"
	"io"
	"log"
	"time"

	pb "ServidorDeStreaming.local/grpc-ServidorDeStreaming/serviciosStreaming"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

func DecodificarReproducir(reader io.Reader, canalSincronizacion chan struct{}) {
	streamer, format, err := mp3.Decode(io.NopCloser(reader))
	if err != nil {
		log.Fatalf("error decodificando MP3: %v", err)
	}
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/2))

	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		close(canalSincronizacion)
	})))
}
func RecibirCancion(
	stream pb.AudioService_EnviarCancionMedianteStreamClient,
	writer *io.PipeWriter,
	canalSincronizacion chan struct{},
) {
	noFragmento := 0
	for {
		fragmento, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("Canción recibida completa.")
			writer.Close()
			break
		}
		if err != nil {
			log.Fatalf("Error recibiendo chunk: %m", err)
		}

		noFragmento++
		fmt.Printf("\n Fragmento #%d recibido (%d bytes) reproduciendo ...",
			noFragmento, len(fragmento.Data))

		if _, err := writer.Write(fragmento.Data); err != nil {
			log.Printf("Error escribiendo en pipe: %m", err)
			break
		}
	}

	// Esperar hasta que termine la reproducción
	<-canalSincronizacion
	fmt.Println("✅ Reproducción finalizada.")
}
