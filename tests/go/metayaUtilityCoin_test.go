package test

import (
	"testing"

	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/stretchr/testify/assert"

	"github.com/0xRaito/metaya-smart-contracts/tests/go/metayaUtilityCoin"
	"github.com/0xRaito/metaya-smart-contracts/tests/go/test"
)

func TestMetayaUtilityCoinDeployContracts(t *testing.T) {
	b := test.NewBlockchain()

	fungibleTokenAddress, metayaUtilityCoinAddress, _ := metayaUtilityCoin.DeployContracts(t, b)

	t.Run("Should have initialized Supply field correctly", func(t *testing.T) {
		supply := test.ExecuteScriptAndCheck(
			t, b,
			metayaUtilityCoin.GetSupplyScript(fungibleTokenAddress.String(), metayaUtilityCoinAddress.String()),
			nil,
		)
		assert.EqualValues(t, test.CadenceUFix64("0.0"), supply)
	})
}

func TestMetayaUtilityCoinSetupAccount(t *testing.T) {
	b := test.NewBlockchain()

	t.Run("Should be able to create empty vault that does not affect supply", func(t *testing.T) {
		fungibleAddr, metayaUtilityCoinAddr, _ := metayaUtilityCoin.DeployContracts(t, b)

		userAddress, _ := metayaUtilityCoin.CreateAccount(t, b, fungibleAddr, metayaUtilityCoinAddr)

		userBalance := test.ExecuteScriptAndCheck(
			t, b,
			metayaUtilityCoin.GetBalanceScript(fungibleAddr.String(), metayaUtilityCoinAddr.String()),
			[][]byte{jsoncdc.MustEncode(cadence.Address(userAddress))},
		)
		assert.EqualValues(t, test.CadenceUFix64("0.0"), userBalance)

		supply := test.ExecuteScriptAndCheck(
			t, b,
			metayaUtilityCoin.GetSupplyScript(fungibleAddr.String(), metayaUtilityCoinAddr.String()),
			nil,
		)
		assert.EqualValues(t, test.CadenceUFix64("0.0"), supply)
	})
}

func TestMetayaUtilityCoinMintTokens(t *testing.T) {
	b := test.NewBlockchain()

	fungibleTokenAddress, metayaUtilityCoinAddress, metayaUtilityCoinSigner := metayaUtilityCoin.DeployContracts(t, b)

	userAddress, _ := metayaUtilityCoin.CreateAccount(t, b, fungibleTokenAddress, metayaUtilityCoinAddress)

	t.Run("Should not be able to mint zero tokens", func(t *testing.T) {
		metayaUtilityCoin.MintTokens(
			t, b,
			fungibleTokenAddress,
			metayaUtilityCoinAddress,
			metayaUtilityCoinSigner,
			userAddress,
			"0.0",
			true,
		)
	})

	t.Run("Should be able to mint tokens, deposit, update balance and total supply", func(t *testing.T) {

		metayaUtilityCoin.MintTokens(
			t, b,
			fungibleTokenAddress,
			metayaUtilityCoinAddress,
			metayaUtilityCoinSigner,
			userAddress,
			"50.0",
			false,
		)

		// Assert that vault balance is correct
		userBalance := test.ExecuteScriptAndCheck(
			t, b,
			metayaUtilityCoin.GetBalanceScript(fungibleTokenAddress.String(), metayaUtilityCoinAddress.String()),
			[][]byte{jsoncdc.MustEncode(cadence.Address(userAddress))},
		)

		assert.EqualValues(t, test.CadenceUFix64("50.0"), userBalance)

		// Assert that total supply is correct
		supply := test.ExecuteScriptAndCheck(
			t, b,
			metayaUtilityCoin.GetSupplyScript(fungibleTokenAddress.String(), metayaUtilityCoinAddress.String()),
			nil,
		)

		assert.EqualValues(t, test.CadenceUFix64("50.0"), supply)
	})
}

