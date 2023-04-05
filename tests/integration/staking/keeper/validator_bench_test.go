package keeper_test

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktestutil "github.com/cosmos/cosmos-sdk/x/bank/testutil"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

func BenchmarkGetValidator(b *testing.B) {
	// 900 is the max number we are allowed to use in order to avoid simtestutil.CreateTestPubKeys
	// panic: encoding/hex: odd length hex string
	powersNumber := 900

	var totalPower int64
	powers := make([]int64, powersNumber)
	for i := range powers {
		powers[i] = int64(i)
		totalPower += int64(i)
	}

	app, ctx, _, valAddrs, vals := initValidators(b, totalPower, len(powers), powers)

	for _, validator := range vals {
		app.StakingKeeper.SetValidator(ctx, validator)
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for _, addr := range valAddrs {
			_, _ = app.StakingKeeper.GetValidator(ctx, addr)
		}
	}
}

func BenchmarkGetValidatorDelegations(b *testing.B) {
	powersNumber := 10

	var totalPower int64
	powers := make([]int64, powersNumber)
	for i := range powers {
		powers[i] = int64(i)
		totalPower += int64(i)
	}

	app, ctx, _, valAddrs, vals := initValidators(b, totalPower, len(powers), powers)

	for _, validator := range vals {
		app.StakingKeeper.SetValidator(ctx, validator)
	}

	delegationsNum := 1000

	delegators := make([]sdk.AccAddress, delegationsNum)

	for _, val := range valAddrs {
		for i := 0; i < delegationsNum; i++ {
			delegator := sdk.AccAddress(fmt.Sprintf("address%d", i))

			banktestutil.FundAccount(app.BankKeeper, ctx, delegator,
				sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(int64(i)))))

			NewDel := types.NewDelegation(delegator, val, sdk.NewDec(int64(i)))
			app.StakingKeeper.SetDelegation(ctx, NewDel)
			delegators[i] = delegator
		}
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		app.StakingKeeper.UpdateValidatorDelegations(ctx, valAddrs[0], sdk.ValAddress("val"))
	}
}

func BenchmarkGetValidatorDelegationsOld(b *testing.B) {

	powersNumber := 10

	var totalPower int64
	powers := make([]int64, powersNumber)
	for i := range powers {
		powers[i] = int64(i)
		totalPower += int64(i)
	}

	app, ctx, _, valAddrs, vals := initValidators(b, totalPower, len(powers), powers)

	for _, validator := range vals {
		app.StakingKeeper.SetValidator(ctx, validator)
	}

	delegationsNum := 1000

	delegators := make([]sdk.AccAddress, delegationsNum)

	for _, val := range valAddrs {
		for i := 0; i < delegationsNum; i++ {
			delegator := sdk.AccAddress(fmt.Sprintf("address%d", i))

			banktestutil.FundAccount(app.BankKeeper, ctx, delegator,
				sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(int64(i)))))

			NewDel := types.NewDelegation(delegator, val, sdk.NewDec(int64(i)))
			app.StakingKeeper.SetDelegation(ctx, NewDel)
			delegators[i] = delegator
		}
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		app.StakingKeeper.UpdateValidatorDelegationsOld(ctx, valAddrs[0], sdk.ValAddress("val"))
	}
}
