package test

import (
	"testing"

	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/stretchr/testify/assert"

	"github.com/0xRaito/metaya-smart-contracts/tests/go/test"
	"github.com/0xRaito/metaya-smart-contracts/tests/go/metayaUtilityCoin"
	"github.com/0xRaito/metaya-smart-contracts/tests/go/soundlinksDID"
	"github.com/0xRaito/metaya-smart-contracts/tests/go/metayaNFT"
	"github.com/0xRaito/metaya-smart-contracts/tests/go/metayaBeneficiaryCut"
	"github.com/0xRaito/metaya-smart-contracts/tests/go/nftStorefront"
	
)

func TestNFTStorefrontDeployContracts(t *testing.T) {
	b := test.NewBlockchain()
	nftStorefront.DeployContracts(t, b)
}

func TestNFTStorefrontSetupAccount(t *testing.T) {
	b := test.NewBlockchain()

	contracts := nftStorefront.DeployContracts(t, b)

	t.Run("Should be able to create an empty Storefront", func(t *testing.T) {
		userAddress, userSigner, _ := test.CreateAccount(t, b)
		nftStorefront.SetupAccount(t, b, userAddress, userSigner, contracts.NFTStorefrontAddress)
	})
}

func TestNFTStorefrontCreateListing(t *testing.T) {
	b := test.NewBlockchain()

	contracts := nftStorefront.DeployContracts(t, b)

	sellerAddress, sellerSigner := nftStorefront.CreateAccount(t, b, contracts)
	buyerAddress, buyerSigner := nftStorefront.CreatePurchaserAccount(t, b, contracts)

	copyrightOwner1Address, _ := nftStorefront.CreateAccount(t, b, contracts)
	metayaBeneficiaryCut.CreateCopyrightOwner(
		t, b,
		"Copyright Owner 001",
		copyrightOwner1Address,
		contracts.FungibleTokenAddress, contracts.MetayaUtilityCoinAddress, contracts.MetayaBeneficiaryCutAddress,
		contracts.MetayaBeneficiaryCutSigner,
	)

	copyrightOwner2Address, _ := nftStorefront.CreateAccount(t, b, contracts)
	metayaBeneficiaryCut.CreateCopyrightOwner(
		t, b,
		"Copyright Owner 002",
		copyrightOwner2Address,
		contracts.FungibleTokenAddress, contracts.MetayaUtilityCoinAddress, contracts.MetayaBeneficiaryCutAddress,
		contracts.MetayaBeneficiaryCutSigner,
	)

	tempMetayaAddress, _ := nftStorefront.CreateAccount(t, b, contracts)
	tx := flow.NewTransaction().
		SetScript(metayaBeneficiaryCut.SetMetayaCapabilityTransaction(contracts.FungibleTokenAddress.String(), contracts.MetayaUtilityCoinAddress.String(), contracts.MetayaBeneficiaryCutAddress.String())).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(contracts.MetayaBeneficiaryCutAddress)

	_ = tx.AddArgument(cadence.Address(tempMetayaAddress))

	test.SignAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, contracts.MetayaBeneficiaryCutAddress},
		[]crypto.Signer{b.ServiceKey().Signer(), contracts.MetayaBeneficiaryCutSigner},
		false,
	)

	// Set market cut percentage
	pair1 := cadence.KeyValuePair{Key: cadence.String("Copyright Owner 001"), Value: test.CadenceUFix64("0.01")}
	pair2 := cadence.KeyValuePair{Key: cadence.String("Copyright Owner 002"), Value: test.CadenceUFix64("0.04")}
	recipientPairs := make([]cadence.KeyValuePair, 2)
	recipientPairs[0] = pair1
	recipientPairs[1] = pair2

	tx = flow.NewTransaction().
		SetScript(metayaBeneficiaryCut.SetMarketCutPercentageTransaction(contracts.MetayaBeneficiaryCutAddress.String())).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(contracts.MetayaBeneficiaryCutAddress)

	_ = tx.AddArgument(cadence.NewUInt32(1))
	_ = tx.AddArgument(cadence.NewDictionary(recipientPairs))

	test.SignAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, contracts.MetayaBeneficiaryCutAddress},
		[]crypto.Signer{b.ServiceKey().Signer(), contracts.MetayaBeneficiaryCutSigner},
		false,
	)

	// Create play
	title := cadence.String("Title")
	artwork1 := cadence.String("Artwork 001")

	tx = flow.NewTransaction().
		SetScript(metaya.CreatePlayTransaction(contracts)).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(contracts.MetayaAddress)

	metadata := []cadence.KeyValuePair{{Key: title, Value: artwork1}}
	play := cadence.NewDictionary(metadata)
	_ = tx.AddArgument(play)

	test.SignAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, contracts.MetayaAddress},
		[]crypto.Signer{b.ServiceKey().Signer(), contracts.MetayaSigner},
		false,
	)

	// Create set
	tx = flow.NewTransaction().
		SetScript(metaya.CreateSetTransaction(contracts)).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(contracts.MetayaAddress)

	_ = tx.AddArgument(cadence.String("Set 001"))

	test.SignAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, contracts.MetayaAddress},
		[]crypto.Signer{b.ServiceKey().Signer(), contracts.MetayaSigner},
		false,
	)

	// Add plays to Set
	tx = flow.NewTransaction().
		SetScript(metaya.AddPlayToSetTransaction(contracts)).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(contracts.MetayaAddress)

	_ = tx.AddArgument(cadence.NewUInt32(1))
	_ = tx.AddArgument(cadence.NewUInt32(1))

	test.SignAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, contracts.MetayaAddress},
		[]crypto.Signer{b.ServiceKey().Signer(), contracts.MetayaSigner},
		false,
	)

	// Mint moments to sellerAddress
	hashs := []string{"FFFFFFFFFFFFFFFFFFFA","FFFFFFFFFFFFFFFFFFFB"}

	soundlinksDID.SetupAccount(
		t, b,
		contracts.MetayaAddress,
		contracts.MetayaSigner,
		contracts.NonFungibleTokenAddress,
		contracts.SoundlinksDIDAddress,
	)

	soundlinksDID.MintDIDs(
		t, b,
		contracts.NonFungibleTokenAddress, contracts.SoundlinksDIDAddress, contracts.MetayaAddress,
		contracts.SoundlinksDIDSigner,
		hashs,
	)

	metaya.MintNFTs(
		t, b,
		contracts,
		uint32(1), uint32(1),
		uint32(2),
		sellerAddress,
		false,
	)

	t.Run("Should be able to create a sale listing and list it", func(t *testing.T) {
		tokenToList := uint64(1)
		tokenPrice := "100.0"
		
		// The seller account lists the item
		listingResourceID := nftStorefront.ListItem(
			t, b,
			contracts,
			sellerAddress,
			sellerSigner,
			tokenToList,
			tokenPrice,
			false,
		)

		ids := test.ExecuteScriptAndCheck(
			t, b,
			nftStorefront.NFTStorefrontGetListingIdsScript(contracts),
			[][]byte{jsoncdc.MustEncode(cadence.Address(sellerAddress))},
		)
		idsArray := cadence.NewArray([]cadence.Value{cadence.NewUInt64(listingResourceID)})
		assert.EqualValues(t, idsArray, ids)

		test.ExecuteScriptAndCheck(
			t, b,
			nftStorefront.NFTStorefrontGetListingDetailsScript(contracts),
			[][]byte{jsoncdc.MustEncode(cadence.Address(sellerAddress)),jsoncdc.MustEncode(cadence.NewUInt64(listingResourceID))},
		)
		
		assert.EqualValues(t, idsArray, ids)

		t.Run("Should be able to remove a sale listing", func(t *testing.T) {
			nftStorefront.RemoveItem(
				t, b,
				contracts,
				sellerAddress,
				sellerSigner,
				listingResourceID,
				false,
			)
		})
	})

	t.Run("Should be able to purchase a sale listing by IUC", func(t *testing.T) {
		tokenToList := uint64(1)
		tokenPrice := "100.0"
	
		// other seller account lists the item
		listingResourceID := nftStorefront.ListItem(
			t, b,
			contracts,
			sellerAddress,
			sellerSigner,
			tokenToList,
			tokenPrice,
			false,
		)

		// Fund the purchase
		metayaUtilityCoin.MintTokens(
			t, b,
			contracts.FungibleTokenAddress,
			contracts.MetayaUtilityCoinAddress,
			contracts.MetayaUtilityCoinSigner,
			buyerAddress,
			"100.0",
			false,
		)

		// Make the purchase by IUC
		nftStorefront.PurchaseItemByIUC(
			t, b,
			contracts,
			buyerAddress,
			buyerSigner,
			sellerAddress,
			listingResourceID,
			false,
		)

		result := test.ExecuteScriptAndCheck(
			t, b,
			metaya.GetCollectionIdsScript(contracts),
			[][]byte{jsoncdc.MustEncode(cadence.Address(sellerAddress))},
		)
		idsArray := cadence.NewArray([]cadence.Value{cadence.NewUInt64(2)})
		assert.EqualValues(t, idsArray, result)

		result = test.ExecuteScriptAndCheck(
			t, b,
			metaya.GetCollectionIdsScript(contracts),
			[][]byte{jsoncdc.MustEncode(cadence.Address(buyerAddress))},
		)
		idsArray = cadence.NewArray([]cadence.Value{cadence.NewUInt64(1)})
		assert.EqualValues(t, idsArray, result)

		buyerBalance := test.ExecuteScriptAndCheck(
			t, b,
			metayaUtilityCoin.GetBalanceScript(contracts.FungibleTokenAddress.String(), contracts.MetayaUtilityCoinAddress.String()),
			[][]byte{jsoncdc.MustEncode(cadence.Address(buyerAddress))},
		)
		assert.EqualValues(t, test.CadenceUFix64("0.0"), buyerBalance)

		copyrightOwner1Balance := test.ExecuteScriptAndCheck(
			t, b,
			metayaUtilityCoin.GetBalanceScript(contracts.FungibleTokenAddress.String(), contracts.MetayaUtilityCoinAddress.String()),
			[][]byte{jsoncdc.MustEncode(cadence.Address(copyrightOwner1Address))},
		)
		assert.EqualValues(t, test.CadenceUFix64("1.0"), copyrightOwner1Balance)

		copyrightOwner2Balance := test.ExecuteScriptAndCheck(
			t, b,
			metayaUtilityCoin.GetBalanceScript(contracts.FungibleTokenAddress.String(), contracts.MetayaUtilityCoinAddress.String()),
			[][]byte{jsoncdc.MustEncode(cadence.Address(copyrightOwner2Address))},
		)
		assert.EqualValues(t, test.CadenceUFix64("4.0"), copyrightOwner2Balance)

		tempMetayaBalance := test.ExecuteScriptAndCheck(
			t, b,
			metayaUtilityCoin.GetBalanceScript(contracts.FungibleTokenAddress.String(), contracts.MetayaUtilityCoinAddress.String()),
			[][]byte{jsoncdc.MustEncode(cadence.Address(tempMetayaAddress))},
		)
		assert.EqualValues(t, test.CadenceUFix64("3.0"), tempMetayaBalance)

		sellerBalance := test.ExecuteScriptAndCheck(
			t, b,
			metayaUtilityCoin.GetBalanceScript(contracts.FungibleTokenAddress.String(), contracts.MetayaUtilityCoinAddress.String()),
			[][]byte{jsoncdc.MustEncode(cadence.Address(sellerAddress))},
		)
		assert.EqualValues(t, test.CadenceUFix64("92.0"), sellerBalance)
	})

	t.Run("Should be able to purchase a sale listing by cash", func(t *testing.T) {
		tokenToList := uint64(2)
		tokenPrice := "100.0"
	
		// other seller account lists the item
		listingResourceID := nftStorefront.ListItem(
			t, b,
			contracts,
			sellerAddress,
			sellerSigner,
			tokenToList,
			tokenPrice,
			false,
		)

		// Make the purchase by cash
		nftStorefront.PurchaseItemByCash(
			t, b,
			contracts,
			sellerAddress,
			buyerAddress,
			listingResourceID,
			false,
		)

		result := test.ExecuteScriptAndCheck(
			t, b,
			metaya.GetCollectionIdsScript(contracts),
			[][]byte{jsoncdc.MustEncode(cadence.Address(sellerAddress))},
		)
		idsArray := cadence.NewArray([]cadence.Value{})
		assert.EqualValues(t, idsArray, result)

		result = test.ExecuteScriptAndCheck(
			t, b,
			metaya.GetCollectionIdsScript(contracts),
			[][]byte{jsoncdc.MustEncode(cadence.Address(buyerAddress))},
		)
		idsArray = cadence.NewArray([]cadence.Value{cadence.NewUInt64(2),cadence.NewUInt64(1)})
		assert.EqualValues(t, idsArray, result)

		buyerBalance := test.ExecuteScriptAndCheck(
			t, b,
			metayaUtilityCoin.GetBalanceScript(contracts.FungibleTokenAddress.String(), contracts.MetayaUtilityCoinAddress.String()),
			[][]byte{jsoncdc.MustEncode(cadence.Address(buyerAddress))},
		)
		assert.EqualValues(t, test.CadenceUFix64("0.0"), buyerBalance)

		copyrightOwner1Balance := test.ExecuteScriptAndCheck(
			t, b,
			metayaUtilityCoin.GetBalanceScript(contracts.FungibleTokenAddress.String(), contracts.MetayaUtilityCoinAddress.String()),
			[][]byte{jsoncdc.MustEncode(cadence.Address(copyrightOwner1Address))},
		)
		assert.EqualValues(t, test.CadenceUFix64("2.0"), copyrightOwner1Balance)

		copyrightOwner2Balance := test.ExecuteScriptAndCheck(
			t, b,
			metayaUtilityCoin.GetBalanceScript(contracts.FungibleTokenAddress.String(), contracts.MetayaUtilityCoinAddress.String()),
			[][]byte{jsoncdc.MustEncode(cadence.Address(copyrightOwner2Address))},
		)
		assert.EqualValues(t, test.CadenceUFix64("8.0"), copyrightOwner2Balance)

		tempMetayaBalance := test.ExecuteScriptAndCheck(
			t, b,
			metayaUtilityCoin.GetBalanceScript(contracts.FungibleTokenAddress.String(), contracts.MetayaUtilityCoinAddress.String()),
			[][]byte{jsoncdc.MustEncode(cadence.Address(tempMetayaAddress))},
		)
		assert.EqualValues(t, test.CadenceUFix64("6.0"), tempMetayaBalance)

		sellerBalance := test.ExecuteScriptAndCheck(
			t, b,
			metayaUtilityCoin.GetBalanceScript(contracts.FungibleTokenAddress.String(), contracts.MetayaUtilityCoinAddress.String()),
			[][]byte{jsoncdc.MustEncode(cadence.Address(sellerAddress))},
		)
		assert.EqualValues(t, test.CadenceUFix64("184.0"), sellerBalance)

		supply := test.ExecuteScriptAndCheck(
			t, b,
			metayaUtilityCoin.GetSupplyScript(contracts.FungibleTokenAddress.String(), contracts.MetayaUtilityCoinAddress.String()),
			nil,
		)
		assert.EqualValues(t, test.CadenceUFix64("200.0"), supply)
	})
}
