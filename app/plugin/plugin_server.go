// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package plugin

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/lasthyphen/dijetsnodego/app"

	pluginpb "github.com/lasthyphen/dijetsnodego/proto/pb/plugin"
)

// Server wraps a node so it can be served with the hashicorp plugin harness
type Server struct {
	pluginpb.UnsafeNodeServer
	app app.App
}

func NewServer(app app.App) *Server {
	return &Server{
		app: app,
	}
}

func (s *Server) Start(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, s.app.Start()
}

func (s *Server) Stop(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, s.app.Stop()
}

func (s *Server) ExitCode(context.Context, *emptypb.Empty) (*pluginpb.ExitCodeResponse, error) {
	exitCode, err := s.app.ExitCode()
	return &pluginpb.ExitCodeResponse{
		ExitCode: int32(exitCode),
	}, err
}
