package asset

import (
	"bytes"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/tendermint/tendermint/libs/bech32"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	SeparateKeyPrefix   = []byte{0x3A}
	TokenStoreKeyPrefix = []byte{0x01}
	WhitelistKeyPrefix  = []byte{0x02}
	ForbidAddrKeyPrefix = []byte{0x03}
)

// Keeper encodes/decodes tokens using the go-amino (binary)
// encoding/decoding library.
type TokenKeeper struct {
	// The (unexposed) key used to access the store from the Context.
	key sdk.StoreKey

	ak  auth.AccountKeeper
	fck auth.FeeCollectionKeeper

	// The codec codec for binary encoding/decoding of token.
	cdc *codec.Codec

	paramSubspace params.Subspace
}

// NewKeeper returns a new Keeper that uses go-amino to
// (binary) encode and decode concrete Token.
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, paramstore params.Subspace,
	ak auth.AccountKeeper, fck auth.FeeCollectionKeeper) TokenKeeper {

	return TokenKeeper{
		key:           key,
		ak:            ak,
		fck:           fck,
		cdc:           cdc,
		paramSubspace: paramstore.WithKeyTable(ParamKeyTable()),
	}
}

func (tk TokenKeeper) checkPrecondition(ctx sdk.Context, msg sdk.Msg, symbol string, owner sdk.AccAddress) (Token, sdk.Error) {
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	token := tk.GetToken(ctx, symbol)
	if token == nil {
		return nil, ErrorTokenNotFound(fmt.Sprintf("token %s not found", symbol))
	}

	if !token.GetOwner().Equals(owner) {
		return nil, ErrorInvalidTokenOwner("Only token owner can do this action")
	}

	return token, nil
}

// setToken  implements token Keeper.
func (tk TokenKeeper) setToken(ctx sdk.Context, token Token) sdk.Error {
	symbol := token.GetSymbol()
	store := ctx.KVStore(tk.key)

	bz, err := tk.cdc.MarshalBinaryBare(token)
	if err != nil {
		return sdk.ErrInternal(err.Error())
	}
	store.Set(TokenStoreKey(symbol), bz)
	return nil
}

// removeToken remove an token for the asset mapper store.
func (tk TokenKeeper) removeToken(ctx sdk.Context, token Token) {
	symbol := token.GetSymbol()
	store := ctx.KVStore(tk.key)
	store.Delete(TokenStoreKey(symbol))
}

// GetToken implements token Keeper.
func (tk TokenKeeper) GetToken(ctx sdk.Context, symbol string) Token {
	store := ctx.KVStore(tk.key)
	bz := store.Get(TokenStoreKey(symbol))
	if bz == nil {
		return nil
	}
	return tk.decodeToken(bz)
}

// GetAllTokens returns all tokens in the token Keeper.
func (tk TokenKeeper) GetAllTokens(ctx sdk.Context) []Token {
	tokens := make([]Token, 0)

	tk.IterateTokenValue(ctx, func(token Token) (stop bool) {
		tokens = append(tokens, token)
		return false
	})

	return tokens
}

// IterateToken implements token Keeper
func (tk TokenKeeper) IterateTokenValue(ctx sdk.Context, process func(Token) (stop bool)) {
	store := ctx.KVStore(tk.key)
	iter := sdk.KVStorePrefixIterator(store, TokenStoreKeyPrefix)
	defer iter.Close()
	for {
		if !iter.Valid() {
			return
		}
		acc := tk.decodeToken(iter.Value())
		if process(acc) {
			return
		}
		iter.Next()
	}
}

//IssueToken - new token and store
func (tk TokenKeeper) IssueToken(ctx sdk.Context, msg MsgIssueToken) sdk.Error {
	token, err := NewToken(msg.Name, msg.Symbol, msg.TotalSupply, msg.Owner,
		msg.Mintable, msg.Burnable, msg.AddrForbiddable, msg.TokenForbiddable)

	if err != nil {
		return err
	}

	if tk.IsTokenExists(ctx, token.Symbol) {
		return ErrorDuplicateTokenSymbol(fmt.Sprintf("token symbol already exists in store"))
	}

	return tk.setToken(ctx, token)
}

//TransferOwnership - transfer token owner
func (tk TokenKeeper) TransferOwnership(ctx sdk.Context, msg MsgTransferOwnership) sdk.Error {
	token, err := tk.checkPrecondition(ctx, msg, msg.Symbol, msg.OriginalOwner)
	if err != nil {
		return err
	}

	if err := token.SetOwner(msg.NewOwner); err != nil {
		return ErrorInvalidTokenOwner("token new owner is invalid")
	}

	return tk.setToken(ctx, token)
}

