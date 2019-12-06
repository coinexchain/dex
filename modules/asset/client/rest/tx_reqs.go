package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

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
	// the flowing 4 reqs defines the properties of a whitelist or forbidden addr request's body.
	addWhiteListReq struct {
		BaseReq   rest.BaseReq     `json:"base_req" yaml:"base_req"`
		Whitelist []sdk.AccAddress `json:"whitelist" yaml:"whitelist"`
	}
	removeWhiteListReq struct {
		BaseReq   rest.BaseReq     `json:"base_req" yaml:"base_req"`
		Whitelist []sdk.AccAddress `json:"whitelist" yaml:"whitelist"`
	}
	forbidAddrReq struct {
		BaseReq   rest.BaseReq     `json:"base_req" yaml:"base_req"`
		Addresses []sdk.AccAddress `json:"addresses" yaml:"addresses"`
	}
	unforbidAddrReq struct {
		BaseReq   rest.BaseReq     `json:"base_req" yaml:"base_req"`
		Addresses []sdk.AccAddress `json:"addresses" yaml:"addresses"`
	}
	// modifyTokenInfoReq defines the properties of a modify token info request's body.
	modifyTokenInfoReq struct {
		BaseReq          rest.BaseReq `json:"base_req" yaml:"base_req"`
		URL              *string      `json:"url,omitempty" yaml:"url,omitempty"`
		Description      *string      `json:"description,omitempty" yaml:"description,omitempty"`
		Identity         *string      `json:"identity,omitempty" yaml:"identity,omitempty"`
		Name             *string      `json:"name" yaml:"name"`
		TotalSupply      *string      `json:"total_supply" yaml:"total_supply"`
		Mintable         *bool        `json:"mintable" yaml:"mintable"`
		Burnable         *bool        `json:"burnable" yaml:"burnable"`
		AddrForbiddable  *bool        `json:"addr_forbiddable" yaml:"addr_forbiddable"`
		TokenForbiddable *bool        `json:"token_forbiddable" yaml:"token_forbiddable"`
	}
)

func (req *issueReq) New() restutil.RestReq {
	return new(issueReq)
}
func (req *issueReq) GetBaseReq() *rest.BaseReq {
	return &req.BaseReq
}
func (req *issueReq) GetMsg(r *http.Request, owner sdk.AccAddress) (sdk.Msg, error) {
	amt, ok := sdk.NewIntFromString(req.TotalSupply)
	if !ok {
		return nil, types.ErrInvalidTokenSupply(req.TotalSupply)
	}
	return types.NewMsgIssueToken(req.Name, req.Symbol, amt, owner,
		req.Mintable, req.Burnable, req.AddrForbiddable, req.TokenForbiddable,
		req.URL, req.Description, req.Identity), nil
}

func (req *transferOwnerReq) New() restutil.RestReq {
	return new(transferOwnerReq)
}
func (req *transferOwnerReq) GetBaseReq() *rest.BaseReq {
	return &req.BaseReq
}
func (req *transferOwnerReq) GetMsg(r *http.Request, owner sdk.AccAddress) (sdk.Msg, error) {
	symbol := getSymbol(r)
	return types.NewMsgTransferOwnership(symbol, owner, req.NewOwner), nil
}

func (req *mintTokenReq) New() restutil.RestReq {
	return new(mintTokenReq)
}
func (req *mintTokenReq) GetBaseReq() *rest.BaseReq {
	return &req.BaseReq
}
func (req *mintTokenReq) GetMsg(r *http.Request, owner sdk.AccAddress) (sdk.Msg, error) {
	symbol := getSymbol(r)
	amt, ok := sdk.NewIntFromString(req.Amount)
	if !ok {
		return nil, types.ErrInvalidTokenMintAmt(req.Amount)
	}
	msg := types.NewMsgMintToken(symbol, amt, owner)
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	return types.NewMsgMintToken(symbol, amt, owner), nil
}

func (req *burnTokenReq) New() restutil.RestReq {
	return new(burnTokenReq)
}
func (req *burnTokenReq) GetBaseReq() *rest.BaseReq {
	return &req.BaseReq
}
func (req *burnTokenReq) GetMsg(r *http.Request, owner sdk.AccAddress) (sdk.Msg, error) {
	symbol := getSymbol(r)
	amt, ok := sdk.NewIntFromString(req.Amount)
	if !ok {
		return nil, types.ErrInvalidTokenBurnAmt(req.Amount)
	}
	return types.NewMsgBurnToken(symbol, amt, owner), nil
}

