package model

import "encoding/json"

type ScoreType = int16

type ScoreViewBonusRule struct {
	Days  int32
	Bonus ScoreType
}

func (rule *ScoreViewBonusRule) MarshalJSON() ([]byte, error) {
	return json.Marshal([]int64{
		int64(rule.Days),
		int64(rule.Bonus),
	})
}

func (rule *ScoreViewBonusRule) UnmarshalJSON(b []byte) error {
	var val []int64
	err := json.Unmarshal(b, &val)
	if err != nil {
		return err
	}

	rule.Days = int32(val[0])
	rule.Bonus = ScoreType(val[1])
	return nil
}
