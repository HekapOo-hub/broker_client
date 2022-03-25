package main

import (
	"context"
	"github.com/HekapOo-hub/broker_client/internal/config"
	"github.com/HekapOo-hub/broker_client/internal/proto/positionpb"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"time"
)

func main() {
	ctx := context.Background()
	ports, err := config.GetPorts()
	if err != nil {
		log.Warnf("%v", err)
		return
	}
	positionConn, err := grpc.Dial(ports.Position, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Warnf("dial to position server: %v", err)
		return
	}
	positionClient := positionpb.NewPositionServiceClient(positionConn)

	/*orderConn, err := grpc.Dial(ports.Order, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Warnf("dial to order server: %v", err)
		return
	}

	orderClient := orderpb.NewOrderServiceClient(orderConn)
	*/
	positionStream, err := positionClient.GetProfitLoss(ctx, &positionpb.AccountID{Value: "1234"})
	if err != nil {
		log.Warnf("%v", err)
	}
	go func() {
		for {
			resp, err := positionStream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				log.Fatalf("can not receive %v", err)
			}

			if resp == nil {
				log.Warnf("error with stream upload func in recieving!!!")
			}

			log.Infof("positionID: %v        Profit/Loss: %v", resp.PositionID, resp.Value)

			log.Info()

		}
	}()
	time.Sleep(time.Minute)
}
