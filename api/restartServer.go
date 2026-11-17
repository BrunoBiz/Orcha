package api

import "context"

func (server *Server) postRestartServer(ctx context.Context, input *ServerRequest) (*ServerResponse, error) {
	return server.runServerOperation(ctx, input, "restart")
}