func (req *forbidTokenReq) New() restutil.RestReq {
	return new(forbidTokenReq)
}
func (req *forbidTokenReq) GetBaseReq() *rest.BaseReq {
	return &req.BaseReq
}
func (req *forbidTokenReq) GetMsg(r *http.Request, owner sdk.AccAddress) (sdk.Msg, error) {
	symbol := getSymbol(r)
	return types.NewMsgForbidToken(symbol, owner), nil
}

func (req *unForbidTokenReq) New() restutil.RestReq {
	return new(unForbidTokenReq)
}
func (req *unForbidTokenReq) GetBaseReq() *rest.BaseReq {
	return &req.BaseReq
}
func (req *unForbidTokenReq) GetMsg(r *http.Request, owner sdk.AccAddress) (sdk.Msg, error) {
	symbol := getSymbol(r)
	return types.NewMsgUnForbidToken(symbol, owner), nil
}

func (req *addWhiteListReq) New() restutil.RestReq {
	return new(addWhiteListReq)
}
func (req *addWhiteListReq) GetBaseReq() *rest.BaseReq {
	return &req.BaseReq
}
func (req *addWhiteListReq) GetMsg(r *http.Request, owner sdk.AccAddress) (sdk.Msg, error) {
	symbol := getSymbol(r)
	return types.NewMsgAddTokenWhitelist(symbol, owner, req.Whitelist), nil
}

func (req *removeWhiteListReq) New() restutil.RestReq {
	return new(removeWhiteListReq)
}
func (req *removeWhiteListReq) GetBaseReq() *rest.BaseReq {
	return &req.BaseReq
}
func (req *removeWhiteListReq) GetMsg(r *http.Request, owner sdk.AccAddress) (sdk.Msg, error) {
	symbol := getSymbol(r)
	return types.NewMsgRemoveTokenWhitelist(symbol, owner, req.Whitelist), nil
}

func (req *forbidAddrReq) New() restutil.RestReq {
	return new(forbidAddrReq)
}
func (req *forbidAddrReq) GetBaseReq() *rest.BaseReq {
	return &req.BaseReq
}
func (req *forbidAddrReq) GetMsg(r *http.Request, owner sdk.AccAddress) (sdk.Msg, error) {
	symbol := getSymbol(r)
	return types.NewMsgForbidAddr(symbol, owner, req.Addresses), nil
}

func (req *unforbidAddrReq) New() restutil.RestReq {
	return new(unforbidAddrReq)
}
func (req *unforbidAddrReq) GetBaseReq() *rest.BaseReq {
	return &req.BaseReq
}
func (req *unforbidAddrReq) GetMsg(r *http.Request, owner sdk.AccAddress) (sdk.Msg, error) {
	symbol := getSymbol(r)
	return types.NewMsgUnForbidAddr(symbol, owner, req.Addresses), nil
}

func (req *modifyTokenInfoReq) New() restutil.RestReq {
	return new(modifyTokenInfoReq)
}
func (req *modifyTokenInfoReq) GetBaseReq() *rest.BaseReq {
	return &req.BaseReq
}
func (req *modifyTokenInfoReq) GetMsg(r *http.Request, owner sdk.AccAddress) (sdk.Msg, error) {
	symbol := getSymbol(r)
	url := getNewTokenStrInfo(req.URL)
	description := getNewTokenStrInfo(req.Description)
	identity := getNewTokenStrInfo(req.Identity)
	name := getNewTokenStrInfo(req.Name)
	supply := getNewTokenStrInfo(req.TotalSupply)
	mintable := getNewTokenBoolInfo(req.Mintable)
	burnable := getNewTokenBoolInfo(req.Burnable)
	addrForbiddable := getNewTokenBoolInfo(req.AddrForbiddable)
	tokenForbiddable := getNewTokenBoolInfo(req.TokenForbiddable)

	return types.NewMsgModifyTokenInfo(symbol, url, description, identity, owner,
		name, supply, mintable, burnable, addrForbiddable, tokenForbiddable), nil
}

func getNewTokenStrInfo(ptr *string) string {
	if ptr != nil {
		return *ptr
	}
	return types.DoNotModifyTokenInfo
}
func getNewTokenBoolInfo(ptr *bool) string {
	if ptr != nil {
		return fmt.Sprintf("%v", *ptr)
	}
	return types.DoNotModifyTokenInfo
}

func getSymbol(r *http.Request) string {
	vars := mux.Vars(r)
	return vars[symbol]
}
