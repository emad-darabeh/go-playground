package main

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	pb "grpc-over-unix-sockets/proto"
)

const (
	minBackoff = 1 * time.Second
	maxBackoff = 30 * time.Second
)

// nodeServer handles inbound Send calls from the peer.
type nodeServer struct {
	pb.UnimplementedLinkServiceServer
	mySocket string
}

func (s *nodeServer) Send(ctx context.Context, msg *pb.Message) (*pb.Ack, error) {
	slog.Info("[PEER->ME]",
		"node", s.mySocket,
		"id", msg.Id,
		"command", msg.Command,
		"payload_len", len(msg.Payload),
	)
	return &pb.Ack{Status: "ok"}, nil
}

func newPeerConn(peerSocket string) (*grpc.ClientConn, error) {
	return grpc.NewClient(
		"unix://"+peerSocket,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// Force unix network — gRPC would otherwise attempt DNS on the path.
		grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
			var d net.Dialer
			return d.DialContext(ctx, "unix", addr)
		}),
	)
}

func backoff(attempt int) time.Duration {
	d := time.Duration(math.Pow(2, float64(attempt))) * minBackoff
	if d > maxBackoff {
		return maxBackoff
	}
	return d
}

// sendLoop periodically sends a message to the peer, retrying on transient errors.
func sendLoop(ctx context.Context, conn *grpc.ClientConn, mySocket string) {
	client := pb.NewLinkServiceClient(conn)
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	seq := 0

	for {
		select {
		case <-ctx.Done():
			return
		case t := <-ticker.C:
			seq++
			msg := &pb.Message{
				Id:      fmt.Sprintf("%s-%d", mySocket, seq),
				Command: "ping",
				Payload: []byte(fmt.Sprintf(`{"ts":%d,"seq":%d}`, t.UnixMilli(), seq)),
			}

			var attempt int
			for {
				callCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
				ack, err := client.Send(callCtx, msg)
				cancel()

				if err == nil {
					slog.Info("[ME->PEER]",
						"node", mySocket,
						"id", msg.Id,
						"command", msg.Command,
						"ack", ack.Status,
					)
					break
				}

				if ctx.Err() != nil {
					return
				}

				c := status.Code(err)
				if c != codes.Unavailable && c != codes.DeadlineExceeded && c != codes.Unknown {
					slog.Error("non-retryable send error", "node", mySocket, "err", err)
					break
				}

				wait := backoff(attempt)
				slog.Warn("peer unavailable, retrying", "node", mySocket, "err", err, "backoff_s", wait.Seconds())
				select {
				case <-time.After(wait):
					attempt++
				case <-ctx.Done():
					return
				}
			}
		}
	}
}

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	mySocket := os.Getenv("MY_SOCKET")
	peerSocket := os.Getenv("PEER_SOCKET")
	if mySocket == "" || peerSocket == "" {
		slog.Error("MY_SOCKET and PEER_SOCKET env vars must both be set")
		os.Exit(1)
	}

	if err := os.Remove(mySocket); err != nil && !os.IsNotExist(err) {
		slog.Error("failed to remove stale socket", "path", mySocket, "err", err)
		os.Exit(1)
	}

	lis, err := net.Listen("unix", mySocket)
	if err != nil {
		slog.Error("failed to listen on socket", "path", mySocket, "err", err)
		os.Exit(1)
	}
	if err := os.Chmod(mySocket, 0666); err != nil {
		slog.Error("failed to chmod socket", "path", mySocket, "err", err)
		os.Exit(1)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterLinkServiceServer(grpcServer, &nodeServer{mySocket: mySocket})

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		slog.Info("listening for peer connections", "socket", mySocket)
		if err := grpcServer.Serve(lis); err != nil {
			slog.Error("grpc server error", "err", err)
		}
	}()

	conn, err := newPeerConn(peerSocket)
	if err != nil {
		slog.Error("failed to create peer gRPC client", "peer", peerSocket, "err", err)
		os.Exit(1)
	}
	defer conn.Close()

	go sendLoop(ctx, conn, mySocket)

	<-ctx.Done()
	slog.Info("shutdown signal received", "node", mySocket)
	grpcServer.GracefulStop()
}