func TestMetayaUtilityCoinTransfersTokens(t *testing.T) {
	b := test.NewBlockchain()

	fungibleTokenAddress, metayaUtilityCoinAddress, metayaUtilityCoinSigner := metayaUtilityCoin.DeployContracts(t, b)

	userAddress, _ := metayaUtilityCoin.CreateAccount(t, b, fungibleTokenAddress, metayaUtilityCoinAddress)

	// Mint 1000 new MetayaUtilityCoin into the metayaUtilityCoin contract account
	metayaUtilityCoin.MintTokens(
		t, b,
		fungibleTokenAddress,
		metayaUtilityCoinAddress,
		metayaUtilityCoinSigner,
		metayaUtilityCoinAddress,
		"1000.0",
		false,
	)

	t.Run("Should not be able to withdraw more than the balance of the vault", func(t *testing.T) {
		tx := flow.NewTransaction().
			SetScript(metayaUtilityCoin.TransferTokensTransaction(fungibleTokenAddress.String(), metayaUtilityCoinAddress.String())).
			SetGasLimit(100).
			SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
			SetPayer(b.ServiceKey().Address).
			AddAuthorizer(metayaUtilityCoinAddress)

		_ = tx.AddArgument(test.CadenceUFix64("30000.0"))
		_ = tx.AddArgument(cadence.NewAddress(userAddress))

		test.SignAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, metayaUtilityCoinAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), metayaUtilityCoinSigner},
			true,
		)

		// Assert that vault balances are correct

		metayaUtilityCoinBalance := test.ExecuteScriptAndCheck(
			t, b,
			metayaUtilityCoin.GetBalanceScript(fungibleTokenAddress.String(), metayaUtilityCoinAddress.String()),
			[][]byte{jsoncdc.MustEncode(cadence.Address(metayaUtilityCoinAddress))},
		)

		assert.EqualValues(t, test.CadenceUFix64("1000.0"), metayaUtilityCoinBalance)

		userBalance := test.ExecuteScriptAndCheck(
			t, b,
			metayaUtilityCoin.GetBalanceScript(fungibleTokenAddress.String(), metayaUtilityCoinAddress.String()),
			[][]byte{jsoncdc.MustEncode(cadence.Address(userAddress))},
		)

		assert.EqualValues(t, test.CadenceUFix64("0.0"), userBalance)
	})

	t.Run("Should be able to withdraw and deposit tokens from a vault", func(t *testing.T) {
		tx := flow.NewTransaction().
			SetScript(metayaUtilityCoin.TransferTokensTransaction(fungibleTokenAddress.String(), metayaUtilityCoinAddress.String())).
			SetGasLimit(100).
			SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
			SetPayer(b.ServiceKey().Address).
			AddAuthorizer(metayaUtilityCoinAddress)

		_ = tx.AddArgument(test.CadenceUFix64("300.0"))
		_ = tx.AddArgument(cadence.NewAddress(userAddress))

		test.SignAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, metayaUtilityCoinAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), metayaUtilityCoinSigner},
			false,
		)

		// Assert that vault balances are correct

		metayaUtilityCoinBalance := test.ExecuteScriptAndCheck(
			t, b,
			metayaUtilityCoin.GetBalanceScript(fungibleTokenAddress.String(), metayaUtilityCoinAddress.String()),
			[][]byte{jsoncdc.MustEncode(cadence.Address(metayaUtilityCoinAddress))},
		)

		assert.EqualValues(t, test.CadenceUFix64("700.0"), metayaUtilityCoinBalance)

		userBalance := test.ExecuteScriptAndCheck(
			t, b,
			metayaUtilityCoin.GetBalanceScript(fungibleTokenAddress.String(), metayaUtilityCoinAddress.String()),
			[][]byte{jsoncdc.MustEncode(cadence.Address(userAddress))},
		)

		assert.EqualValues(t, test.CadenceUFix64("300.0"), userBalance)

		supply := test.ExecuteScriptAndCheck(
			t, b,
			metayaUtilityCoin.GetSupplyScript(fungibleTokenAddress.String(), metayaUtilityCoinAddress.String()),
			nil,
		)
		assert.EqualValues(t, test.CadenceUFix64("1000.0"), supply)
	})

	t.Run("Should be able to transfer to multiple accounts", func(t *testing.T) {

		recipient1Address := cadence.Address(userAddress)
		recipient1Amount := test.CadenceUFix64("300.0")

		pair := cadence.KeyValuePair{Key: recipient1Address, Value: recipient1Amount}
		recipientPairs := make([]cadence.KeyValuePair, 1)
		recipientPairs[0] = pair

		tx := flow.NewTransaction().
			SetScript(metayaUtilityCoin.TransferManyAccountTransaction(fungibleTokenAddress.String(), metayaUtilityCoinAddress.String())).
			SetGasLimit(100).
			SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
			SetPayer(b.ServiceKey().Address).
			AddAuthorizer(metayaUtilityCoinAddress)

		_ = tx.AddArgument(cadence.NewDictionary(recipientPairs))

		test.SignAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, metayaUtilityCoinAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), metayaUtilityCoinSigner},
			false,
		)

		// Assert that the vaults' balances are correct
		
		metayaUtilityCoinBalance := test.ExecuteScriptAndCheck(
			t, b,
			metayaUtilityCoin.GetBalanceScript(fungibleTokenAddress.String(), metayaUtilityCoinAddress.String()),
			[][]byte{jsoncdc.MustEncode(cadence.Address(metayaUtilityCoinAddress))},
		)

		assert.EqualValues(t, test.CadenceUFix64("400.0"), metayaUtilityCoinBalance)

		userBalance := test.ExecuteScriptAndCheck(
			t, b,
			metayaUtilityCoin.GetBalanceScript(fungibleTokenAddress.String(), metayaUtilityCoinAddress.String()),
			[][]byte{jsoncdc.MustEncode(cadence.Address(userAddress))},
		)

		assert.EqualValues(t, test.CadenceUFix64("600.0"), userBalance)

		supply := test.ExecuteScriptAndCheck(
			t, b,
			metayaUtilityCoin.GetSupplyScript(fungibleTokenAddress.String(), metayaUtilityCoinAddress.String()),
			nil,
		)
		assert.EqualValues(t, test.CadenceUFix64("1000.0"), supply)
	})
}

