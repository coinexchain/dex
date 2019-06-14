package msgqueue

type GenesisState struct {
	Brokers string `json:"brokers"`
	Topics  string `json:"topics"`
}

func NewGenesisState(brokers string, topics string) GenesisState {
	return GenesisState{
		Brokers: brokers,
		Topics:  topics,
	}
}

// DefaultGenesisState - Return a default genesis state
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Brokers: "",
		Topics:  "",
	}
}

func InitGenesis(k *Producer, data GenesisState) {
	k.SetParam(data)
}

func ExportGenesis(k Producer) GenesisState {
	return k.GetParam()
}
