package asset

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

func init() {
	reservedSymbolMap = make(map[string]int)
	for i, symbol := range reserved {
		reservedSymbolMap[symbol] = i
	}
}