func TestMetayaUtilityCoinBurnTokens(t *testing.T) {
	b := test.NewBlockchain()

	fungibleTokenAddress, metayaUtilityCoinAddress, metayaUtilityCoinSigner := metayaUtilityCoin.DeployContracts(t, b)

	userAddress, userSigner := metayaUtilityCoin.CreateAccount(t, b, fungibleTokenAddress, metayaUtilityCoinAddress)

	// Mint 1000 new MetayaUtilityCoin into the metayaUtilityCoin contract account
	metayaUtilityCoin.MintTokens(
		t, b,
		fungibleTokenAddress,
		metayaUtilityCoinAddress,
		metayaUtilityCoinSigner,
		metayaUtilityCoinAddress,
		"1000.0",
		false,
	)

	// Transfer 400 MetayaUtilityCoin to the userAddress
	tx := flow.NewTransaction().
			SetScript(metayaUtilityCoin.TransferTokensTransaction(fungibleTokenAddress.String(), metayaUtilityCoinAddress.String())).
			SetGasLimit(100).
			SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
			SetPayer(b.ServiceKey().Address).
			AddAuthorizer(metayaUtilityCoinAddress)

		_ = tx.AddArgument(test.CadenceUFix64("400.0"))
		_ = tx.AddArgument(cadence.NewAddress(userAddress))

		test.SignAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, metayaUtilityCoinAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), metayaUtilityCoinSigner},
			false,
		)

	t.Run("Should subtract tokens from supply when they are destroyed by admin", func(t *testing.T) {
		tx := flow.NewTransaction().
			SetScript(metayaUtilityCoin.BurnTokensByAdminTransaction(fungibleTokenAddress.String(), metayaUtilityCoinAddress.String())).
			SetGasLimit(100).
			SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
			SetPayer(b.ServiceKey().Address).
			AddAuthorizer(metayaUtilityCoinAddress)

		_ = tx.AddArgument(test.CadenceUFix64("300.0"))

		test.SignAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, metayaUtilityCoinAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), metayaUtilityCoinSigner},
			false,
		)

		// Assert that the vaults' balances are correct
		
		metayaUtilityCoinBalance := test.ExecuteScriptAndCheck(
			t, b,
			metayaUtilityCoin.GetBalanceScript(fungibleTokenAddress.String(), metayaUtilityCoinAddress.String()),
			[][]byte{jsoncdc.MustEncode(cadence.Address(metayaUtilityCoinAddress))},
		)

		assert.EqualValues(t, test.CadenceUFix64("300.0"), metayaUtilityCoinBalance)

		userBalance := test.ExecuteScriptAndCheck(
			t, b,
			metayaUtilityCoin.GetBalanceScript(fungibleTokenAddress.String(), metayaUtilityCoinAddress.String()),
			[][]byte{jsoncdc.MustEncode(cadence.Address(userAddress))},
		)

		assert.EqualValues(t, test.CadenceUFix64("400.0"), userBalance)

		supply := test.ExecuteScriptAndCheck(
			t, b,
			metayaUtilityCoin.GetSupplyScript(fungibleTokenAddress.String(), metayaUtilityCoinAddress.String()),
			nil,
		)
		assert.EqualValues(t, test.CadenceUFix64("700.0"), supply)
	})

	t.Run("Should subtract tokens from supply when they are destroyed by user", func(t *testing.T) {
		tx := flow.NewTransaction().
		SetScript(metayaUtilityCoin.BurnTokensByUserTransaction(fungibleTokenAddress.String(), metayaUtilityCoinAddress.String())).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(userAddress).
		AddAuthorizer(metayaUtilityCoinAddress)

	_ = tx.AddArgument(test.CadenceUFix64("200.0"))

	test.SignAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, userAddress, metayaUtilityCoinAddress},
		[]crypto.Signer{b.ServiceKey().Signer(), userSigner, metayaUtilityCoinSigner},
		false,
	)

	// Assert that the vaults' balances are correct
	
	metayaUtilityCoinBalance := test.ExecuteScriptAndCheck(
		t, b,
		metayaUtilityCoin.GetBalanceScript(fungibleTokenAddress.String(), metayaUtilityCoinAddress.String()),
		[][]byte{jsoncdc.MustEncode(cadence.Address(metayaUtilityCoinAddress))},
	)

	assert.EqualValues(t, test.CadenceUFix64("300.0"), metayaUtilityCoinBalance)

	userBalance := test.ExecuteScriptAndCheck(
		t, b,
		metayaUtilityCoin.GetBalanceScript(fungibleTokenAddress.String(), metayaUtilityCoinAddress.String()),
		[][]byte{jsoncdc.MustEncode(cadence.Address(userAddress))},
	)

	assert.EqualValues(t, test.CadenceUFix64("200.0"), userBalance)

	supply := test.ExecuteScriptAndCheck(
		t, b,
		metayaUtilityCoin.GetSupplyScript(fungibleTokenAddress.String(), metayaUtilityCoinAddress.String()),
		nil,
	)
	assert.EqualValues(t, test.CadenceUFix64("500.0"), supply)
	})
}
