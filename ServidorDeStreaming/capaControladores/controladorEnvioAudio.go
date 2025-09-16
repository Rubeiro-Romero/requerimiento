package capacontroladores

import (
	capafachadaservices "ServidorDeStreaming.local/grpc-ServidorDeStreaming/capaFachadaServices"
	pb "ServidorDeStreaming.local/grpc-ServidorDeStreaming/serviciosStreaming"
)

type ControladorServidor struct {
	pb.UnimplementedAudioServiceServer
}

// Implementación del procedimiento remoto
func (s *ControladorServidor) EnviarCancionMedianteStream(req *pb.PeticionDTO, stream pb.AudioService_EnviarCancionMedianteStreamServer) error {

	// Usamos la fachada directamente
	return capafachadaservices.StreamAudioFile(
		req.Titulo,
		func(data []byte) error {
			return stream.Send(&pb.FragmentoCancion{Data: data})
		})
}
