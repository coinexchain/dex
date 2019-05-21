package incentive

type Keeper struct {
	feeCollectionKeeper FeeCollectionKeeper
	bankKeeper          BankKeeper
}

func NewKeeper(fck FeeCollectionKeeper, bk BankKeeper) Keeper {

	return Keeper{
		feeCollectionKeeper: fck,
		bankKeeper:          bk,
	}
}
