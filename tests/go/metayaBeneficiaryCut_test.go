package test

import (
	"testing"

	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/stretchr/testify/assert"

	"github.com/0xRaito/metaya-smart-contracts/tests/go/metayaBeneficiaryCut"
	"github.com/0xRaito/metaya-smart-contracts/tests/go/metayaUtilityCoin"
	"github.com/0xRaito/metaya-smart-contracts/tests/go/test"
)

func TestMetayaBeneficiaryCutDeployContracts(t *testing.T) {
	b := test.NewBlockchain()

	_, _, _, metayaBeneficiaryCutAddress, _ := metayaBeneficiaryCut.DeployContracts(t, b)

	t.Run("Should have initialized field correctly", func(t *testing.T) {
		metayaAddress := test.ExecuteScriptAndCheck(
			t, b,
			metayaBeneficiaryCut.GetMetayaAddressScript(metayaBeneficiaryCutAddress.String()),
			nil,
		)
		assert.EqualValues(t, cadence.Address(metayaBeneficiaryCutAddress), metayaAddress)

		metayaMarketCutPercentage := test.ExecuteScriptAndCheck(
			t, b,
			metayaBeneficiaryCut.GetMetayaMarketCutPercentageScript(metayaBeneficiaryCutAddress.String()),
			nil,
		)
		assert.EqualValues(t, test.CadenceUFix64("0.03"), metayaMarketCutPercentage)
	})
}

