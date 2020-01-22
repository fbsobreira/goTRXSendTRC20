package main

import (
	"context"
	"encoding/hex"
	"flag"
	"log"
	"math/big"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"

	proto "github.com/golang/protobuf/proto"
	"github.com/tronprotocol/grpc-gateway/api"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "path of server config file")
	flag.Parse()

	zapCfg := zap.NewDevelopmentConfig()
	zapCfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	zapCfg.Level.SetLevel(zap.InfoLevel)
	l, err := zapCfg.Build()
	if err != nil {
		log.Panic("Failed to init zap global logger, no zap log will be shown till zap is properly initialized: ", err)
	}
	zap.ReplaceGlobals(l)

	config := loadConfig(configPath)

	conn, err := grpc.Dial(config.API.FullNode, grpc.WithInsecure())
	if err != nil {
		zap.L().Fatal("Cannot connect", zap.Error(err))
		return
	}
	defer conn.Close()
	ctx := context.Background()
	cli := api.NewWalletClient(conn)

	acc := NewAccount(config.PrivateKey)
	zap.L().Info("Address", zap.String("Address", acc.Address))
	tx := acc.SendTRC20(ctx, getBytesFromAddress("TKTcfBEKpp5ZRPwmiZ8SfLx8W7CDZ7PHCY"), getBytesFromAddress("TKTcfBEKpp5ZRPwmiZ8SfLx8W7CDZ7PHCY"), *big.NewInt(10), cli)
	err = acc.sign(tx)
	if err != nil {
		zap.L().Error("Failed to sign", zap.Error(err))
	}
	rawDataBytes, _ := proto.Marshal(tx.GetRawData())
	zap.L().Info("RAW", zap.Any("TX", hex.EncodeToString(rawDataBytes)))

	// broadcast
	result, err := cli.BroadcastTransaction(ctx, tx)

	zap.L().Info("Result", zap.Bool("Success", result.Result), zap.String("Message", hex.EncodeToString(result.Message)))

}
