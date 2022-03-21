package main

import (
	"context"
	"github.com/HekapOo-hub/broker_client/internal/proto/orderpb"
	"github.com/HekapOo-hub/broker_client/internal/proto/positionpb"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"time"
)

const (
	PositionPort = ":50005"
	OrderPort    = ":50004"
)

func main() {
	ctx := context.Background()
	positionConn, err := grpc.Dial(PositionPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Warnf("dial to position server: %v", err)
		return
	}
	positionClient := positionpb.NewPositionServiceClient(positionConn)

	orderConn, err := grpc.Dial(OrderPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Warnf("dial to order server: %v", err)
		return
	}

	orderClient := orderpb.NewOrderServiceClient(orderConn)

	order1 := orderpb.Order{Symbol: "silver", AccountID: "1234", Quantity: 0.1, Side: "BUY", GuaranteedStopLoss: false}

	order2 := orderpb.Order{Symbol: "gold", AccountID: "1234", Quantity: 0.1, Side: "SELL", GuaranteedStopLoss: false, Leverage: true}

	_, err = orderClient.Create(ctx, &order1)
	if err != nil {
		log.Warnf("%v", err)
		return
	}
	_, err = orderClient.Create(ctx, &order2)
	if err != nil {
		log.Warnf("%v", err)
		return
	}

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
