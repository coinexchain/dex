package asset

import (
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// only coinexdex owner can issue reserved symbol token
var reservedSymbolMap map[string]int

//nolint
var reserved = []string{
	// coin market cap currency 200
	"btc", "eth", "xrp", "ltc", "bch", "eos", "bnb", "bsv", "usdt", "xlm",
	"ada", "trx", "xmr", "dash", "miota", "atom", "etc", "neo", "xtz", "xem",
	"mkr", "ont", "zec", "btg", "cro", "vet", "bat", "link", "doge", "usdc",
	"qtum", "omg", "dcr", "btt", "tusd", "hot", "rvn", "waves", "lsk", "bcd",
	"rep", "nano", "npxs", "zil", "zrx", "kmd", "pax", "icx", "bts", "bcn",
	"dgb", "ht", "aoa", "mona", "gxc", "xvg", "btm", "ae", "iost", "dent",
	"steem", "solve", "sc", "qbit", "theta", "enj", "etp", "abbc", "hc", "snt",
	"thr", "elf", "ardr", "kcs", "mco", "strat", "wtc", "wax", "gnt", "nas",
	"maid", "hedg", "dai", "inb", "xin", "cnx", "mxm", "true", "pai", "vest",
	"xzc", "nuls", "ark", "zen", "lrc", "dgd", "loom", "aion", "mana", "fct",
	"ppt", "ela", "nex", "orbs", "nexo", "matic", "san", "new", "cccx", "tfuel",
	"ode", "r", "celr", "egt", "powr", "wicc", "la", "moac", "net", "rdd",
	"wan", "knc", "ipc", "etn", "ignis", "fun", "bnt", "ftm", "lamb", "poly",
	"nrg", "ecoreal", "bczero", "qash", "pivx", "eng", "brd", "fsn", "ekt", "abt",
	"qkc", "storj", "iotx", "grin", "etz", "ren", "eurs", "snx", "sys", "nxt",
	"rlc", "qnt", "repo", "meta", "veri", "dgtx", "rif", "gas", "pay", "grs",
	"jct", "ctxc", "csc", "itc", "rhoc", "cvc", "cmt", "c20", "bix", "mft",
	"tomo", "utk", "mith", "emc2", "vtc", "mhc", "cpt", "ugas", "part", "cnd",
	"agi", "cosm", "ino", "apl", "gno", "sky", "nkn", "gbyte", "medx", "sxdt",
	"edo", "ttc", "cennz", "dtr", "gbc", "noah", "nxs", "icn", "lba", "dac",

	// ISO 8601 fiat currency 93
	"usd", "all", "dzd", "ars", "amd", "aud", "azn", "bhd", "bdt", "byn",
	"bmd", "bob", "bam", "brl", "bgn", "khr", "cad", "clp", "cny", "cop",
	"crc", "hrk", "cup", "czk", "dkk", "dop", "egp", "eur", "gel", "ghs",
	"gtq", "hnl", "hkd", "huf", "isk", "inr", "idr", "irr", "iqd", "ils",
	"jmd", "jpy", "jod", "kzt", "kes", "kwd", "kgs", "lbp", "mkd", "myr",
	"mur", "mxn", "mdl", "mnt", "mad", "mmk", "nad", "npr", "twd", "nzd",
	"nio", "ngn", "nok", "omr", "pkr", "pab", "pen", "php", "pln", "gbp",
	"qar", "ron", "rub", "sar", "rsd", "sgd", "zar", "krw", "ssp", "ves",
	"lkr", "sek", "chf", "thb", "ttd", "tnd", "try", "ugx", "uah", "aed",
	"uyu", "uzs", "vnd",

	"libra",

	// precious metals 4
	"xau", "xag", "xpt", "xpd",
}

// GenesisState - all asset state that must be provided at genesis
type GenesisState struct {
	Params     Params   `json:"params"`
	Tokens     []Token  `json:"tokens"`
	Whitelist  []string `json:"whitelist"`
	ForbidAddr []string `json:"forbid_addr"`
}

func init() {
	reservedSymbolMap = make(map[string]int)
	for i, symbol := range reserved {
		reservedSymbolMap[symbol] = i
	}
}

// NewGenesisState - Create a new genesis state
func NewGenesisState(params Params, tokens []Token, whitelist []string, forbidAddr []string) GenesisState {
	return GenesisState{
		Params:     params,
		Tokens:     tokens,
		Whitelist:  whitelist,
		ForbidAddr: forbidAddr,
	}
}

// DefaultGenesisState - Return a default genesis state
func DefaultGenesisState() GenesisState {
	return NewGenesisState(DefaultParams(), []Token{}, []string{}, []string{})
}

// InitGenesis - Init store state from genesis data
func InitGenesis(ctx sdk.Context, keeper BaseKeeper, data GenesisState) {
	keeper.SetParams(ctx, data.Params)

	for _, token := range data.Tokens {
		if err := keeper.setToken(ctx, token); err != nil {
			panic(err)
		}
	}
	for _, addr := range data.Whitelist {
		if err := keeper.setAddrKey(ctx, WhitelistKeyPrefix, addr); err != nil {
			panic(err)
		}
	}
	for _, addr := range data.ForbidAddr {
		if err := keeper.setAddrKey(ctx, ForbidAddrKeyPrefix, addr); err != nil {
			panic(err)
		}
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper
func ExportGenesis(ctx sdk.Context, keeper BaseKeeper) GenesisState {
	return NewGenesisState(
		keeper.GetParams(ctx),
		keeper.GetAllTokens(ctx),
		keeper.GetAllAddrKeys(ctx, WhitelistKeyPrefix),
		keeper.GetAllAddrKeys(ctx, ForbidAddrKeyPrefix))
}

// ValidateGenesis performs basic validation of asset genesis data returning an
// error for any failed validation criteria.
func (data GenesisState) Validate() error {
	if err := data.Params.ValidateGenesis(); err != nil {
		return err
	}

	for _, token := range data.Tokens {
		if err := token.Validate(); err != nil {
			return err
		}
	}

	tokenSymbols := make(map[string]interface{})
	for _, token := range data.Tokens {
		if _, exists := tokenSymbols[token.GetSymbol()]; exists {
			return errors.New("duplicate token symbol found in GenesisState")
		}

		tokenSymbols[token.GetSymbol()] = nil
	}

	return nil
}
