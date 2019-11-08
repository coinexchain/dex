package types

// only coinexdex owner can issue reserved symbol token
var reservedSymbolMap map[string]int

//nolint
var reserved = []string{
	// coin market cap currency Top 200
	"btc", "eth", "xrp", "bch", "usdt", "ltc", "eos", "bnb", "bsv", "xlm",
	"trx", "ada", "xmr", "leo", "ht", "link", "xtz", "neo", "miota", "atom",
	"mkr", "dash", "etc", "ont", "usdc", "cro", "xem", "bat", "doge", "vet",
	"zec", "pax", "hedg", "qtum", "dcr", "zrx", "tusd", "hot", "btg", "cennz",
	"vsys", "rvn", "omg", "zb", "btm", "nano", "luna", "rep", "snx", "abbc",
	"algo", "ekt", "kcs", "dai", "btt", "lsk", "bcd", "icx", "dgb", "sc",
	"hc", "qnt", "kmd", "waves", "bts", "bcn", "theta", "sxp", "kbc", "iost",
	"mona", "ftt", "ae", "mco", "dx", "xvg", "maid", "nexo", "seele", "zil",
	"nrg", "ardr", "aoa", "rlc", "chz", "enj", "steem", "rif", "elf", "snt",
	"ren", "npxs", "gnt", "xzc", "crpt", "new", "hpt", "ilc", "solve", "ode",
	"zen", "nex", "lamb", "gxc", "etn", "matic", "xmx", "mof", "eurs", "drg",
	"win", "bcv", "ela", "wtc", "mana", "strat", "dgtx", "nas", "aion", "etp",
	"brd", "grin", "beam", "pai", "nuls", "knc", "tt", "tnt", "lrc", "wicc",
	"cvc", "true", "ppt", "fct", "ignis", "rdd", "dgd", "ftm", "fet", "ark",
	"rcn", "gt", "wan", "fun", "tomo", "you", "ant", "r", "hbar", "eng",
	"waxp", "iotx", "loom", "lina", "hyn", "dcn", "orbs", "edc", "bhp", "busd",
	"qash", "man", "powr", "bnt", "c20", "dag", "dent", "storj", "mtl", "nxs",
	"uni", "san", "gno", "kan", "cocos", "rox", "stream", "grs", "abt", "bix",
	"cs", "cmt", "edo", "one", "tel", "gas", "fx", "medx", "adn", "celr",
	"sys", "wxt", "erd", "wgr", "lend", "mda", "dpt", "loki", "iris", "egt",

	// coinex exchange currency: 96
	"ont", "grin", "gnt", "atom", "qtum", "eth", "btc", "omg", "ae", "wings",
	"btm", "eosc", "dcr", "loom", "olt", "neo", "icx", "bat", "etc", "kmd",
	"xzc", "dash", "iota", "btt", "vet", "okb", "xtz", "link", "stx", "algo",
	"hot", "cnn", "tusd", "lsk", "ltc", "zil", "xrp", "zec", "vsys", "bch",
	"ada", "cmt", "cody", "nano", "zrx", "hydro", "trx", "cet", "seele", "doge",
	"dot", "ftt", "bsv", "hc", "btu", "iost", "sai", "xmr", "dero", "eos",
	"spice", "xvg", "spok", "ong", "rvn", "sys", "usdh", "xem", "dgb", "bnb",
	"nnb", "xlm", "waves", "ela", "ardr", "ult", "usdt", "usdc", "lamb", "sc",
	"kan", "ckb", "ht", "pax", "ctxc", "gusd", "gram", "trtl", "vtho", "akro",
	"whc", "lfc", "tct", "gas", "egt", "fcny",

	// ISO 8601 fiat currency: 93
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

	// precious metals: 4
	"xau", "xag", "xpt", "xpd",

	// extra currency
	"libra","cet","rmb","dcep","coinex","viabtc","matrix","bitmain",
}

func init() {
	reservedSymbolMap = make(map[string]int)
	for i, symbol := range reserved {
		reservedSymbolMap[symbol] = i
	}
}

func IsReservedSymbol(symbol string) bool {
	var _, found = reservedSymbolMap[symbol]
	return found
}

func GetReservedSymbols() []string {
	return reserved
}
