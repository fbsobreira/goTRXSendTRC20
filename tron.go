package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	proto "github.com/golang/protobuf/proto"
	"go.uber.org/zap"
	"golang.org/x/crypto/sha3"

	"crypto/ecdsa"

	"github.com/tronprotocol/grpc-gateway/api"
	"github.com/tronprotocol/grpc-gateway/core"

	"github.com/sasaxie/go-client-api/common/base58"
	"github.com/sasaxie/go-client-api/common/hexutil"
)

// TRXAccount Status
type TRXAccount struct {
	Address string
	Key     *ecdsa.PrivateKey
}

func NewAccount(pk string) *TRXAccount {
	bytesPk, _ := hex.DecodeString(pk)
	pkECDSA, _ := crypto.ToECDSA(bytesPk)
	acc := TRXAccount{Key: pkECDSA}

	address := crypto.PubkeyToAddress(pkECDSA.PublicKey).Hex()
	address = "41" + address[2:]
	addressBytes, _ := hex.DecodeString(address)
	acc.Address = getAddressFromBytes(addressBytes)
	return &acc
}

func (s *TRXAccount) SendTRC20(ctx context.Context, trcAddress []byte, toAddress []byte, value big.Int, cli api.WalletClient) *core.Transaction {
	transferFnSignature := []byte("transfer(address,uint256)")
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]
	fmt.Println(hexutil.Encode(methodID))                   // 0xa9059cbb
	paddedAddress := common.LeftPadBytes(toAddress[1:], 32) // remove first byte
	fmt.Println(hexutil.Encode(paddedAddress))

	paddedAmount := common.LeftPadBytes(value.Bytes(), 32)
	fmt.Println(hexutil.Encode(paddedAmount))
	var dataPack []byte
	dataPack = append(dataPack, methodID...)
	dataPack = append(dataPack, paddedAddress...)
	dataPack = append(dataPack, paddedAmount...)
	fmt.Println(hexutil.Encode(dataPack))

	contract := &core.TriggerSmartContract{
		OwnerAddress:    getBytesFromAddress(s.Address),
		ContractAddress: trcAddress,
		Data:            dataPack,
	}

	tx, err := cli.TriggerContract(ctx, contract)
	if err != nil {
		zap.L().Fatal("Error trigger Contract", zap.Error(err))
	}
	return tx.GetTransaction()
}

func getBytesFromAddress(address string) []byte {
	return base58.DecodeCheck(address)
}

func getAddressFromBytes(address []byte) string {
	return base58.EncodeCheck(address)
}

func (s *TRXAccount) sign(tx *core.Transaction) error {
	// hash
	rawDataBytes, err := proto.Marshal(tx.GetRawData())
	if err != nil {
		return err
	}
	// use sha3 algorithm
	h := sha256.New()
	h.Write(rawDataBytes)
	hash := h.Sum(nil)

	// sign
	signature, err := crypto.Sign(hash, s.Key)
	if err != nil {
		return err
	}
	// add signature
	tx.Signature = append(tx.Signature, signature)
	//fmt.Println(hex.EncodeToString(account.Key.D.Bytes()))
	return nil
}