//MintToken - mint token
func (tk TokenKeeper) MintToken(ctx sdk.Context, msg MsgMintToken) sdk.Error {
	token, err := tk.checkPrecondition(ctx, msg, msg.Symbol, msg.OwnerAddress)
	if err != nil {
		return err
	}

	if !token.GetMintable() {
		return ErrorInvalidTokenMint(fmt.Sprintf("token %s do not support mint", msg.Symbol))
	}

	if err := token.SetTotalMint(msg.Amount + token.GetTotalMint()); err != nil {
		return ErrorInvalidTokenMint(err.Error())
	}

	if err := token.SetTotalSupply(msg.Amount + token.GetTotalSupply()); err != nil {
		return ErrorInvalidTokenSupply(err.Error())
	}

	return tk.setToken(ctx, token)
}

//BurnToken - burn token
func (tk TokenKeeper) BurnToken(ctx sdk.Context, msg MsgBurnToken) sdk.Error {
	token, err := tk.checkPrecondition(ctx, msg, msg.Symbol, msg.OwnerAddress)
	if err != nil {
		return err
	}

	if !token.GetBurnable() {
		return ErrorInvalidTokenBurn(fmt.Sprintf("token %s do not support burn", msg.Symbol))
	}

	if err := token.SetTotalBurn(msg.Amount + token.GetTotalBurn()); err != nil {
		return ErrorInvalidTokenBurn(err.Error())
	}

	if err := token.SetTotalSupply(token.GetTotalSupply() - msg.Amount); err != nil {
		return ErrorInvalidTokenSupply(err.Error())
	}

	return tk.setToken(ctx, token)
}

// -----------------------------------------------------------------------------
// Forbid token with whitelist

func (tk TokenKeeper) addWhitelist(ctx sdk.Context, symbol string, whitelist []sdk.AccAddress) sdk.Error {
	store := ctx.KVStore(tk.key)
	for _, addr := range whitelist {
		store.Set(PrefixAddrStoreKey(WhitelistKeyPrefix, symbol, addr), []byte{})
	}

	return nil
}

func (tk TokenKeeper) removeWhitelist(ctx sdk.Context, symbol string, whitelist []sdk.AccAddress) sdk.Error {
	store := ctx.KVStore(tk.key)
	for _, addr := range whitelist {
		store.Delete(PrefixAddrStoreKey(WhitelistKeyPrefix, symbol, addr))
	}

	return nil
}

//ForbidToken - forbid token
func (tk TokenKeeper) ForbidToken(ctx sdk.Context, msg MsgForbidToken) sdk.Error {
	token, err := tk.checkPrecondition(ctx, msg, msg.Symbol, msg.OwnerAddress)
	if err != nil {
		return err
	}

	if !token.GetTokenForbiddable() {
		return ErrorInvalidTokenForbidden(fmt.Sprintf("token %s do not support forbid token", msg.Symbol))
	}
	if token.GetIsForbidden() {
		return ErrorInvalidTokenForbidden(fmt.Sprintf("token %s has been forbidden", msg.Symbol))
	}
	token.SetIsForbidden(true)

	return tk.setToken(ctx, token)
}

//UnForbidToken - unforbid token
func (tk TokenKeeper) UnForbidToken(ctx sdk.Context, msg MsgUnForbidToken) sdk.Error {
	token, err := tk.checkPrecondition(ctx, msg, msg.Symbol, msg.OwnerAddress)
	if err != nil {
		return err
	}

	if !token.GetTokenForbiddable() {
		return ErrorInvalidTokenForbidden(fmt.Sprintf("token %s do not support unforbid token", msg.Symbol))
	}
	if !token.GetIsForbidden() {
		return ErrorInvalidTokenForbidden(fmt.Sprintf("token %s has not been forbidden", msg.Symbol))
	}
	token.SetIsForbidden(false)

	return tk.setToken(ctx, token)
}

//AddTokenWhitelist - add token forbid whitelist
func (tk TokenKeeper) AddTokenWhitelist(ctx sdk.Context, msg MsgAddTokenWhitelist) sdk.Error {
	token, err := tk.checkPrecondition(ctx, msg, msg.Symbol, msg.OwnerAddress)
	if err != nil {
		return err
	}

	if !token.GetTokenForbiddable() {
		return ErrorInvalidTokenForbidden(fmt.Sprintf("token %s do not support forbid token and add whitelist", msg.Symbol))
	}
	if err = tk.addWhitelist(ctx, msg.Symbol, msg.Whitelist); err != nil {
		return ErrorInvalidTokenWhitelist(fmt.Sprintf("token whitelist is invalid"))
	}
	return nil
}

