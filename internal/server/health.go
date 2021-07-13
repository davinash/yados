package server

import (
	"context"
	"time"

	pb "github.com/davinash/yados/internal/proto/gen"
)

// UpdateHealthStatus updates the locally the health status of other peers
func (server *YadosServer) UpdateHealthStatus(ctx context.Context, request *pb.HealthStatusRequest) (*pb.HealthStatusReply, error) {
	server.logger.Debug("Received UpdateHealthStatus")

	server.mutex.Lock()
	defer server.mutex.Unlock()

	remotePeer := server.peers[request.Member.Name]
	remotePeer.ReceivedHCLastAt = time.Now().Unix()
	remotePeer.Status = pb.Member_Healthy
	return &pb.HealthStatusReply{}, nil
}

//StartHealthCheck starts the health check go routine
func (server *YadosServer) StartHealthCheck() error {
	server.wg.Add(1)
	go func() {
		defer server.wg.Done()
		ticker := time.NewTicker(time.Duration(server.hcTriggerDuration) * time.Second)
		server.performCheck()
		for {
			select {
			case <-ticker.C:
				server.SendHealthToPeers()
				server.performCheck()
			case <-server.healthCheckChan:
				ticker.Stop()
				server.logger.Info("Received close signal, shutting down")
				return
			}
		}
	}()
	server.logger.Info("Health Check is started ...")
	return nil
}

// StopHealthCheck Stops the health check go-routine
func (server *YadosServer) StopHealthCheck() {
	close(server.healthCheckChan)
}

//SendHealthToPeers Send the current status of this server to all the peers
func (server *YadosServer) SendHealthToPeers() {
	server.mutex.Lock()
	defer server.mutex.Unlock()

	for _, peer := range server.peers {
		conn, remotePeer, err := GetPeerConn(peer.Address, peer.Port)
		if err != nil {
			return
		}
		server.logger.Debugf("Sending Health Status to %s:%d", peer.Address, peer.Port)
		_, err = remotePeer.UpdateHealthStatus(context.Background(), &pb.HealthStatusRequest{Member: server.self})
		if err != nil {
			server.logger.Warnf("Failed to SendHealthStatus to %s:%d, error = %v", peer.Address, peer.Port, err)
		}

		err = conn.Close()
		if err != nil {
			server.logger.Warnf("Failed to close connection for %s:%d, error = %v", peer.Address, peer.Port, err)
		}
	}
}

func (server *YadosServer) performCheck() {
	server.logger.Debug("Performing health check")
	currentTime := time.Now().Unix()
	server.mutex.Lock()
	for _, peer := range server.peers {
		server.logger.Infof("%s:[%s:%d] => %v", peer.Name, peer.Address, peer.Port, peer.ReceivedHCLastAt)
		diff := currentTime - currentTime
		if diff > server.hcMaxWaitTime {
			server.logger.Infof("%s:[%s:%d] => Did not respond", peer.Name, peer.Address, peer.Port)
			peer.Status = pb.Member_UnHealthy
		}
	}
	server.mutex.Unlock()
}
