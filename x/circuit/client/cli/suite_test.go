package cli_test

import (
	"bytes"
	"io"
	"testing"

	"cosmossdk.io/x/circuit"
	abci "github.com/cometbft/cometbft/abci/types"
	rpcclientmock "github.com/cometbft/cometbft/rpc/client/mock"
	"github.com/stretchr/testify/suite"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/testutil"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	testutilmod "github.com/cosmos/cosmos-sdk/types/module/testutil"
)

type CLITestSuite struct {
	suite.Suite

	kr        keyring.Keyring
	baseCtx   client.Context
	encCfg    testutilmod.TestEncodingConfig
	clientCtx client.Context

	accounts []sdk.AccAddress
}

func TestCLITestSuite(t *testing.T) {
	suite.Run(t, new(CLITestSuite))
}

func (s *CLITestSuite) SetupSuite() {
	s.T().Log("setting up cli test suite")

	s.encCfg = testutilmod.MakeTestEncodingConfig(circuit.AppModuleBasic{})
	s.kr = keyring.NewInMemory(s.encCfg.Codec)
	s.baseCtx = client.Context{}.
		WithKeyring(s.kr).
		WithTxConfig(s.encCfg.TxConfig).
		WithCodec(s.encCfg.Codec).
		WithClient(clitestutil.MockCometRPC{Client: rpcclientmock.Client{}}).
		WithAccountRetriever(client.MockAccountRetriever{}).
		WithOutput(io.Discard).
		WithChainID("test-chain")

	var outBuf bytes.Buffer
	ctxGen := func() client.Context {
		bz, _ := s.encCfg.Codec.Marshal(&sdk.TxResponse{})
		c := clitestutil.NewMockCometRPC(abci.ResponseQuery{
			Value: bz,
		})

		return s.baseCtx.WithClient(c)
	}
	s.clientCtx = ctxGen().WithOutput(&outBuf)

	accounts := testutil.CreateKeyringAccounts(s.T(), s.kr, 4)

	for _, v := range accounts {
		s.accounts = append(s.accounts, v.Address)
	}
	s.accounts[1] = accounts[1].Address
}