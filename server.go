package main

import (
	"os"
	"strconv"

	"github.com/MRegterschot/GbxRemoteGo/events"
	"github.com/MRegterschot/GbxRemoteGo/gbxclient"
	"github.com/MRegterschot/GoController/config"
	"go.uber.org/zap"
)

type Server struct {
	Client *gbxclient.GbxClient
}

func NewServer() *Server {
	server := &Server{
		Client: gbxclient.NewGbxClient(gbxclient.Options{}),
	}

	// Register event handlers
	onConnectionChan := make(chan interface{})
	server.Client.Events.On("connect", onConnectionChan)
	go handleConnect(onConnectionChan)

	onDisconnectChan := make(chan interface{})
	server.Client.Events.On("disconnect", onDisconnectChan)
	go handleDisconnect(onDisconnectChan)

	server.Client.OnEcho = append(server.Client.OnEcho, server.onEcho)

	return server
}

func (s *Server) Connect() error {
	host, port := config.AppEnv.Host, config.AppEnv.Port
	zap.L().Info("Connecting to server", zap.String("host", host), zap.Int("port", port))
	if err := s.Client.Connect(host, port); err != nil {
		zap.L().Error("Failed to connect to server", zap.Error(err))
		return err
	}
	zap.L().Info("Connected to server")
	return nil
}

func (s *Server) Authenticate() error {
	user, pass := config.AppEnv.User, config.AppEnv.Pass
	zap.L().Info("Authenticating with server", zap.String("user", user))
	if err := s.Client.Authenticate(user, pass); err != nil {
		zap.L().Error("Failed to authenticate with server", zap.Error(err))
		return err
	}
	zap.L().Info("Authenticated with server")
	return nil
}

func (s *Server) Disconnect() error {
	zap.L().Info("Disconnecting from server")
	if err := s.Client.Disconnect(); err != nil {
		zap.L().Error("Failed to disconnect from server", zap.Error(err))
		return err
	}
	zap.L().Info("Disconnected from server")
	return nil
}

func handleConnect(eventChan chan interface{}) {
	for {
		select {
		case event := <-eventChan:
			if connected, ok := event.(bool); ok {
				if connected {
					zap.L().Info("Connected to server")
				} else {
					zap.L().Info("Disconnected from server")
				}
			} else {
				zap.L().Error("Invalid event type for connect.")
			}
		}
	}
}

func handleDisconnect(eventChan chan interface{}) {
	for {
		select {
		case event := <-eventChan:
			if msg, ok := event.(string); ok {
				zap.L().Info("Disconnected from server", zap.String("message", msg))
			} else {
				zap.L().Error("Invalid event type for disconnect.")
			}
		}
	}
}

func (s *Server) onEcho(client *gbxclient.GbxClient, echoEvent events.EchoEventArgs) {
	public, err := strconv.Atoi(echoEvent.Public)
	if err != nil {
		zap.L().Error("Failed to convert public to int", zap.Error(err))
		return
	}

	if echoEvent.Internal == "GoController" && public != GetController().StartTime {
		zap.L().Fatal("Another instance of GoController has started! Exiting...")
		os.Exit(1)
	}
}