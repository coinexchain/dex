package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/dex/types"
)

// RouterKey is the name of the bankx module
const (
	ModuleName = "distrx"
	RouterKey  = ModuleName
)

var _ sdk.Msg = MsgDonateToCommunityPool{}

// msg struct for validator withdraw
type MsgDonateToCommunityPool struct {
	FromAddr sdk.AccAddress `json:"from_addr"`
	Amount   sdk.Coins      `json:"amount"`
}

func NewMsgDonateToCommunityPool(addr sdk.AccAddress, amt sdk.Coins) MsgDonateToCommunityPool {
	return MsgDonateToCommunityPool{
		FromAddr: addr,
		Amount:   amt,
	}
}

func (msg MsgDonateToCommunityPool) Route() string { return ModuleName }
func (msg MsgDonateToCommunityPool) Type() string  { return "donate_to_community_pool" }

// Return address that must sign over msg.GetSignBytes()
func (msg MsgDonateToCommunityPool) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.FromAddr.Bytes())}
}

// get the bytes for the message signer to sign on
func (msg MsgDonateToCommunityPool) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// quick validity check
func (msg MsgDonateToCommunityPool) ValidateBasic() sdk.Error {
	if msg.FromAddr.Empty() {
		return ErrorInvalidFromAddr()
	}

	if msg.Amount.Len() != 1 {
		return ErrorInvalidDonation("invalid donation length")
	}
	if msg.Amount[0].Denom != types.DefaultBondDenom {
		return ErrorInvalidDonation("donation's denom must be cet")
	}

	if !msg.Amount.IsValid() {
		return ErrorInvalidDonation("invalid donation")
	}
	return nil
}
