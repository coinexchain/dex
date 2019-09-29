package stakingx

import (
	"fmt"
	"testing"

	yaml "gopkg.in/yaml.v2"

	"github.com/stretchr/testify/require"

	"github.com/tendermint/tendermint/crypto/ed25519"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

// https://trello.com/c/ggorEHXt
func TestPubKeyToYAML(t *testing.T) {
	var pubKey = ed25519.PubKeyEd25519{153, 85, 161, 9, 255, 2, 202, 213, 158, 38, 76, 202, 196, 254, 102, 184, 243, 189, 86, 185, 13, 109, 80, 89, 2, 82, 100, 28, 166, 53, 122, 182}
	resp := sdk.TxResponse{
		Tx: auth.StdTx{
			Msgs: []sdk.Msg{
				staking.MsgCreateValidator{
					PubKey: pubKey,
				},
			},
			Signatures: []auth.StdSignature{
				{PubKey: pubKey},
			},
		},
	}

	var toPrint fmt.Stringer = resp
	out, err := yaml.Marshal(&toPrint)
	require.NoError(t, err)
	require.Equal(t, `height: 0
txhash: ""
code: 0
data: ""
rawlog: ""
logs: []
info: ""
gaswanted: 0
gasused: 0
events: []
codespace: ""
tx:
  msg:
  - description:
      moniker: ""
      identity: ""
      website: ""
      details: ""
    commission:
      rate: <nil>
      max_rate: <nil>
      max_change_rate: <nil>
    min_self_delegation: <nil>
    delegator_address: ""
    validator_address: ""
    pubkey:
    - 153
    - 85
    - 161
    - 9
    - 255
    - 2
    - 202
    - 213
    - 158
    - 38
    - 76
    - 202
    - 196
    - 254
    - 102
    - 184
    - 243
    - 189
    - 86
    - 185
    - 13
    - 109
    - 80
    - 89
    - 2
    - 82
    - 100
    - 28
    - 166
    - 53
    - 122
    - 182
    value:
      denom: ""
      amount: <nil>
  fee:
    amount: []
    gas: 0
  signatures:
  - |
    pubkey: cosmospub1zcjduepqn926zz0lqt9dt83xfn9vflnxhrem644ep4k4qkgz2fjpef3402mqj54wjh
    signature: ""
  memo: ""
timestamp: ""
`, string(out))
}
