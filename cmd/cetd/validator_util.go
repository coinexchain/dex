package main

import (
	"fmt"
	"sort"

	tm "github.com/tendermint/tendermint/types"
)

type validators []validatorInfo
type validatorInfo struct {
	name  string
	power int64
}

func (vs validators) Len() int {
	return len(vs)
}
func (vs validators) Less(i, j int) bool {
	return vs[i].power < vs[j].power
}
func (vs validators) Swap(i, j int) {
	tmp := vs[i]
	vs[i] = vs[j]
	vs[j] = tmp
}

func listValidators(genDoc *tm.GenesisDoc) {
	totalPower := int64(0)
	vs := make([]validatorInfo, len(genDoc.Validators))
	for i, v := range genDoc.Validators {
		totalPower += v.Power
		vs[i] = validatorInfo{
			name:  v.Name,
			power: v.Power,
		}
	}
	sort.Sort(sort.Reverse(validators(vs)))
	for i, v := range vs {
		ratio := float64(v.power) / float64(totalPower) * 100.0
		fmt.Printf("#%02d\t%12d\t%5.2f%%\t%s\n",
			i+1, v.power, ratio, v.name)
	}
}
