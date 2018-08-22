package rpc

import (
	"context"
	"fmt"

	"github.com/iost-official/Go-IOS-Protocol/account"
	"github.com/iost-official/Go-IOS-Protocol/core/event"
	"github.com/iost-official/Go-IOS-Protocol/core/new_tx"
	"github.com/iost-official/Go-IOS-Protocol/ilog"
	"time"
)

//go:generate mockgen -destination mock_rpc/mock_rpc.go -package rpc_mock github.com/iost-official/Go-IOS-Protocol/new_rpc ApisServer

// RPCServer is the class of RPC server
type RPCServer struct {
}

// newRPCServer
func newRPCServer() *RPCServer {
	s := &RPCServer{}
	return s
}

// GetHeight ...
func (s *RPCServer) GetHeight(ctx context.Context, void *VoidReq) (*HeightRes, error) {
	return &HeightRes{
		Height: bchain.Length(),
	}, nil
}

// GetTxByHash ...
func (s *RPCServer) GetTxByHash(ctx context.Context, hash *HashReq) (*tx.TxRaw, error) {
	if hash == nil {
		return nil, fmt.Errorf("argument cannot be nil pointer")
	}
	txHash := hash.Hash

	trx, err := txdb.Get(txHash)
	if err != nil {
		return nil, err
	}
	txRaw := trx.ToTxRaw()
	return txRaw, nil
}

// GetBlockByHash ...
func (s *RPCServer) GetBlockByHash(ctx context.Context, blkHashReq *BlockByHashReq) (*BlockInfo, error) {
	if blkHashReq == nil {
		return nil, fmt.Errorf("argument cannot be nil pointer")
	}

	hash := blkHashReq.Hash
	complete := blkHashReq.Complete

	blk, _ := bchain.GetBlockByHash(hash)
	if blk == nil {
		blk, _ = bc.GetBlockByHash(hash)
	}
	if blk == nil {
		return nil, fmt.Errorf("cant find the block")
	}
	blkInfo := &BlockInfo{
		Head:   &blk.Head,
		Txs:    make([]*tx.TxRaw, 0),
		Txhash: make([][]byte, 0),
	}
	for _, trx := range blk.Txs {
		if complete {
			blkInfo.Txs = append(blkInfo.Txs, trx.ToTxRaw())
		} else {
			blkInfo.Txhash = append(blkInfo.Txhash, trx.Hash())
		}
	}
	return blkInfo, nil
}

// GetBlockByNum ...
func (s *RPCServer) GetBlockByNum(ctx context.Context, blkNumReq *BlockByNumReq) (*BlockInfo, error) {
	if blkNumReq == nil {
		return nil, fmt.Errorf("argument cannot be nil pointer")
	}

	num := blkNumReq.Num
	complete := blkNumReq.Complete

	blk, _ := bchain.GetBlockByNumber(num)
	if blk == nil {
		blk, _ = bc.GetBlockByNumber(num)
	}
	if blk == nil {
		return nil, fmt.Errorf("cant find the block")
	}
	blkInfo := &BlockInfo{
		Head:   &blk.Head,
		Txs:    make([]*tx.TxRaw, 0),
		Txhash: make([][]byte, 0),
	}
	for _, trx := range blk.Txs {
		if complete {
			blkInfo.Txs = append(blkInfo.Txs, trx.ToTxRaw())
		} else {
			blkInfo.Txhash = append(blkInfo.Txhash, trx.Hash())
		}
	}
	return blkInfo, nil
}

// GetState ...
func (s *RPCServer) GetState(ctx context.Context, key *GetStateReq) (*GetStateRes, error) {
	if key == nil {
		return nil, fmt.Errorf("argument cannot be nil pointer")
	}
	return &GetStateRes{
		Value: visitor.BasicHandler.Get(key.Key),
	}, nil
}

// GetBalance ...
func (s *RPCServer) GetBalance(ctx context.Context, key *GetBalanceReq) (*GetBalanceRes, error) {
	if key == nil {
		return nil, fmt.Errorf("argument cannot be nil pointer")
	}
	return &GetBalanceRes{
		Balance: visitor.BalanceHandler.Balance(account.GetIDByPubkey(key.Pubkey)),
	}, nil
}

// SendRawTx ...
func (s *RPCServer) SendRawTx(ctx context.Context, rawTx *RawTxReq) (*SendRawTxRes, error) {
	if rawTx == nil {
		return nil, fmt.Errorf("argument cannot be nil pointer")
	}
	var trx tx.Tx
	err := trx.Decode(rawTx.Data)
	if err != nil {
		return nil, err
	}
	// add servi
	//tx.RecordTx(trx, tx.Data.Self())
	/*
		ret := txpool.TxPoolS.AddTx(trx)
		if ret != txpool.Success {
			return nil, fmt.Errorf("tx err:%v", ret)
		}
	*/
	res := SendRawTxRes{}
	res.Hash = trx.Hash()
	return &res, nil
}

// EstimateGas ...
func (s *RPCServer) EstimateGas(ctx context.Context, rawTx *RawTxReq) (*GasRes, error) {
	return nil, nil
}

func (s *RpcServer) Subscribe(req *SubscribeReq, res Apis_SubscribeServer) error {
	ec := event.GetEventCollectorInstance()
	sub := event.NewSubscription(100, req.Topics)
	ec.Subscribe(sub)
	defer ec.Unsubscribe(sub)

	timerChan := time.NewTicker(time.Minute).C
forloop:
	for {
		select {
		case <-timerChan:
			ilog.Debug("timeup in subscribe send")
			break forloop
		case ev := <-sub.ReadChan():
			err := res.Send(&SubscribeRes{Ev: ev})
			if err != nil {
				return err
			}
		default:
		}
	}
	return nil
}
