package rest

import (
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/dex/modules/bankx/internal/types"
	dex "github.com/coinexchain/dex/types"
)

func TestCmd(t *testing.T) {
	dex.InitSdkConfig()
	sendReq := sendReq{
		Amount:     dex.NewCetCoins(100000000),
		UnlockTime: 0,
	}
	addr, _ := sdk.AccAddressFromBech32("coinex1px8alypku5j84qlwzdpynhn4nyrkagaytu5u4a")
	req := &http.Request{Method: "POST", URL: nil}
	req = mux.SetURLVars(req, map[string]string{"address": "coinex1px8alypku5j84qlwzdpynhn4nyrkagaytu5u4a"})
	msg, _ := sendReq.GetMsg(req, addr)
	assert.Equal(t, types.MsgSend{
		FromAddress: addr,
		ToAddress:   addr,
		Amount:      dex.NewCetCoins(100000000),
		UnlockTime:  0,
	}, msg)

	memoReq := memoReq{
		Required: true,
	}
	msg, _ = memoReq.GetMsg(req, addr)
	assert.Equal(t, types.MsgSetMemoRequired{
		Address:  addr,
		Required: true,
	}, msg)
}
