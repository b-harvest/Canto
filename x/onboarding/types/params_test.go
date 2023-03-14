package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParamsMisc(t *testing.T) {
	params := DefaultParams()
	require.NotEmpty(t, params.ParamSetPairs())
	kt := ParamKeyTable()
	require.NotEmpty(t, kt)
}

func TestParamsValidate(t *testing.T) {
	testCases := []struct {
		name     string
		params   Params
		expError bool
	}{
		{
			"empty params",
			Params{},
			false,
		},
		{
			"default params",
			DefaultParams(),
			false,
		},
		{
			"custom params",
			NewParams(true),
			false,
		},
		{
			"invalid duration",
			NewParams(true),
			true,
		},
	}

	for _, tc := range testCases {
		err := tc.params.Validate()
		if tc.expError {
			require.Error(t, err, tc.name)
		} else {
			require.NoError(t, err, tc.name)
		}
	}
}

func TestValidate(t *testing.T) {
	require.Error(t, validateBool(""))
	require.NoError(t, validateBool(true))
}
