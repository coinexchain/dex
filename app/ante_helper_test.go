package app

import (
	"reflect"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov"

	"github.com/coinexchain/cet-sdk/types"
)

var (
	testAddr, _ = sdk.AccAddressFromBech32("test-addr")
	abcCoins    = types.NewCoins("abc", 100)
	cetCoins    = types.NewCetCoins(1000)
	ah          = anteHelper{}
)

func TestAnteHelper_CheckMsgDeposit(t *testing.T) {

	tests := []struct {
		name string
		msg  gov.MsgDeposit
		want error
	}{
		{
			name: "deposit abc coins",
			msg:  gov.NewMsgDeposit(testAddr, 1, abcCoins),
			want: sdk.ErrInvalidCoins("tx not allowed to deposit other coins than cet"),
		},
		{
			name: "deposit cet coins",
			msg:  gov.NewMsgDeposit(testAddr, 1, cetCoins),
			want: nil,
		},
		{
			name: "deposit multiple coins",
			msg:  gov.NewMsgDeposit(testAddr, 1, sdk.NewCoins(abcCoins[0], cetCoins[0])),
			want: sdk.ErrInvalidCoins("tx not allowed to deposit other coins than cet"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ah.checkMsgDeposit(tt.msg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MsgIssueToken.ValidateBasic() = %v, want %v", got, tt.want)
			}
		})
	}
}
