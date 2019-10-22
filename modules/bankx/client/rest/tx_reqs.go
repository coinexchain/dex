package rest

import (
	"net/http"

	"github.com/gorilla/mux"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/coinexchain/dex/client/restutil"
	"github.com/coinexchain/dex/modules/bankx/internal/types"
)

type (
	sendReq struct {
		BaseReq    rest.BaseReq `json:"base_req"`
		Amount     sdk.Coins    `json:"amount"`
		UnlockTime int64        `json:"unlock_time"`
	}

	memoReq struct {
		BaseReq  rest.BaseReq `json:"base_req"`
		Required bool         `json:"memo_required"`
	}

	sendSupervisedReq struct {
		BaseReq    rest.BaseReq `json:"base_req"`
		Amount     sdk.Coin     `json:"amount"`
		UnlockTime int64        `json:"unlock_time"`
		Sender     string       `json:"sender,omitempty"`
		Supervisor string       `json:"supervisor,omitempty"`
		Reward     int64        `json:"reward,omitempty"`
		Operation  byte         `json:"operation"`
	}
)

func (req *sendReq) New() restutil.RestReq {
	return new(sendReq)
}
func (req *sendReq) GetBaseReq() *rest.BaseReq {
	return &req.BaseReq
}
func (req *sendReq) GetMsg(r *http.Request, sender sdk.AccAddress) (sdk.Msg, error) {
	toAddr := getAddr(r)
	return types.NewMsgSend(sender, toAddr, req.Amount, req.UnlockTime), nil
}

func (req *memoReq) New() restutil.RestReq {
	return new(memoReq)
}
func (req *memoReq) GetBaseReq() *rest.BaseReq {
	return &req.BaseReq
}
func (req *memoReq) GetMsg(r *http.Request, addr sdk.AccAddress) (sdk.Msg, error) {
	return types.NewMsgSetTransferMemoRequired(addr, req.Required), nil
}

func (req *sendSupervisedReq) New() restutil.RestReq {
	return new(sendSupervisedReq)
}
func (req *sendSupervisedReq) GetBaseReq() *rest.BaseReq {
	return &req.BaseReq
}
func (req *sendSupervisedReq) GetMsg(r *http.Request, addr sdk.AccAddress) (sdk.Msg, error) {
	toAddr := getAddr(r)

	var fromAddr sdk.AccAddress
	var supervisorAddr sdk.AccAddress
	var err error

	if req.Sender != "" {
		if fromAddr, err = sdk.AccAddressFromBech32(req.Sender); err != nil {
			return nil, err
		}
	}
	if req.Supervisor != "" {
		if supervisorAddr, err = sdk.AccAddressFromBech32(req.Supervisor); err != nil {
			return nil, err
		}
	}

	if req.Operation == types.Return || req.Operation == types.EarlierUnlockBySupervisor {
		supervisorAddr = addr
	} else {
		fromAddr = addr
	}
	return types.NewMsgSupervisedSend(fromAddr, supervisorAddr, toAddr, req.Amount, req.UnlockTime,
		req.Reward, req.Operation), nil
}

func getAddr(r *http.Request) sdk.AccAddress {
	vars := mux.Vars(r)
	addr, err := sdk.AccAddressFromBech32(vars["address"])
	if err != nil {
		return nil
	}
	return addr
}