//RemoveTokenWhitelist - remove token forbid whitelist
func (tk TokenKeeper) RemoveTokenWhitelist(ctx sdk.Context, msg MsgRemoveTokenWhitelist) sdk.Error {
	token, err := tk.checkPrecondition(ctx, msg, msg.Symbol, msg.OwnerAddress)
	if err != nil {
		return err
	}

	if !token.GetTokenForbiddable() {
		return ErrorInvalidTokenForbidden(fmt.Sprintf("token %s do not support forbid token and remove whitelist", msg.Symbol))
	}
	if err = tk.removeWhitelist(ctx, msg.Symbol, msg.Whitelist); err != nil {
		return ErrorInvalidTokenWhitelist(fmt.Sprintf("token whitelist is invalid"))
	}
	return nil
}

// GetWhitelist returns whitelist in the token Keeper.
func (tk TokenKeeper) GetWhitelist(ctx sdk.Context, symbol string) []sdk.AccAddress {
	whitelist := make([]sdk.AccAddress, 0)
	keyPrefix := append(append(WhitelistKeyPrefix, symbol...), SeparateKeyPrefix...)

	tk.IterateAddrKeys(ctx, keyPrefix, func(key []byte) (stop bool) {
		addr := key[len(WhitelistKeyPrefix)+len(symbol)+len(SeparateKeyPrefix):]
		whitelist = append(whitelist, addr)
		return false
	})

	return whitelist
}

// -----------------------------------------------------------------------------
// Forbid address

func (tk TokenKeeper) addForbidAddress(ctx sdk.Context, symbol string, addresses []sdk.AccAddress) sdk.Error {
	store := ctx.KVStore(tk.key)
	for _, addr := range addresses {
		store.Set(PrefixAddrStoreKey(ForbidAddrKeyPrefix, symbol, addr), []byte{})
	}

	return nil
}

func (tk TokenKeeper) removeForbidAddress(ctx sdk.Context, symbol string, addresses []sdk.AccAddress) sdk.Error {
	store := ctx.KVStore(tk.key)
	for _, addr := range addresses {
		store.Delete(PrefixAddrStoreKey(ForbidAddrKeyPrefix, symbol, addr))
	}

	return nil
}

//ForbidAddress - add forbid addresses
func (tk TokenKeeper) ForbidAddress(ctx sdk.Context, msg MsgForbidAddr) sdk.Error {
	token, err := tk.checkPrecondition(ctx, msg, msg.Symbol, msg.OwnerAddr)
	if err != nil {
		return err
	}

	if !token.GetAddrForbiddable() {
		return ErrorInvalidTokenForbidden(fmt.Sprintf("token %s do not support forbid address", msg.Symbol))
	}
	if err = tk.addForbidAddress(ctx, msg.Symbol, msg.ForbidAddr); err != nil {
		return ErrorInvalidAddress(fmt.Sprintf("forbid addr is invalid"))
	}
	return nil
}

//UnForbidAddress - remove forbid addresses
func (tk TokenKeeper) UnForbidAddress(ctx sdk.Context, msg MsgUnForbidAddr) sdk.Error {
	token, err := tk.checkPrecondition(ctx, msg, msg.Symbol, msg.OwnerAddr)
	if err != nil {
		return err
	}

	if !token.GetAddrForbiddable() {
		return ErrorInvalidTokenForbidden(fmt.Sprintf("token %s do not support unforbid address", msg.Symbol))
	}
	if err = tk.removeForbidAddress(ctx, msg.Symbol, msg.UnForbidAddr); err != nil {
		return ErrorInvalidAddress(fmt.Sprintf("unforbid addr is invalid"))
	}
	return nil
}

// GetForbidAddr returns all forbidden addr in the token Keeper.
func (tk TokenKeeper) GetForbiddenAddr(ctx sdk.Context, symbol string) []sdk.AccAddress {
	addresses := make([]sdk.AccAddress, 0)
	keyPrefix := append(append(ForbidAddrKeyPrefix, symbol...), SeparateKeyPrefix...)

	tk.IterateAddrKeys(ctx, keyPrefix, func(key []byte) (stop bool) {
		addr := key[len(ForbidAddrKeyPrefix)+len(symbol)+len(SeparateKeyPrefix):]
		addresses = append(addresses, addr)
		return false
	})

	return addresses
}

// -----------------------------------------------------------------------------
// ExpectedAssertStatusKeeper

