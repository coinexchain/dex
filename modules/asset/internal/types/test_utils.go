package types

import (
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _, _, testAddr = keyPubAddr()

func keyPubAddr() (crypto.PrivKey, crypto.PubKey, sdk.AccAddress) {
	key := secp256k1.GenPrivKey()
	pub := key.PubKey()
	addr := sdk.AccAddress(pub.Address())
	return key, pub, addr
}

var TestIdentityString = "552A83BA62F9B1F8"

func mockAddrList() (list []sdk.AccAddress) {
	var addr1, _ = sdk.AccAddressFromBech32("coinex1y5kdxnzn2tfwayyntf2n28q8q2s80mcul852ke")
	var addr2, _ = sdk.AccAddressFromBech32("coinex133w8vwj73s4h2uynqft9gyyy52cr6rg8dskv3h")
	var addr3, _ = sdk.AccAddressFromBech32("coinex1zvf0hx6rpz0n7dkuzu34s39dnsyr8eygqs8h3q")

	list = append(list, addr1)
	list = append(list, addr2)
	list = append(list, addr3)
	return
}

var nilAddr, _ = sdk.AccAddressFromBech32("")
