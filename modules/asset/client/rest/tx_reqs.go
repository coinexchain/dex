package rest

import (
	"net/http"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/coinexchain/dex/client/restutil"
	"github.com/coinexchain/dex/modules/asset/internal/types"
)

type (
	// issueReq defines the properties of a issue token request's body
	issueReq struct {
		BaseReq          rest.BaseReq `json:"base_req" yaml:"base_req"`
		Name             string       `json:"name" yaml:"name"`
		Symbol           string       `json:"symbol" yaml:"symbol"`
		TotalSupply      string       `json:"total_supply" yaml:"total_supply"`
		Mintable         bool         `json:"mintable" yaml:"mintable"`
		Burnable         bool         `json:"burnable" yaml:"burnable"`
		AddrForbiddable  bool         `json:"addr_forbiddable" yaml:"addr_forbiddable"`
		TokenForbiddable bool         `json:"token_forbiddable" yaml:"token_forbiddable"`
		URL              string       `json:"url" yaml:"url"`
		Description      string       `json:"description" yaml:"description"`
		Identity         string       `json:"identity" yaml:"identity"`
	}

	// transferOwnerReq defines the properties of a transfer ownership request's body.
	transferOwnerReq struct {
		BaseReq  rest.BaseReq   `json:"base_req" yaml:"base_req"`
		NewOwner sdk.AccAddress `json:"new_owner" yaml:"new_owner"`
	}

	// mintTokenReq defines the properties of a mint token request's body.
	mintTokenReq struct {
		BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
		Amount  string       `json:"amount" yaml:"amount"`
	}

	// burnTokenReq defines the properties of a burn token request's body.
	burnTokenReq struct {
		BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
		Amount  string       `json:"amount" yaml:"amount"`
	}

	// forbidTokenReq defines the properties of a forbid token request's body.
	forbidTokenReq struct {
		BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
	}
	// unforbidTokenReq defines the properties of a unforbid token request's body.
	unForbidTokenReq struct {
		BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
	}
	// addrListReq defines the properties of a whitelist or forbidden addr request's body.
	addrListReq struct {
		BaseReq  rest.BaseReq     `json:"base_req" yaml:"base_req"`
		AddrList []sdk.AccAddress `json:"addr_list" yaml:"addr_list"`
	}
	// modifyTokenInfoReq defines the properties of a modify token info request's body.
	modifyTokenInfoReq struct {
		BaseReq     rest.BaseReq `json:"base_req" yaml:"base_req"`
		URL         *string      `json:"url,omitempty" yaml:"url,omitempty"`
		Description *string      `json:"description,omitempty" yaml:"description,omitempty"`
		Identity    *string      `json:"identity,omitempty" yaml:"identity,omitempty"`
	}
)

func (req *transferOwnerReq) New() restutil.RestReq {
	return new(transferOwnerReq)
}
func (req *transferOwnerReq) GetBaseReq() *rest.BaseReq {
	return &req.BaseReq
}
func (req *transferOwnerReq) GetMsg(r *http.Request, sender sdk.AccAddress) (sdk.Msg, error) {
	return types.NewMsgTransferOwnership(symbol, sender, req.NewOwner), nil
}

//func (req *mintTokenReq) New() restutil.RestReq {
//	return new(mintTokenReq)
//}
//func (req *mintTokenReq) GetBaseReq() *rest.BaseReq {
//	return &req.BaseReq
//}
//func (req *mintTokenReq) GetMsg(r *http.Request, sender sdk.AccAddress) (sdk.Msg, error) {
//	return types.NewMsgMintToken(symbol, amt, owner)
//}
