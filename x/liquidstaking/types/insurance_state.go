package types

import "fmt"

func (is InsuranceState) Equal(is2 InsuranceState) bool {
	return is.TotalInsuranceTokens.Equal(is2.TotalInsuranceTokens) &&
		is.TotalPairedInsuranceTokens.Equal(is2.TotalPairedInsuranceTokens) &&
		is.TotalUnpairingInsuranceTokens.Equal(is2.TotalUnpairingInsuranceTokens) &&
		is.TotalRemainingInsuranceCommissions.Equal(is2.TotalRemainingInsuranceCommissions)
}

func (is InsuranceState) IsZeroState() bool {
	return is.TotalPairedInsuranceTokens.IsZero() &&
		// Total insurances includes Pairing insurances, so we should skip this
		// nas.TotalInsuranceTokens.IsZero() &&
		is.TotalUnpairingInsuranceTokens.IsZero() &&
		is.TotalRemainingInsuranceCommissions.IsZero()
}

func (is InsuranceState) String() string {
	// Print all fields with field name
	return fmt.Sprintf(`InsuranceState:
	TotalInsuranceTokens:       %s
	TotalPairedInsuranceTokens: %s
    TotalUnpairingInsuranceTokens: %s
    TotalRemainingInsuranceCommissions: %s
`,
		is.TotalInsuranceTokens,
		is.TotalPairedInsuranceTokens,
		is.TotalUnpairingInsuranceTokens,
		is.TotalRemainingInsuranceCommissions,
	)
}
