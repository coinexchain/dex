package app

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/mock"
	"github.com/tendermint/tendermint/proxy"
	sm "github.com/tendermint/tendermint/state"
	"github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
	"github.com/pkg/profile"

	"github.com/coinexchain/dex/modules/bankx"
	"github.com/coinexchain/dex/modules/market"
	"github.com/coinexchain/dex/msgqueue"
	"github.com/coinexchain/dex/testutil"
	dex "github.com/coinexchain/dex/types"
)

var gCdc = codec.New()

func initCodec() {
	gCdc.RegisterInterface((*crypto.PubKey)(nil), nil)
	gCdc.RegisterInterface((*crypto.PrivKey)(nil), nil)
	gCdc.RegisterInterface((*sdk.Msg)(nil), nil)

	gCdc.RegisterConcrete(secp256k1.PubKeySecp256k1{}, "tendermint/PubKeySecp256k1", nil)
	gCdc.RegisterConcrete(secp256k1.PrivKeySecp256k1{}, "tendermint/PrivKeySecp256k1", nil)
	gCdc.RegisterConcrete(auth.StdTx{}, "cosmos-sdk/StdTx", nil)
	gCdc.RegisterConcrete(bankx.MsgSend{}, "bankx/MsgSend", nil)
	gCdc.RegisterConcrete(market.MsgCreateTradingPair{}, "market/MsgCreateTradingPair", nil)
	gCdc.RegisterConcrete(market.MsgCreateOrder{}, "market/MsgCreateOrder", nil)

}

func generateAccount(accountNum int) ([]auth.BaseAccount, []crypto.PrivKey) {
	res := make([]auth.BaseAccount, 0, accountNum)
	keys := make([]crypto.PrivKey, 0, accountNum)
	coins := sdk.NewCoins(sdk.NewInt64Coin("cet", 260000000000000), sdk.NewInt64Coin("eth", 100000000000))
	for i := 0; i < accountNum; i++ {
		key, _, fromAddr := testutil.KeyPubAddr()
		acc := auth.BaseAccount{Address: fromAddr, Coins: coins}
		res = append(res, acc)
		keys = append(keys, key)
	}
	return res, keys
}

func makeState(nVals, height int, stateDB dbm.DB) (sm.State, map[string]types.PrivValidator) {
	vals := make([]types.GenesisValidator, nVals)
	privVals := make(map[string]types.PrivValidator, nVals)
	for i := 0; i < nVals; i++ {
		secret := []byte(fmt.Sprintf("test%d", i))
		pk := ed25519.GenPrivKeyFromSecret(secret)
		valAddr := pk.PubKey().Address()
		vals[i] = types.GenesisValidator{
			Address: valAddr,
			PubKey:  pk.PubKey(),
			Power:   1000,
			Name:    fmt.Sprintf("test%d", i),
		}
		privVals[valAddr.String()] = types.NewMockPVWithParams(pk, false, false)
	}
	s, _ := sm.MakeGenesisState(&types.GenesisDoc{
		ChainID:    "execution_chain",
		Validators: vals,
		AppHash:    nil,
	})

	sm.SaveState(stateDB, s)

	for i := 1; i < height; i++ {
		s.LastBlockHeight++
		s.LastValidators = s.Validators.Copy()
		sm.SaveState(stateDB, s)
	}
	return s, privVals
}

func TestBlockExec(t *testing.T) {
	initCodec()
	var AccountNum = 4

	stateDB := dbm.NewDB("testBlock", dbm.GoLevelDBBackend, "./")
	defer func() {
		err := os.RemoveAll("./testBlock.db")
		require.Nil(t, err)
	}()

	// init account and app
	accs, keys := generateAccount(AccountNum)
	app := initAppWithBaseAccounts(accs...)
	app.msgQueProducer = msgqueue.NewProducerFromConfig([]string{"nop"}, "auth,bank", false, nil)
	proxyApp := proxy.NewAppConns(proxy.NewLocalClientCreator(app))
	err := proxyApp.Start()
	require.Nil(t, err)
	defer proxyApp.Stop()

	blockExec := sm.NewBlockExecutor(stateDB, log.TestingLogger(), proxyApp.Consensus(),
		mock.Mempool{}, sm.MockEvidencePool{})

	// generate txs
	txcount := 3000
	txs := make([]types.Tx, 0, txcount)
	for i := 0; i < txcount; i++ {
		//fromAddr, toAddr := accs[i%AccountNum].Address, accs[(i+1)%AccountNum].Address
		fromAddr, toAddr := accs[0].Address, accs[0].Address
		coins := dex.NewCetCoins(200000000)
		msg := bankx.NewMsgSend(fromAddr, toAddr, coins, 0)
		tx := newStdTxBuilder().
			Msgs(msg).GasAndFee(100000, 2000000).
			AccNumSeqKey(0, uint64(i), keys[0]).Build()
		txBytes, err := gCdc.MarshalBinaryLengthPrefixed(tx)
		require.Nil(t, err)
		txs = append(txs, txBytes)
	}

	// generate block
	state, _ := makeState(1, 1, stateDB)
	block, _ := state.MakeBlock(1, txs, new(types.Commit), nil, state.Validators.GetProposer().Address)
	require.Nil(t, err)
	require.NotNil(t, block)
	blockID := types.BlockID{Hash: block.Hash(), PartsHeader: block.MakePartSet(types.BlockPartSizeBytes).Header()}

	// exec block
	defer profile.Start(profile.CPUProfile).Stop()
	now := time.Now()
	state, err = blockExec.ApplyBlock(state, blockID, block)
	fmt.Println("exec block time :  ", time.Now().Sub(now).Seconds())
	fmt.Println("exec block time :  ", time.Now().Sub(now).Milliseconds())
	require.Nil(t, err)
}
