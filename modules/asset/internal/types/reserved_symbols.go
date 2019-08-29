package types

// only coinexdex owner can issue reserved symbol token
var reservedSymbolMap map[string]int

//nolint
var reserved = []string{
	// coin market cap currency Top 200
	"btc", "eth", "xrp", "bch", "ltc", "usdt", "bnb", "eos", "bsv", "xlm",
	"xmr", "ada", "leo", "ht", "trx", "dash", "etc", "xtz", "miota", "link",
	"neo", "mkr", "usdc", "xem", "atom", "ont", "cro", "zec", "doge", "vsys",
	"hedg", "dcr", "vet", "bat", "qtum", "pax", "btg", "tusd", "zb", "omg",
	"lsk", "rvn", "nano", "kcs", "bcd", "inb", "algo", "waves", "btt", "nrg",
	"hot", "icx", "theta", "bcn", "lamb", "dgb", "bts", "zrx", "hc", "npxs",
	"rep", "iost", "aoa", "egt", "maid", "mona", "nex", "kmd", "dai", "qnt",
	"btm", "sc", "xvg", "gnt", "rif", "ae", "zil", "etp", "steem", "snt",
	"ardr", "ren", "mco", "wtc", "enj", "xzc", "abbc", "gxc", "snx", "solve",
	"wax", "strat", "beam", "ela", "nexo", "elf", "grin", "zen", "pai", "eurs",
	"ftm", "r", "wan", "fct", "etn", "ode", "new", "dent", "qash", "rdd",
	"mana", "nuls", "nas", "dgd", "tomo", "lrc", "dgtx", "crpt", "san", "knc",
	"matic", "loom", "wicc", "ppt", "fsn", "qkc", "lina", "eng", "aion", "fet",
	"bix", "true", "ignis", "orbs", "ark", "one", "bnt", "bhp", "tel", "fx",
	"powr", "xmx", "brd", "bcv", "cmt", "ttc", "celr", "c20", "valor", "tfuel",
	"storj", "bmc", "edo", "rhoc", "iotx", "mtl", "ant", "pivx", "hyn", "abt",
	"seele", "pzm", "fun", "agi", "gno", "divi", "rlc", "ret", "gbyte", "poly",
	"nxt", "bft", "gas", "grs", "dac", "ugas", "rox", "nxs", "sys", "cs",
	"cvnt", "ptt", "cos", "cvc", "ctxc", "vtc", "pla", "itc", "cbt", "gmat",
	"noah", "tkn", "tnt", "medx", "kan", "mith", "hpb", "erd", "trio", "man",

	// coinex exchange currency: 90
	"ont", "grin", "gnt", "atom", "qtum", "eth", "btc", "ae", "wings", "btm",
	"dcr", "olt", "neo", "icx", "bat", "etc", "kmd", "xzc", "dash", "iota",
	"btt", "vet", "okb", "stx", "algo", "link", "cnn", "tusd", "lsk", "ltc",
	"omg", "zil", "xrp", "zec", "hot", "bch", "ada", "cody", "nano", "zrx",
	"hydro", "trx", "cet", "seele", "doge", "dot", "ftt", "bsv", "hc", "btu",
	"xmr", "dero", "eos", "xtz", "spice", "xvg", "ong", "rvn", "sys", "usdh",
	"loom", "xem", "dgb", "bnb", "nnb", "xlm", "waves", "ult", "usdt", "usdc",
	"lamb", "wwb", "kan", "ckb", "ht", "pax", "ctxc", "gusd", "gram", "cmt",
	"trtl", "vtho", "akro", "whc", "sc", "lfc", "tct", "gas", "egt", "fcny",

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
	"libra", "cet", "rmb",
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