//IsTokenForbidden - check whether the coin's owner has forbidden "symbol", forbidding transmission and exchange.
func (tk TokenKeeper) IsTokenForbidden(ctx sdk.Context, symbol string) bool {
	token := tk.GetToken(ctx, symbol)
	if token != nil {
		return token.GetIsForbidden()
	}

	return true
}

// IsTokenExists - check whether there is a coin named "symbol"
func (tk TokenKeeper) IsTokenExists(ctx sdk.Context, symbol string) bool {
	return tk.GetToken(ctx, symbol) != nil
}

// IsTokenIssuer - check whether addr is a token issuer
func (tk TokenKeeper) IsTokenIssuer(ctx sdk.Context, symbol string, addr sdk.AccAddress) bool {
	if addr.Empty() {
		return false
	}

	token := tk.GetToken(ctx, symbol)
	return token != nil && token.GetOwner().Equals(addr)
}

func (tk TokenKeeper) IsForbiddenByTokenIssuer(ctx sdk.Context, symbol string, addr sdk.AccAddress) bool {
	token := tk.GetToken(ctx, symbol)
	if token == nil {
		return true
	}

	for _, forbidden := range tk.GetForbiddenAddr(ctx, symbol) {
		if addr.Equals(forbidden) {
			return true
		}
	}

	if !token.GetIsForbidden() {
		return false
	}

	for _, forbidden := range tk.GetWhitelist(ctx, symbol) {
		if addr.Equals(forbidden) {
			return false
		}
	}

	return true
}

// -----------------------------------------------------------------------------
// Params

// SetParams sets the asset module's parameters.
func (tk TokenKeeper) SetParams(ctx sdk.Context, params Params) {
	tk.paramSubspace.SetParamSet(ctx, &params)
}

// GetParams gets the asset module's parameters.
func (tk TokenKeeper) GetParams(ctx sdk.Context) (params Params) {
	tk.paramSubspace.GetParamSet(ctx, &params)
	return
}

// -----------------------------------------------------------------------------
// Misc.

// TokenStoreKey turn an token symbol to key used to get it from the asset store
func TokenStoreKey(symbol string) []byte {
	return append(TokenStoreKeyPrefix, []byte(symbol)...)
}

func (tk TokenKeeper) decodeToken(bz []byte) (token Token) {
	if err := tk.cdc.UnmarshalBinaryBare(bz, &token); err != nil {
		panic(err)
	}
	return
}

// PrefixAddrStoreKey - return paramsKeyPrefix-Symbol-:-AccAddress KEY
func PrefixAddrStoreKey(prefix []byte, symbol string, addr sdk.AccAddress) []byte {
	return append(append(append(prefix, symbol...), SeparateKeyPrefix...), addr...)
}

// GetAllAddrKeys return [] symbol:addr key. for get all whitelists or forbidden addresses string
func (tk TokenKeeper) GetAllAddrKeys(ctx sdk.Context, prefix []byte) []string {
	res := make([]string, 0)
	bech32PrefixAccAddr := sdk.GetConfig().GetBech32AccountAddrPrefix()

	tk.IterateAddrKeys(ctx, prefix, func(key []byte) (stop bool) {
		i := bytes.Index(key, SeparateKeyPrefix) + len(SeparateKeyPrefix)
		bech32Addr, err := bech32.ConvertAndEncode(bech32PrefixAccAddr, key[i:])
		if err != nil {
			panic(err)
		}
		s := string(key[len(prefix):i]) + bech32Addr
		res = append(res, s)
		return false
	})

	return res
}

// IterateAddrKeys implements token Keeper
func (tk TokenKeeper) IterateAddrKeys(ctx sdk.Context, prefix []byte, process func(key []byte) (stop bool)) {
	store := ctx.KVStore(tk.key)
	iter := sdk.KVStorePrefixIterator(store, prefix)
	defer iter.Close()
	for {
		if !iter.Valid() {
			return
		}
		key := iter.Key()
		if process(key) {
			return
		}
		iter.Next()
	}
}

func (tk TokenKeeper) setAddrKey(ctx sdk.Context, prefix []byte, addr string) error {
	store := ctx.KVStore(tk.key)
	index := strings.Index(addr, string(SeparateKeyPrefix))

	accBech32, err := sdk.AccAddressFromBech32(string([]byte(addr)[index+1:]))
	if err != nil {
		return err
	}
	key := PrefixAddrStoreKey(prefix, string([]byte(addr)[:index]), accBech32)
	store.Set(key, []byte{})

	return nil
}

func (tk TokenKeeper) GetReservedSymbols() []string {
	return reserved
}
