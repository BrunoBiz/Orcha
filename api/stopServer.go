package api

import "context"

func (server *Server) postStopServer(ctx context.Context, input *ServerRequest) (*ServerResponse, error) {
	return server.runServerOperation(ctx, input, "stop")
}
