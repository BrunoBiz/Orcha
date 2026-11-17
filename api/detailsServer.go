package api

import "context"

func (server *Server) postDetailsServer(ctx context.Context, input *ServerRequest) (*ServerResponse, error) {
	return server.runServerOperation(ctx, input, "details")
}