func TestMetayaBeneficiaryCut(t *testing.T) {
	b := test.NewBlockchain()

	fungibleTokenAddress, metayaUtilityCoinAddress, _, metayaBeneficiaryCutAddress, metayaBeneficiaryCutSigner := metayaBeneficiaryCut.DeployContracts(t, b)

	copyrightOwner1Address, copyrightOwner1Signer, _ := test.CreateAccount(t, b)
	metayaUtilityCoin.SetupAccount(
		t, b,
		copyrightOwner1Address,
		copyrightOwner1Signer,
		fungibleTokenAddress,
		metayaUtilityCoinAddress,
	)

	copyrightOwner2Address, copyrightOwner2Signer, _ := test.CreateAccount(t, b)
	metayaUtilityCoin.SetupAccount(
		t, b,
		copyrightOwner2Address,
		copyrightOwner2Signer,
		fungibleTokenAddress,
		metayaUtilityCoinAddress,
	)

	commonweal1Address, commonweal1Signer, _ := test.CreateAccount(t, b)
	metayaUtilityCoin.SetupAccount(
		t, b,
		commonweal1Address,
		commonweal1Signer,
		fungibleTokenAddress,
		metayaUtilityCoinAddress,
	)

	commonweal2Address, commonweal2Signer, _ := test.CreateAccount(t, b)
	metayaUtilityCoin.SetupAccount(
		t, b,
		commonweal2Address,
		commonweal2Signer,
		fungibleTokenAddress,
		metayaUtilityCoinAddress,
	)

	tempMetayaAddress, tempMetayaSigner, _ := test.CreateAccount(t, b)
	metayaUtilityCoin.SetupAccount(
		t, b,
		tempMetayaAddress,
		tempMetayaSigner,
		fungibleTokenAddress,
		metayaUtilityCoinAddress,
	)

	t.Run("Should be able to add a copyright owner", func(t *testing.T) {
		metayaBeneficiaryCut.CreateCopyrightOwner(
			t, b,
			"Copyright Owner 001",
			copyrightOwner1Address,
			fungibleTokenAddress, metayaUtilityCoinAddress, metayaBeneficiaryCutAddress,
			metayaBeneficiaryCutSigner,
		)

		// Assert that vault balances are correct

		copyrightOwnerNames := test.ExecuteScriptAndCheck(
			t, b,
			metayaBeneficiaryCut.GetCopyrightOwnerNamesScript(metayaBeneficiaryCutAddress.String()),
			nil,
		)
		assert.EqualValues(t, cadence.NewArray([]cadence.Value{cadence.String("Copyright Owner 001")}), copyrightOwnerNames)

		copyrightOwnerAmount := test.ExecuteScriptAndCheck(
			t, b,
			metayaBeneficiaryCut.GetCopyrightOwnerAmountScript(metayaBeneficiaryCutAddress.String()),
			nil,
		)
		assert.EqualValues(t, cadence.NewInt(1), copyrightOwnerAmount)

		copyrightOwnerContain := test.ExecuteScriptAndCheck(
			t, b,
			metayaBeneficiaryCut.GetCopyrightOwnerContainScript(metayaBeneficiaryCutAddress.String()),
			[][]byte{jsoncdc.MustEncode(cadence.String("Copyright Owner 001"))},
		)
		assert.EqualValues(t, cadence.NewBool(true), copyrightOwnerContain)

		copyrightOwnerAddress := test.ExecuteScriptAndCheck(
			t, b,
			metayaBeneficiaryCut.GetCopyrightOwnerAddressByNameScript(metayaBeneficiaryCutAddress.String()),
			[][]byte{jsoncdc.MustEncode(cadence.String("Copyright Owner 001"))},
		)
		assert.EqualValues(t, cadence.Address(copyrightOwner1Address), copyrightOwnerAddress)
	})

	t.Run("Should be able to del a copyright owner", func(t *testing.T) {
		metayaBeneficiaryCut.CreateCopyrightOwner(
			t, b,
			"Copyright Owner 001",
			copyrightOwner1Address,
			fungibleTokenAddress, metayaUtilityCoinAddress, metayaBeneficiaryCutAddress,
			metayaBeneficiaryCutSigner,
		)

		tx := flow.NewTransaction().
			SetScript(metayaBeneficiaryCut.DelCopyrightOwnerTransaction(metayaBeneficiaryCutAddress.String())).
			SetGasLimit(100).
			SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
			SetPayer(b.ServiceKey().Address).
			AddAuthorizer(metayaBeneficiaryCutAddress)

		_ = tx.AddArgument(cadence.String("Copyright Owner 001"))

		test.SignAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, metayaBeneficiaryCutAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), metayaBeneficiaryCutSigner},
			false,
		)

		// Assert that vault balances are correct

		copyrightOwnerNames := test.ExecuteScriptAndCheck(
			t, b,
			metayaBeneficiaryCut.GetCopyrightOwnerNamesScript(metayaBeneficiaryCutAddress.String()),
			nil,
		)
		assert.EqualValues(t, cadence.NewArray([]cadence.Value{}), copyrightOwnerNames)

		copyrightOwnerAmount := test.ExecuteScriptAndCheck(
			t, b,
			metayaBeneficiaryCut.GetCopyrightOwnerAmountScript(metayaBeneficiaryCutAddress.String()),
			nil,
		)
		assert.EqualValues(t, cadence.NewInt(0), copyrightOwnerAmount)

		copyrightOwnerContain := test.ExecuteScriptAndCheck(
			t, b,
			metayaBeneficiaryCut.GetCopyrightOwnerContainScript(metayaBeneficiaryCutAddress.String()),
			[][]byte{jsoncdc.MustEncode(cadence.String("Copyright Owner 001"))},
		)
		assert.EqualValues(t, cadence.NewBool(false), copyrightOwnerContain)
	})

	t.Run("Should be able to add a commonweal", func(t *testing.T) {
		metayaBeneficiaryCut.CreateCommonweal(
			t, b,
			"Commonweal 001",
			"0.002",
			commonweal1Address,
			fungibleTokenAddress, metayaUtilityCoinAddress, metayaBeneficiaryCutAddress,
			metayaBeneficiaryCutSigner,
		)

		// Assert that vault balances are correct

		commonwealNames := test.ExecuteScriptAndCheck(
			t, b,
			metayaBeneficiaryCut.GetCommonwealNamesScript(metayaBeneficiaryCutAddress.String()),
			nil,
		)
		assert.EqualValues(t, cadence.NewArray([]cadence.Value{cadence.String("Commonweal 001")}), commonwealNames)

		commonwealCutPercentageByName := test.ExecuteScriptAndCheck(
			t, b,
			metayaBeneficiaryCut.GetCommonwealCutPercentageByNameScript(metayaBeneficiaryCutAddress.String()),
			[][]byte{jsoncdc.MustEncode(cadence.String("Commonweal 001"))},
		)
		assert.EqualValues(t, cadence.NewOptional(test.CadenceUFix64("0.002")), commonwealCutPercentageByName)
	})

	t.Run("Should be able to del a commonweal", func(t *testing.T) {
		metayaBeneficiaryCut.CreateCommonweal(
			t, b,
			"Commonweal 001",
			"0.002",
			commonweal1Address,
			fungibleTokenAddress, metayaUtilityCoinAddress, metayaBeneficiaryCutAddress,
			metayaBeneficiaryCutSigner,
		)

		tx := flow.NewTransaction().
			SetScript(metayaBeneficiaryCut.DelCommonwealTransaction(metayaBeneficiaryCutAddress.String())).
			SetGasLimit(100).
			SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
			SetPayer(b.ServiceKey().Address).
			AddAuthorizer(metayaBeneficiaryCutAddress)

		_ = tx.AddArgument(cadence.String("Commonweal 001"))

		test.SignAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, metayaBeneficiaryCutAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), metayaBeneficiaryCutSigner},
			false,
		)

		// Assert that vault balances are correct

		commonwealNames := test.ExecuteScriptAndCheck(
			t, b,
			metayaBeneficiaryCut.GetCommonwealNamesScript(metayaBeneficiaryCutAddress.String()),
			nil,
		)
		assert.EqualValues(t, cadence.NewArray([]cadence.Value{}), commonwealNames)

		commonwealCutPercentageByName := test.ExecuteScriptAndCheck(
			t, b,
			metayaBeneficiaryCut.GetCommonwealCutPercentageByNameScript(metayaBeneficiaryCutAddress.String()),
			[][]byte{jsoncdc.MustEncode(cadence.String("Commonweal 001"))},
		)
		assert.EqualValues(t, cadence.NewOptional(nil), commonwealCutPercentageByName)
	})

	t.Run("Should be able to set metaya capability", func(t *testing.T) {
		tx := flow.NewTransaction().
			SetScript(metayaBeneficiaryCut.SetMetayaCapabilityTransaction(fungibleTokenAddress.String(), metayaUtilityCoinAddress.String(), metayaBeneficiaryCutAddress.String())).
			SetGasLimit(100).
			SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
			SetPayer(b.ServiceKey().Address).
			AddAuthorizer(metayaBeneficiaryCutAddress)

		_ = tx.AddArgument(cadence.Address(tempMetayaAddress))

		test.SignAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, metayaBeneficiaryCutAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), metayaBeneficiaryCutSigner},
			false,
		)

		// Assert that vault balances are correct

		metayaAddress := test.ExecuteScriptAndCheck(
			t, b,
			metayaBeneficiaryCut.GetMetayaAddressScript(metayaBeneficiaryCutAddress.String()),
			nil,
		)
		assert.EqualValues(t, cadence.Address(tempMetayaAddress), metayaAddress)
	})

	t.Run("Should be able to set metaya marketCutPercentage", func(t *testing.T) {
		tx := flow.NewTransaction().
			SetScript(metayaBeneficiaryCut.SetMetayaMarketCutPercentageTransaction(metayaBeneficiaryCutAddress.String())).
			SetGasLimit(100).
			SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
			SetPayer(b.ServiceKey().Address).
			AddAuthorizer(metayaBeneficiaryCutAddress)

		_ = tx.AddArgument(test.CadenceUFix64("0.05"))

		test.SignAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, metayaBeneficiaryCutAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), metayaBeneficiaryCutSigner},
			false,
		)

		// Assert that vault balances are correct

		metayaAddress := test.ExecuteScriptAndCheck(
			t, b,
			metayaBeneficiaryCut.GetMetayaMarketCutPercentageScript(metayaBeneficiaryCutAddress.String()),
			nil,
		)
		assert.EqualValues(t, test.CadenceUFix64("0.05"), metayaAddress)
	})

	t.Run("Should be able to add sell cutPercentage in store", func(t *testing.T) {
		metayaBeneficiaryCut.CreateCopyrightOwner(
			t, b,
			"Copyright Owner 001",
			copyrightOwner1Address,
			fungibleTokenAddress, metayaUtilityCoinAddress, metayaBeneficiaryCutAddress,
			metayaBeneficiaryCutSigner,
		)

		metayaBeneficiaryCut.CreateCopyrightOwner(
			t, b,
			"Copyright Owner 002",
			copyrightOwner2Address,
			fungibleTokenAddress, metayaUtilityCoinAddress, metayaBeneficiaryCutAddress,
			metayaBeneficiaryCutSigner,
		)
		
		metayaBeneficiaryCut.CreateCopyrightOwner(
			t, b,
			"Metaya",
			copyrightOwner2Address,
			fungibleTokenAddress, metayaUtilityCoinAddress, metayaBeneficiaryCutAddress,
			metayaBeneficiaryCutSigner,
		)

		pair1 := cadence.KeyValuePair{Key: cadence.String("Copyright Owner 001"), Value: test.CadenceUFix64("0.5")}
		pair2 := cadence.KeyValuePair{Key: cadence.String("Copyright Owner 002"), Value: test.CadenceUFix64("0.3")}
		pair3 := cadence.KeyValuePair{Key: cadence.String("Metaya"), Value: test.CadenceUFix64("0.2")}
		recipientPairs := make([]cadence.KeyValuePair, 3)
		recipientPairs[0] = pair1
		recipientPairs[1] = pair2
		recipientPairs[2] = pair3

		tx := flow.NewTransaction().
			SetScript(metayaBeneficiaryCut.SetStoreCutPercentageTransaction(metayaBeneficiaryCutAddress.String())).
			SetGasLimit(100).
			SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
			SetPayer(b.ServiceKey().Address).
			AddAuthorizer(metayaBeneficiaryCutAddress)

		_ = tx.AddArgument(cadence.NewUInt32(1))
		_ = tx.AddArgument(cadence.NewDictionary(recipientPairs))

		test.SignAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, metayaBeneficiaryCutAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), metayaBeneficiaryCutSigner},
			false,
		)

		// Assert that vault balances are correct

		storeCutPercentagesAmount := test.ExecuteScriptAndCheck(
			t, b,
			metayaBeneficiaryCut.GetStoreCutPercentagesAmountScript(metayaBeneficiaryCutAddress.String()),
			nil,
		)
		assert.EqualValues(t, cadence.NewInt(1), storeCutPercentagesAmount)

		storeCutPercentageByName1 := test.ExecuteScriptAndCheck(
			t, b,
			metayaBeneficiaryCut.GetStoreCutPercentageByNameScript(metayaBeneficiaryCutAddress.String()),
			[][]byte{jsoncdc.MustEncode(cadence.UInt32(1)),jsoncdc.MustEncode(cadence.String("Copyright Owner 001"))},
		)
		assert.EqualValues(t, cadence.NewOptional(test.CadenceUFix64("0.5")), storeCutPercentageByName1)

		storeCutPercentageByName2 := test.ExecuteScriptAndCheck(
			t, b,
			metayaBeneficiaryCut.GetStoreCutPercentageByNameScript(metayaBeneficiaryCutAddress.String()),
			[][]byte{jsoncdc.MustEncode(cadence.UInt32(1)),jsoncdc.MustEncode(cadence.String("Copyright Owner 002"))},
		)
		assert.EqualValues(t, cadence.NewOptional(test.CadenceUFix64("0.3")), storeCutPercentageByName2)

		storeCutPercentageByName3 := test.ExecuteScriptAndCheck(
			t, b,
			metayaBeneficiaryCut.GetStoreCutPercentageByNameScript(metayaBeneficiaryCutAddress.String()),
			[][]byte{jsoncdc.MustEncode(cadence.UInt32(1)),jsoncdc.MustEncode(cadence.String("Metaya"))},
		)
		assert.EqualValues(t, cadence.NewOptional(test.CadenceUFix64("0.2")), storeCutPercentageByName3)

		t.Run("Should be able to del sell cutPercentage in store", func(t *testing.T) {
			tx := flow.NewTransaction().
			SetScript(metayaBeneficiaryCut.DelStoreCutPercentageTransaction(metayaBeneficiaryCutAddress.String())).
			SetGasLimit(100).
			SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
			SetPayer(b.ServiceKey().Address).
			AddAuthorizer(metayaBeneficiaryCutAddress)

			_ = tx.AddArgument(cadence.NewUInt32(1))

			test.SignAndSubmit(
				t, b, tx,
				[]flow.Address{b.ServiceKey().Address, metayaBeneficiaryCutAddress},
				[]crypto.Signer{b.ServiceKey().Signer(), metayaBeneficiaryCutSigner},
				false,
			)

			// Assert that vault balances are correct

			storeCutPercentagesAmount := test.ExecuteScriptAndCheck(
				t, b,
				metayaBeneficiaryCut.GetStoreCutPercentagesAmountScript(metayaBeneficiaryCutAddress.String()),
				nil,
			)
			assert.EqualValues(t, cadence.NewInt(0), storeCutPercentagesAmount)

			storeCutPercentageByName1 := test.ExecuteScriptAndCheck(
				t, b,
				metayaBeneficiaryCut.GetStoreCutPercentageByNameScript(metayaBeneficiaryCutAddress.String()),
				[][]byte{jsoncdc.MustEncode(cadence.UInt32(1)),jsoncdc.MustEncode(cadence.String("Copyright Owner 001"))},
			)
			assert.EqualValues(t, cadence.NewOptional(nil), storeCutPercentageByName1)

			storeCutPercentageByName2 := test.ExecuteScriptAndCheck(
				t, b,
				metayaBeneficiaryCut.GetStoreCutPercentageByNameScript(metayaBeneficiaryCutAddress.String()),
				[][]byte{jsoncdc.MustEncode(cadence.UInt32(1)),jsoncdc.MustEncode(cadence.String("Copyright Owner 002"))},
			)
			assert.EqualValues(t, cadence.NewOptional(nil), storeCutPercentageByName2)

			storeCutPercentageByName3 := test.ExecuteScriptAndCheck(
				t, b,
				metayaBeneficiaryCut.GetStoreCutPercentageByNameScript(metayaBeneficiaryCutAddress.String()),
				[][]byte{jsoncdc.MustEncode(cadence.UInt32(1)),jsoncdc.MustEncode(cadence.String("Metaya"))},
			)
			assert.EqualValues(t, cadence.NewOptional(nil), storeCutPercentageByName3)
		})
	})

	t.Run("Should be able to add sell cutPercentage in pack", func(t *testing.T) {
		metayaBeneficiaryCut.CreateCopyrightOwner(
			t, b,
			"Copyright Owner 001",
			copyrightOwner1Address,
			fungibleTokenAddress, metayaUtilityCoinAddress, metayaBeneficiaryCutAddress,
			metayaBeneficiaryCutSigner,
		)

		metayaBeneficiaryCut.CreateCopyrightOwner(
			t, b,
			"Copyright Owner 002",
			copyrightOwner2Address,
			fungibleTokenAddress, metayaUtilityCoinAddress, metayaBeneficiaryCutAddress,
			metayaBeneficiaryCutSigner,
		)
		
		metayaBeneficiaryCut.CreateCopyrightOwner(
			t, b,
			"Metaya",
			copyrightOwner2Address,
			fungibleTokenAddress, metayaUtilityCoinAddress, metayaBeneficiaryCutAddress,
			metayaBeneficiaryCutSigner,
		)

		pair1 := cadence.KeyValuePair{Key: cadence.String("Copyright Owner 001"), Value: test.CadenceUFix64("0.2")}
		pair2 := cadence.KeyValuePair{Key: cadence.String("Copyright Owner 002"), Value: test.CadenceUFix64("0.7")}
		pair3 := cadence.KeyValuePair{Key: cadence.String("Metaya"), Value: test.CadenceUFix64("0.1")}
		recipientPairs := make([]cadence.KeyValuePair, 3)
		recipientPairs[0] = pair1
		recipientPairs[1] = pair2
		recipientPairs[2] = pair3

		tx := flow.NewTransaction().
			SetScript(metayaBeneficiaryCut.SetPackCutPercentageTransaction(metayaBeneficiaryCutAddress.String())).
			SetGasLimit(100).
			SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
			SetPayer(b.ServiceKey().Address).
			AddAuthorizer(metayaBeneficiaryCutAddress)

		_ = tx.AddArgument(cadence.NewUInt32(1))
		_ = tx.AddArgument(cadence.NewDictionary(recipientPairs))

		test.SignAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, metayaBeneficiaryCutAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), metayaBeneficiaryCutSigner},
			false,
		)

		// Assert that vault balances are correct

		packCutPercentagesAmount := test.ExecuteScriptAndCheck(
			t, b,
			metayaBeneficiaryCut.GetPackCutPercentagesAmountScript(metayaBeneficiaryCutAddress.String()),
			nil,
		)
		assert.EqualValues(t, cadence.NewInt(1), packCutPercentagesAmount)

		packCutPercentageByName1 := test.ExecuteScriptAndCheck(
			t, b,
			metayaBeneficiaryCut.GetPackCutPercentageByNameScript(metayaBeneficiaryCutAddress.String()),
			[][]byte{jsoncdc.MustEncode(cadence.UInt32(1)),jsoncdc.MustEncode(cadence.String("Copyright Owner 001"))},
		)
		assert.EqualValues(t, cadence.NewOptional(test.CadenceUFix64("0.2")), packCutPercentageByName1)

		packCutPercentageByName2 := test.ExecuteScriptAndCheck(
			t, b,
			metayaBeneficiaryCut.GetPackCutPercentageByNameScript(metayaBeneficiaryCutAddress.String()),
			[][]byte{jsoncdc.MustEncode(cadence.UInt32(1)),jsoncdc.MustEncode(cadence.String("Copyright Owner 002"))},
		)
		assert.EqualValues(t, cadence.NewOptional(test.CadenceUFix64("0.7")), packCutPercentageByName2)

		packCutPercentageByName3 := test.ExecuteScriptAndCheck(
			t, b,
			metayaBeneficiaryCut.GetPackCutPercentageByNameScript(metayaBeneficiaryCutAddress.String()),
			[][]byte{jsoncdc.MustEncode(cadence.UInt32(1)),jsoncdc.MustEncode(cadence.String("Metaya"))},
		)
		assert.EqualValues(t, cadence.NewOptional(test.CadenceUFix64("0.1")), packCutPercentageByName3)

		t.Run("Should be able to del sell cutPercentage in pack", func(t *testing.T) {
			tx := flow.NewTransaction().
			SetScript(metayaBeneficiaryCut.DelPackCutPercentageTransaction(metayaBeneficiaryCutAddress.String())).
			SetGasLimit(100).
			SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
			SetPayer(b.ServiceKey().Address).
			AddAuthorizer(metayaBeneficiaryCutAddress)

			_ = tx.AddArgument(cadence.NewUInt32(1))

			test.SignAndSubmit(
				t, b, tx,
				[]flow.Address{b.ServiceKey().Address, metayaBeneficiaryCutAddress},
				[]crypto.Signer{b.ServiceKey().Signer(), metayaBeneficiaryCutSigner},
				false,
			)

			// Assert that vault balances are correct

			packCutPercentagesAmount := test.ExecuteScriptAndCheck(
				t, b,
				metayaBeneficiaryCut.GetPackCutPercentagesAmountScript(metayaBeneficiaryCutAddress.String()),
				nil,
			)
			assert.EqualValues(t, cadence.NewInt(0), packCutPercentagesAmount)

			packCutPercentageByName1 := test.ExecuteScriptAndCheck(
				t, b,
				metayaBeneficiaryCut.GetPackCutPercentageByNameScript(metayaBeneficiaryCutAddress.String()),
				[][]byte{jsoncdc.MustEncode(cadence.UInt32(1)),jsoncdc.MustEncode(cadence.String("Copyright Owner 001"))},
			)
			assert.EqualValues(t, cadence.NewOptional(nil), packCutPercentageByName1)

			packCutPercentageByName2 := test.ExecuteScriptAndCheck(
				t, b,
				metayaBeneficiaryCut.GetPackCutPercentageByNameScript(metayaBeneficiaryCutAddress.String()),
				[][]byte{jsoncdc.MustEncode(cadence.UInt32(1)),jsoncdc.MustEncode(cadence.String("Copyright Owner 002"))},
			)
			assert.EqualValues(t, cadence.NewOptional(nil), packCutPercentageByName2)

			packCutPercentageByName3 := test.ExecuteScriptAndCheck(
				t, b,
				metayaBeneficiaryCut.GetPackCutPercentageByNameScript(metayaBeneficiaryCutAddress.String()),
				[][]byte{jsoncdc.MustEncode(cadence.UInt32(1)),jsoncdc.MustEncode(cadence.String("Metaya"))},
			)
			assert.EqualValues(t, cadence.NewOptional(nil), packCutPercentageByName3)
		})
	})

	t.Run("Should be able to add sell cutPercentage in market", func(t *testing.T) {
		metayaBeneficiaryCut.CreateCopyrightOwner(
			t, b,
			"Copyright Owner 001",
			copyrightOwner1Address,
			fungibleTokenAddress, metayaUtilityCoinAddress, metayaBeneficiaryCutAddress,
			metayaBeneficiaryCutSigner,
		)

		metayaBeneficiaryCut.CreateCopyrightOwner(
			t, b,
			"Copyright Owner 002",
			copyrightOwner2Address,
			fungibleTokenAddress, metayaUtilityCoinAddress, metayaBeneficiaryCutAddress,
			metayaBeneficiaryCutSigner,
		)
		
		pair1 := cadence.KeyValuePair{Key: cadence.String("Copyright Owner 001"), Value: test.CadenceUFix64("0.01")}
		pair2 := cadence.KeyValuePair{Key: cadence.String("Copyright Owner 002"), Value: test.CadenceUFix64("0.04")}
		recipientPairs := make([]cadence.KeyValuePair, 2)
		recipientPairs[0] = pair1
		recipientPairs[1] = pair2

		tx := flow.NewTransaction().
			SetScript(metayaBeneficiaryCut.SetMarketCutPercentageTransaction(metayaBeneficiaryCutAddress.String())).
			SetGasLimit(100).
			SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
			SetPayer(b.ServiceKey().Address).
			AddAuthorizer(metayaBeneficiaryCutAddress)

		_ = tx.AddArgument(cadence.NewUInt32(1))
		_ = tx.AddArgument(cadence.NewDictionary(recipientPairs))

		test.SignAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, metayaBeneficiaryCutAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), metayaBeneficiaryCutSigner},
			false,
		)

		// Assert that vault balances are correct

		marketCutPercentagesAmount := test.ExecuteScriptAndCheck(
			t, b,
			metayaBeneficiaryCut.GetMarketCutPercentagesAmountScript(metayaBeneficiaryCutAddress.String()),
			nil,
		)
		assert.EqualValues(t, cadence.NewInt(1), marketCutPercentagesAmount)

		marketCutPercentageByName1 := test.ExecuteScriptAndCheck(
			t, b,
			metayaBeneficiaryCut.GetMarketCutPercentageByNameScript(metayaBeneficiaryCutAddress.String()),
			[][]byte{jsoncdc.MustEncode(cadence.UInt32(1)),jsoncdc.MustEncode(cadence.String("Copyright Owner 001"))},
		)
		assert.EqualValues(t, cadence.NewOptional(test.CadenceUFix64("0.01")), marketCutPercentageByName1)

		marketCutPercentageByName2 := test.ExecuteScriptAndCheck(
			t, b,
			metayaBeneficiaryCut.GetMarketCutPercentageByNameScript(metayaBeneficiaryCutAddress.String()),
			[][]byte{jsoncdc.MustEncode(cadence.UInt32(1)),jsoncdc.MustEncode(cadence.String("Copyright Owner 002"))},
		)
		assert.EqualValues(t, cadence.NewOptional(test.CadenceUFix64("0.04")), marketCutPercentageByName2)

		t.Run("Should be able to del sell cutPercentage in market", func(t *testing.T) {
			tx := flow.NewTransaction().
			SetScript(metayaBeneficiaryCut.DelMarketCutPercentageTransaction(metayaBeneficiaryCutAddress.String())).
			SetGasLimit(100).
			SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
			SetPayer(b.ServiceKey().Address).
			AddAuthorizer(metayaBeneficiaryCutAddress)

			_ = tx.AddArgument(cadence.NewUInt32(1))

			test.SignAndSubmit(
				t, b, tx,
				[]flow.Address{b.ServiceKey().Address, metayaBeneficiaryCutAddress},
				[]crypto.Signer{b.ServiceKey().Signer(), metayaBeneficiaryCutSigner},
				false,
			)

			// Assert that vault balances are correct

			marketCutPercentagesAmount := test.ExecuteScriptAndCheck(
				t, b,
				metayaBeneficiaryCut.GetMarketCutPercentagesAmountScript(metayaBeneficiaryCutAddress.String()),
				nil,
			)
			assert.EqualValues(t, cadence.NewInt(0), marketCutPercentagesAmount)

			marketCutPercentageByName1 := test.ExecuteScriptAndCheck(
				t, b,
				metayaBeneficiaryCut.GetMarketCutPercentageByNameScript(metayaBeneficiaryCutAddress.String()),
				[][]byte{jsoncdc.MustEncode(cadence.UInt32(1)),jsoncdc.MustEncode(cadence.String("Copyright Owner 001"))},
			)
			assert.EqualValues(t, cadence.NewOptional(nil), marketCutPercentageByName1)

			marketCutPercentageByName2 := test.ExecuteScriptAndCheck(
				t, b,
				metayaBeneficiaryCut.GetMarketCutPercentageByNameScript(metayaBeneficiaryCutAddress.String()),
				[][]byte{jsoncdc.MustEncode(cadence.UInt32(1)),jsoncdc.MustEncode(cadence.String("Copyright Owner 002"))},
			)
			assert.EqualValues(t, cadence.NewOptional(nil), marketCutPercentageByName2)
		})
	})
}
