package test

import (
	"testing"

	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	sdktemplates "github.com/onflow/flow-go-sdk/templates"
	sdktest "github.com/onflow/flow-go-sdk/test"
	flowStorageFees "github.com/onflow/flow-core-contracts/lib/go/contracts"
	"github.com/stretchr/testify/assert"
	

	"github.com/0xRaito/metaya-smart-contracts/tests/go/test"
	"github.com/0xRaito/metaya-smart-contracts/tests/go/metayaNFT"
	"github.com/0xRaito/metaya-smart-contracts/tests/go/soundlinksDID"
	"github.com/0xRaito/metaya-smart-contracts/tests/go/nftStorefront"
	"github.com/0xRaito/metaya-smart-contracts/tests/go/metayaUtilityCoin"
	"github.com/0xRaito/metaya-smart-contracts/tests/go/metayaBeneficiaryCut"
)

// This test is for testing the deployment the Metaya contracts
func TestMetayaNFTDeployContracts(t *testing.T) {
	b := test.NewBlockchain()
	
	contracts := metaya.DeployContracts(t, b)

	// Deploy the sharded collection contract
	metayaShardedCollectionCode := metaya.LoadMetayaShardedCollection(contracts.NonFungibleTokenAddress.String(), contracts.MetayaAddress.String())
	metayaShardedCollectionAddress, err := b.CreateAccount(nil, []sdktemplates.Contract{
		{
			Name:   "MetayaShardedCollection",
			Source: string(metayaShardedCollectionCode),
		},
	})
	assert.NoError(t, err)

	_, err = b.CommitBlock()
	assert.NoError(t, err)

	// Should be able to deploy the admin receiver contract as a new account with no keys.
	metayaAdminReceiverCode := metaya.LoadMetayaAdminReceiver(contracts.MetayaAddress.String(), metayaShardedCollectionAddress.String())
	_, err = b.CreateAccount(nil, []sdktemplates.Contract{
		{
			Name:   "MetayaAdminReceiver",
			Source: string(metayaAdminReceiverCode),
		},
	})
	assert.NoError(t, err)

	_, err = b.CommitBlock()
	assert.NoError(t, err)
}

// This test tests the pure functionality of the smart contract
func TestMetayaNFTMintNFTs(t *testing.T) {
	b := test.NewBlockchain()

	contracts := metaya.DeployContracts(t, b)

	soundlinksDID.SetupAccount(
		t, b,
		contracts.MetayaAddress,
		contracts.MetayaSigner,
		contracts.NonFungibleTokenAddress,
		contracts.SoundlinksDIDAddress,
	)

	// Check that that main contract fields were initialized correctly
	result := test.ExecuteScriptAndCheck(
		t, b,
		metaya.GetCurrentSeriesScript(contracts),
		nil,
	)
	assert.EqualValues(t, cadence.NewUInt32(0), result)

	result = test.ExecuteScriptAndCheck(
		t, b,
		metaya.GetNextPlayIDScript(contracts),
		nil,
	)
	assert.EqualValues(t, cadence.NewUInt32(1), result)

	result = test.ExecuteScriptAndCheck(
		t, b,
		metaya.GetNextSetIDScript(contracts),
		nil,
	)
	assert.EqualValues(t, cadence.NewUInt32(1), result)

	result = test.ExecuteScriptAndCheck(
		t, b,
		metaya.GetTotalSupplyScript(contracts),
		nil,
	)
	assert.EqualValues(t, cadence.NewUInt64(0), result)

	// Deploy the sharded collection contract
	metayaShardedCollectionCode := metaya.LoadMetayaShardedCollection(contracts.NonFungibleTokenAddress.String(), contracts.MetayaAddress.String())
	metayaShardedCollectionAddress, err := b.CreateAccount(nil, []sdktemplates.Contract{
		{
			Name:   "MetayaShardedCollection",
			Source: string(metayaShardedCollectionCode),
		},
	})
	assert.NoError(t, err)

	_, err = b.CommitBlock()
	assert.NoError(t, err)

	// Deploy the flowStorageFees contract
	flowStorageFeesCode := flowStorageFees.FlowStorageFees(test.FTAddress.String(), test.FlowTokenAddress.String())
	flowStorageFeesAddress, err := b.CreateAccount(nil, []sdktemplates.Contract{
		{
			Name:   "FlowStorageFees",
			Source: string(flowStorageFeesCode),
		},
	})
	assert.NoError(t, err)

	_, err = b.CommitBlock()
	assert.NoError(t, err)

	// Deploy the nftStorefront contract
	nftStorefrontCode := nftStorefront.LoadNFTStorefront(
		test.FTAddress,
		contracts.NonFungibleTokenAddress,
	)

	nftStorefrontAddress, err := b.CreateAccount(nil, []sdktemplates.Contract{
		{
			Name:   "NFTStorefront",
			Source: string(nftStorefrontCode),
		},
	})
	assert.NoError(t, err)

	_, err = b.CommitBlock()
	assert.NoError(t, err)

	// Create a new user account
	accountKeys := sdktest.AccountKeyGenerator()
	userAccountKey, userSigner := accountKeys.NewWithSigner()
	userAddress, _ := b.CreateAccount([]*flow.AccountKey{userAccountKey}, nil)

	title := cadence.String("Title")
	artwork1 := cadence.String("Artwork 001")
	artwork2 := cadence.String("Artwork 002")
	artwork3 := cadence.String("Artwork 003")
	artwork4 := cadence.String("Artwork 004")

	// Admin sends a transaction to create a play
	t.Run("Should be able to create a new Play", func(t *testing.T) {
		metadata := []cadence.KeyValuePair{{Key: title, Value: artwork1}}
		play := cadence.NewDictionary(metadata)

		tx := flow.NewTransaction().
			SetScript(metaya.CreatePlayTransaction(contracts)).
			SetGasLimit(100).
			SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
			SetPayer(b.ServiceKey().Address).
			AddAuthorizer(contracts.MetayaAddress)

		_ = tx.AddArgument(play)

		test.SignAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, contracts.MetayaAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), contracts.MetayaSigner},
			false,
		)
	})

	// Admin sends transactions to create multiple plays
	t.Run("Should be able to create multiple new Plays", func(t *testing.T) {
		metadata := []cadence.KeyValuePair{{Key: title, Value: artwork2}}
		play := cadence.NewDictionary(metadata)

		tx := flow.NewTransaction().
			SetScript(metaya.CreatePlayTransaction(contracts)).
			SetGasLimit(100).
			SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
			SetPayer(b.ServiceKey().Address).
			AddAuthorizer(contracts.MetayaAddress)

		_ = tx.AddArgument(play)

		test.SignAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, contracts.MetayaAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), contracts.MetayaSigner},
			false,
		)

		metadata = []cadence.KeyValuePair{{Key: title, Value: artwork3}}
		play = cadence.NewDictionary(metadata)

		tx = flow.NewTransaction().
			SetScript(metaya.CreatePlayTransaction(contracts)).
			SetGasLimit(100).
			SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
			SetPayer(b.ServiceKey().Address).
			AddAuthorizer(contracts.MetayaAddress)

		_ = tx.AddArgument(play)

		test.SignAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, contracts.MetayaAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), contracts.MetayaSigner},
			false,
		)

		metadata = []cadence.KeyValuePair{{Key: title, Value: artwork4}}
		play = cadence.NewDictionary(metadata)

		tx = flow.NewTransaction().
			SetScript(metaya.CreatePlayTransaction(contracts)).
			SetGasLimit(100).
			SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
			SetPayer(b.ServiceKey().Address).
			AddAuthorizer(contracts.MetayaAddress)

		_ = tx.AddArgument(play)

		test.SignAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, contracts.MetayaAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), contracts.MetayaSigner},
			false,
		)

		// Check that the return all plays script doesn't fail
		// and that we can return metadata about the plays
		test.ExecuteScriptAndCheck(
			t, b,
			metaya.GetAllPlaysScript(contracts),
			nil,
		)

		result = test.ExecuteScriptAndCheck(
			t, b,
			metaya.GetPlayMetadataFieldScript(contracts),
			[][]byte{jsoncdc.MustEncode(cadence.UInt32(1)), jsoncdc.MustEncode(cadence.String("Title"))},
		)
		assert.EqualValues(t, cadence.String("Artwork 001"), result)

		result = test.ExecuteScriptAndCheck(
			t, b,
			metaya.GetPlayMetadataScript(contracts),
			[][]byte{jsoncdc.MustEncode(cadence.UInt32(2))},
		)

		metadata = []cadence.KeyValuePair{{Key: title, Value: artwork2}}
		assert.EqualValues(t, cadence.NewDictionary(metadata), result)
	})

	// Admin creates a new Set with the name Set 001
	t.Run("Should be able to create a new Set", func(t *testing.T) {
		tx := flow.NewTransaction().
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

		// Check that the set name, ID, and series were initialized correctly.
		result = test.ExecuteScriptAndCheck(
			t, b,
			metaya.GetSetNameScript(contracts),
			[][]byte{jsoncdc.MustEncode(cadence.UInt32(1))},
		)
		assert.EqualValues(t, cadence.String("Set 001"), result)

		result = test.ExecuteScriptAndCheck(
			t, b,
			metaya.GetSetIDsByNameScript(contracts),
			[][]byte{jsoncdc.MustEncode(cadence.String("Set 001"))},
		)
		idsArray := cadence.NewArray([]cadence.Value{cadence.NewUInt32(1)})
		assert.EqualValues(t, idsArray, result)

		result = test.ExecuteScriptAndCheck(
			t, b,
			metaya.GetSetSeriesScript(contracts),
			[][]byte{jsoncdc.MustEncode(cadence.UInt32(1))},
		)
		assert.EqualValues(t, cadence.NewUInt32(0), result)
	})

	t.Run("Should not be able to create play data struct that increment the id counter", func(t *testing.T) {
		tx := flow.NewTransaction().
			SetScript(metaya.CreatePlayStructTransaction(contracts)).
			SetGasLimit(100).
			SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
			SetPayer(b.ServiceKey().Address).
			AddAuthorizer(contracts.MetayaAddress)

		test.SignAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, contracts.MetayaAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), contracts.MetayaSigner},
			false,
		)

		// Check that the play ID and set ID were not incremented
		result = test.ExecuteScriptAndCheck(
			t, b,
			metaya.GetNextPlayIDScript(contracts),
			nil,
		)
		assert.EqualValues(t, cadence.NewUInt32(5), result)
	})

	// Admin sends a transaction that adds play 1 to the set
	t.Run("Should be able to add a play to a Set", func(t *testing.T) {
		tx := flow.NewTransaction().
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
	})

	// Admin sends a transaction that adds plays 2 and 3 to the set
	t.Run("Should be able to add multiple plays to a Set", func(t *testing.T) {
		tx := flow.NewTransaction().
			SetScript(metaya.AddPlaysToSetTransaction(contracts)).
			SetGasLimit(100).
			SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
			SetPayer(b.ServiceKey().Address).
			AddAuthorizer(contracts.MetayaAddress)

		_ = tx.AddArgument(cadence.NewUInt32(1))
		plays := []cadence.Value{cadence.NewUInt32(2), cadence.NewUInt32(3)}
		_ = tx.AddArgument(cadence.NewArray(plays))

		test.SignAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, contracts.MetayaAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), contracts.MetayaSigner},
			false,
		)

		// Make sure the plays were added correctly and the edition isn't retired or locked
		result = test.ExecuteScriptAndCheck(
			t, b,
			metaya.GetPlaysInSetScript(contracts),
			[][]byte{jsoncdc.MustEncode(cadence.UInt32(1))},
		)
		playsArray := cadence.NewArray([]cadence.Value{cadence.NewUInt32(1), cadence.NewUInt32(2), cadence.NewUInt32(3)})
		assert.EqualValues(t, playsArray, result)

		result = test.ExecuteScriptAndCheck(
			t, b,
			metaya.GetEditionRetiredScript(contracts),
			[][]byte{jsoncdc.MustEncode(cadence.UInt32(1)), jsoncdc.MustEncode(cadence.UInt32(1))},
		)
		assert.EqualValues(t, cadence.NewBool(false), result)

		result = test.ExecuteScriptAndCheck(
			t, b,
			metaya.GetSetLockedScript(contracts),
			[][]byte{jsoncdc.MustEncode(cadence.UInt32(1))},
		)
		assert.EqualValues(t, cadence.NewBool(false), result)
	})

	// Admin sends a transaction that locks the set
	t.Run("Should be able to lock a set which stops plays from being added", func(t *testing.T) {
		tx := flow.NewTransaction().
			SetScript(metaya.LockSetTransaction(contracts)).
			SetGasLimit(100).
			SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
			SetPayer(b.ServiceKey().Address).
			AddAuthorizer(contracts.MetayaAddress)

		_ = tx.AddArgument(cadence.NewUInt32(1))

		test.SignAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, contracts.MetayaAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), contracts.MetayaSigner},
			false,
		)

		// This should fail because the set is locked
		tx = flow.NewTransaction().
			SetScript(metaya.AddPlayToSetTransaction(contracts)).
			SetGasLimit(100).
			SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
			SetPayer(b.ServiceKey().Address).
			AddAuthorizer(contracts.MetayaAddress)

		_ = tx.AddArgument(cadence.NewUInt32(1))
		_ = tx.AddArgument(cadence.NewUInt32(4))

		test.SignAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, contracts.MetayaAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), contracts.MetayaSigner},
			true,
		)

		// Script should return that the set is locked
		result = test.ExecuteScriptAndCheck(
			t, b,
			metaya.GetSetLockedScript(contracts),
			[][]byte{jsoncdc.MustEncode(cadence.UInt32(1))},
		)
		assert.EqualValues(t, cadence.NewBool(true), result)
	})

	// Admin mints a moment that stores it in the admin's collection
	t.Run("Should be able to mint a moment", func(t *testing.T) {
		hashs := []string{"FFFFFFFFFFFFFFFFFFF1"}
		soundlinksDID.MintDIDs(
			t, b,
			contracts.NonFungibleTokenAddress, contracts.SoundlinksDIDAddress, contracts.MetayaAddress,
			contracts.SoundlinksDIDSigner, 
			hashs,
		)

		setID := uint32(1)
		playID := uint32(1)
		metaya.MintNFT(
			t, b,
			contracts,
			setID, playID,
			contracts.MetayaAddress,
			false,
		)

		// Make sure the moment was minted correctly and is stored in the collection with the correct data
		result = test.ExecuteScriptAndCheck(
			t, b,
			metaya.GetIdInCollectionScript(contracts),
			[][]byte{jsoncdc.MustEncode(cadence.Address(contracts.MetayaAddress)), jsoncdc.MustEncode(cadence.UInt64(1))},
		)
		assert.EqualValues(t, cadence.NewBool(true), result)

		result = test.ExecuteScriptAndCheck(
			t, b,
			metaya.GetCollectionIdsScript(contracts),
			[][]byte{jsoncdc.MustEncode(cadence.Address(contracts.MetayaAddress))},
		)
		idsArray := cadence.NewArray([]cadence.Value{cadence.NewUInt64(1)})
		assert.EqualValues(t, idsArray, result)

		result = test.ExecuteScriptAndCheck(
			t, b,
			metaya.GetMomentSetIDScript(contracts),
			[][]byte{jsoncdc.MustEncode(cadence.Address(contracts.MetayaAddress)), jsoncdc.MustEncode(cadence.UInt64(1))},
		)
		assert.EqualValues(t, cadence.NewUInt32(1), result)

		result = test.ExecuteScriptAndCheck(
			t, b,
			metaya.GetMetadataScript(contracts),
			[][]byte{jsoncdc.MustEncode(cadence.Address(contracts.MetayaAddress)), jsoncdc.MustEncode(cadence.UInt64(1))},
		)
		metadata := []cadence.KeyValuePair{{Key: title, Value: artwork1}}
		assert.EqualValues(t, cadence.NewDictionary(metadata), result)

		result = test.ExecuteScriptAndCheck(
			t, b,
			metaya.GetMetadataFieldScript(contracts),
			[][]byte{jsoncdc.MustEncode(cadence.Address(contracts.MetayaAddress)), jsoncdc.MustEncode(cadence.UInt64(1)), jsoncdc.MustEncode(cadence.String("Title"))},
		)
		assert.EqualValues(t, cadence.String("Artwork 001"), result)
	})

	// Admin sends a transaction that mints a batch of moments
	t.Run("Should be able to mint a batch of moments", func(t *testing.T) {
		hashs := []string{"FFFFFFFFFFFFFFFFFFF2","FFFFFFFFFFFFFFFFFFF3","FFFFFFFFFFFFFFFFFFF4","FFFFFFFFFFFFFFFFFFF5","FFFFFFFFFFFFFFFFFFF6"}
		soundlinksDID.MintDIDs(
			t, b,
			contracts.NonFungibleTokenAddress, contracts.SoundlinksDIDAddress, contracts.MetayaAddress,
			contracts.SoundlinksDIDSigner,
			hashs,
		)

		setID := uint32(1)
		playID := uint32(3)
		quantity := uint32(5)
		metaya.MintNFTs(
		 	t, b,
		 	contracts,
		 	setID, playID,
		 	quantity,
		 	contracts.MetayaAddress,
			false,
		)

		// Ensure that the correct number of moments have been minted for each edition
		result = test.ExecuteScriptAndCheck(
			t, b,
			metaya.GetNumMomentsInEditionScript(contracts),
			[][]byte{jsoncdc.MustEncode(cadence.UInt32(1)), jsoncdc.MustEncode(cadence.UInt32(1))},
		)
		assert.EqualValues(t, cadence.NewUInt32(1), result)

		result = test.ExecuteScriptAndCheck(
			t, b,
			metaya.GetNumMomentsInEditionScript(contracts),
			[][]byte{jsoncdc.MustEncode(cadence.UInt32(1)), jsoncdc.MustEncode(cadence.UInt32(3))},
		)
		assert.EqualValues(t, cadence.NewUInt32(5), result)

		result = test.ExecuteScriptAndCheck(
			t, b,
			metaya.GetIdInCollectionScript(contracts),
			[][]byte{jsoncdc.MustEncode(cadence.Address(contracts.MetayaAddress)), jsoncdc.MustEncode(cadence.UInt64(1))},
		)
		assert.EqualValues(t, cadence.NewBool(true), result)

		result = test.ExecuteScriptAndCheck(
			t, b,
			metaya.GetCollectionIdsScript(contracts),
			[][]byte{jsoncdc.MustEncode(cadence.Address(contracts.MetayaAddress))},
		)
		idsArray := cadence.NewArray([]cadence.Value{cadence.NewUInt64(4), cadence.NewUInt64(3), cadence.NewUInt64(1), cadence.NewUInt64(6), cadence.NewUInt64(2), cadence.NewUInt64(5)})
		assert.EqualValues(t, idsArray, result)

		result = test.ExecuteScriptAndCheck(
			t, b,
			metaya.GetMomentSetIDScript(contracts),
			[][]byte{jsoncdc.MustEncode(cadence.Address(contracts.MetayaAddress)), jsoncdc.MustEncode(cadence.UInt64(1))},
		)
		assert.EqualValues(t, cadence.NewUInt32(1), result)

		result = test.ExecuteScriptAndCheck(
			t, b,
			metaya.GetMomentSerialNumScript(contracts),
			[][]byte{jsoncdc.MustEncode(cadence.Address(contracts.MetayaAddress)), jsoncdc.MustEncode(cadence.UInt64(3))},
		)
		assert.EqualValues(t, cadence.UInt32(2), result)

		result = test.ExecuteScriptAndCheck(
			t, b,
			metaya.GetMomentSetNameScript(contracts),
			[][]byte{jsoncdc.MustEncode(cadence.Address(contracts.MetayaAddress)), jsoncdc.MustEncode(cadence.UInt64(3))},
		)
		assert.EqualValues(t, cadence.String("Set 001"), result)
	})

	t.Run("Should be able to mint a batch of moments and fulfill a pack", func(t *testing.T) {
		hashs := []string{"FFFFFFFFFFFFFFFFFFF7","FFFFFFFFFFFFFFFFFFF8","FFFFFFFFFFFFFFFFFFF9","FFFFFFFFFFFFFFFFFFF10","FFFFFFFFFFFFFFFFFFF11"}
		soundlinksDID.MintDIDs(
			t, b,
			contracts.NonFungibleTokenAddress, contracts.SoundlinksDIDAddress, contracts.MetayaAddress,
			contracts.SoundlinksDIDSigner,
			hashs,
		)

		setID := uint32(1)
		playID := uint32(3)
		quantity := uint32(5)
		metaya.MintNFTs(
		 	t, b,
		 	contracts,
		 	setID, playID,
		 	quantity,
		 	contracts.MetayaAddress,
			false,
		)

		tx := flow.NewTransaction().
			SetScript(metaya.ProvideMomentTransaction(contracts, metayaShardedCollectionAddress.String(), flowStorageFeesAddress.String())).
			SetGasLimit(300).
			SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
			SetPayer(b.ServiceKey().Address).
			AddAuthorizer(contracts.MetayaAddress)

		_ = tx.AddArgument(cadence.NewAddress(contracts.MetayaAddress))
		ids := []cadence.Value{cadence.NewUInt64(6), cadence.NewUInt64(7), cadence.NewUInt64(8)}
		_ = tx.AddArgument(cadence.NewArray(ids))

		test.SignAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, contracts.MetayaAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), contracts.MetayaSigner},
			false,
		)
	})

	// Admin sends a transaction to retire a play
	t.Run("Should be able to retire a Play which stops minting", func(t *testing.T) {
		tx := flow.NewTransaction().
			SetScript(metaya.RetirePlayFromSetTransaction(contracts)).
			SetGasLimit(300).
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

		// Minting from this play should fail becuase it is retired
		hashs := []string{"FFFFFFFFFFFFFFFFFFF1"}
		soundlinksDID.MintDIDs(
			t, b,
			contracts.NonFungibleTokenAddress, contracts.SoundlinksDIDAddress, contracts.MetayaAddress,
			contracts.SoundlinksDIDSigner, 
			hashs,
		)

		setID := uint32(1)
		playID := uint32(1)
		metaya.MintNFT(
			t, b,
			contracts,
			setID, playID,
			contracts.MetayaAddress,
			true,
		)

		// Make sure this edition is retired
		result = test.ExecuteScriptAndCheck(
			t, b,
			metaya.GetEditionRetiredScript(contracts),
			[][]byte{jsoncdc.MustEncode(cadence.UInt32(1)), jsoncdc.MustEncode(cadence.UInt32(1))},
		)
		assert.EqualValues(t, cadence.NewBool(true), result)
	})

	// Admin sends a transaction that retires all the plays in a set
	t.Run("Should be able to retire all Plays which stops minting", func(t *testing.T) {
		tx := flow.NewTransaction().
			SetScript(metaya.RetireAllPlaysFromSetTransaction(contracts)).
			SetGasLimit(300).
			SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
			SetPayer(b.ServiceKey().Address).
			AddAuthorizer(contracts.MetayaAddress)

		_ = tx.AddArgument(cadence.NewUInt32(1))

		test.SignAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, contracts.MetayaAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), contracts.MetayaSigner},
			false,
		)

		// minting should fail
		hashs := []string{"FFFFFFFFFFFFFFFFFFF1"}
		soundlinksDID.MintDIDs(
			t, b,
			contracts.NonFungibleTokenAddress, contracts.SoundlinksDIDAddress, contracts.MetayaAddress,
			contracts.SoundlinksDIDSigner, 
			hashs,
		)

		setID := uint32(1)
		playID := uint32(3)
		metaya.MintNFT(
			t, b,
			contracts,
			setID, playID,
			contracts.MetayaAddress,
			true,
		)

		metaya.VerifyQuerySetMetadata(t, b, contracts,
			metaya.SetMetadata{
				SetID:  1,
				Name:   "Set 001",
				Series: 0,
				Plays:  []uint32{1, 2, 3},
				//retired {UInt32: Bool}
				Locked: true,
				//numberMintedPerPlay {UInt32: UInt32}})
			})
	})

	// Create a new Collection for a user address
	t.Run("Should be able to create a moment Collection", func(t *testing.T) {
		metaya.SetupAccount(
			t, b,
			userAddress,
			userSigner,
			contracts.MetayaAddress,
		)
	})

	// Admin sends a transaction to transfer a moment to a user
	t.Run("Should be able to transfer a moment", func(t *testing.T) {
		metaya.TransferNFT(
			t, b,
			contracts,
			contracts.MetayaAddress,
			contracts.MetayaSigner,
			nftStorefrontAddress,
			uint64(1),
			userAddress,
			false,
		)

		// Make sure the user received it
		result = test.ExecuteScriptAndCheck(
			t, b,
			metaya.GetIdInCollectionScript(contracts),
			[][]byte{jsoncdc.MustEncode(cadence.Address(userAddress)), jsoncdc.MustEncode(cadence.UInt64(1))},
		)
		assert.EqualValues(t, cadence.NewBool(true), result)
	})

	// Admin sends a transaction to transfer a batch of moments to a user
	t.Run("Should be able to transfer a batch of moments", func(t *testing.T) {
		metaya.TransferNFTs(
			t, b,
			contracts,
			contracts.MetayaAddress,
			contracts.MetayaSigner,
			nftStorefrontAddress,
			[]uint64{2, 3},
			userAddress,
			false,
		)

		// Make sure the user received it
		result = test.ExecuteScriptAndCheck(
			t, b,
			metaya.GetIdInCollectionScript(contracts),
			[][]byte{jsoncdc.MustEncode(cadence.Address(userAddress)), jsoncdc.MustEncode(cadence.UInt64(2))},
		)
		assert.EqualValues(t, cadence.NewBool(true), result)

		result = test.ExecuteScriptAndCheck(
			t, b,
			metaya.GetIdInCollectionScript(contracts),
			[][]byte{jsoncdc.MustEncode(cadence.Address(userAddress)), jsoncdc.MustEncode(cadence.UInt64(3))},
		)
		assert.EqualValues(t, cadence.NewBool(true), result)
	})

	// Admin sends a transaction that creates a new sharded collection for the admin
	t.Run("Should be able to create new sharded moment collection and store it", func(t *testing.T) {
		tx := flow.NewTransaction().
			SetScript(metaya.SetupShardedCollectionTransaction(contracts, metayaShardedCollectionAddress.String())).
			SetGasLimit(400).
			SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
			SetPayer(b.ServiceKey().Address).
			AddAuthorizer(contracts.MetayaAddress)

		_ = tx.AddArgument(cadence.NewUInt64(32))

		test.SignAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, contracts.MetayaAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), contracts.MetayaSigner},
			false,
		)
	})

	// Admin sends a transaction to transfer a batch of moments to a user
	t.Run("Should be able to transfer moments from a sharded collection", func(t *testing.T) {
		// Create a new set 002
		tx := flow.NewTransaction().
			SetScript(metaya.CreateSetTransaction(contracts)).
			SetGasLimit(100).
			SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
			SetPayer(b.ServiceKey().Address).
			AddAuthorizer(contracts.MetayaAddress)

		_ = tx.AddArgument(cadence.String("Set 002"))

		test.SignAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, contracts.MetayaAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), contracts.MetayaSigner},
			false,
		)

		// Add play 4 to the set 002
		tx = flow.NewTransaction().
			SetScript(metaya.AddPlayToSetTransaction(contracts)).
			SetGasLimit(100).
			SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
			SetPayer(b.ServiceKey().Address).
			AddAuthorizer(contracts.MetayaAddress)

		_ = tx.AddArgument(cadence.NewUInt32(2))
		_ = tx.AddArgument(cadence.NewUInt32(4))

		test.SignAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, contracts.MetayaAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), contracts.MetayaSigner},
			false,
		)

		hashs := []string{"FFFFFFFFFFFFFFFFFFF12","FFFFFFFFFFFFFFFFFFF13","FFFFFFFFFFFFFFFFFFF14"}
		soundlinksDID.MintDIDs(
			t, b,
			contracts.NonFungibleTokenAddress, contracts.SoundlinksDIDAddress, contracts.MetayaAddress,
			contracts.SoundlinksDIDSigner,
			hashs,
		)

		setID := uint32(2)
		playID := uint32(4)
		quantity := uint32(3)
		metaya.MintNFTs(
		 	t, b,
		 	contracts,
		 	setID, playID,
		 	quantity,
		 	contracts.MetayaAddress,
			false,
		)

		// Transfer a moment
		tx = flow.NewTransaction().
			SetScript(metaya.TransferFromShardedTransaction(contracts, metayaShardedCollectionAddress.String())).
			SetGasLimit(300).
			SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
			SetPayer(b.ServiceKey().Address).
			AddAuthorizer(contracts.MetayaAddress)

		_ = tx.AddArgument(cadence.NewAddress(userAddress))
		_ = tx.AddArgument(cadence.NewUInt64(12))

		test.SignAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, contracts.MetayaAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), contracts.MetayaSigner},
			false,
		)

		// Make sure the user received them
		result = test.ExecuteScriptAndCheck(
			t, b,
			metaya.GetMomentSeriesScript(contracts),
			[][]byte{jsoncdc.MustEncode(cadence.Address(userAddress)), jsoncdc.MustEncode(cadence.UInt64(12))},
		)
		assert.EqualValues(t, cadence.NewUInt32(0), result)

		// Batch transfer moments
		tx = flow.NewTransaction().
			SetScript(metaya.BatchTransferFromShardedTransaction(contracts, metayaShardedCollectionAddress.String())).
			SetGasLimit(300).
			SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
			SetPayer(b.ServiceKey().Address).
			AddAuthorizer(contracts.MetayaAddress)

		_ = tx.AddArgument(cadence.NewAddress(userAddress))
		ids := []cadence.Value{cadence.NewUInt64(13), cadence.NewUInt64(14)}
		_ = tx.AddArgument(cadence.NewArray(ids))

		test.SignAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, contracts.MetayaAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), contracts.MetayaSigner},
			false,
		)

		// Make sure the user received them
		result = test.ExecuteScriptAndCheck(
			t, b,
			metaya.GetMomentPlayIDScript(contracts),
			[][]byte{jsoncdc.MustEncode(cadence.Address(userAddress)), jsoncdc.MustEncode(cadence.UInt64(13))},
		)
		assert.EqualValues(t, cadence.NewUInt32(4), result)

		result = test.ExecuteScriptAndCheck(
			t, b,
			metaya.GetMomentSetIDScript(contracts),
			[][]byte{jsoncdc.MustEncode(cadence.Address(userAddress)), jsoncdc.MustEncode(cadence.UInt64(14))},
		)
		assert.EqualValues(t, cadence.NewUInt32(2), result)
	})

	// Admin sends a transaction to update the current series
	t.Run("Should be able to change the current series", func(t *testing.T) {
		tx := flow.NewTransaction().
			SetScript(metaya.StartNewSeriesTransaction(contracts)).
			SetGasLimit(100).
			SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
			SetPayer(b.ServiceKey().Address).
			AddAuthorizer(contracts.MetayaAddress)

		test.SignAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, contracts.MetayaAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), contracts.MetayaSigner},
			false,
		)
	})

	// Make sure the contract fields are correct
	result = test.ExecuteScriptAndCheck(
		t, b,
		metaya.GetCurrentSeriesScript(contracts),
		nil,
	)
	assert.EqualValues(t, cadence.NewUInt32(1), result)

	result = test.ExecuteScriptAndCheck(
		t, b,
		metaya.GetNextPlayIDScript(contracts),
		nil,
	)
	assert.EqualValues(t, cadence.NewUInt32(5), result)

	result = test.ExecuteScriptAndCheck(
		t, b,
		metaya.GetNextSetIDScript(contracts),
		nil,
	)
	assert.EqualValues(t, cadence.NewUInt32(3), result)

	result = test.ExecuteScriptAndCheck(
		t, b,
		metaya.GetTotalSupplyScript(contracts),
		nil,
	)
	assert.EqualValues(t, cadence.NewUInt64(14), result)
}

// This test is for ensuring that admin receiver smart contract works correctly
func TestMetayaNFTTransferAdmin(t *testing.T) {
	b := test.NewBlockchain()

	accountKeys := sdktest.AccountKeyGenerator()

	contracts := metaya.DeployContracts(t, b)

	soundlinksDID.SetupAccount(
		t, b,
		contracts.MetayaAddress,
		contracts.MetayaSigner,
		contracts.NonFungibleTokenAddress,
		contracts.SoundlinksDIDAddress,
	)

	// Deploy the sharded collection contract
	metayaShardedCollectionCode := metaya.LoadMetayaShardedCollection(contracts.NonFungibleTokenAddress.String(), contracts.MetayaAddress.String())
	metayaShardedCollectionAddress, err := b.CreateAccount(nil, []sdktemplates.Contract{
		{
			Name:   "MetayaShardedCollection",
			Source: string(metayaShardedCollectionCode),
		},
	})
	assert.NoError(t, err)

	_, err = b.CommitBlock()
	assert.NoError(t, err)

	// Should be able to deploy the admin receiver contract
	metayaAdminReceiverCode := metaya.LoadMetayaAdminReceiver(contracts.MetayaAddress.String(), metayaShardedCollectionAddress.String())
	adminAccountKey, adminSigner := accountKeys.NewWithSigner()
	adminAddress, err := b.CreateAccount([]*flow.AccountKey{adminAccountKey}, []sdktemplates.Contract{
		{
			Name:   "MetayaAdminReceiver",
			Source: string(metayaAdminReceiverCode),
		},
	})
	assert.NoError(t, err)

	_, err = b.CommitBlock()
	assert.NoError(t, err)

	title := cadence.String("Title")
	artwork1 := cadence.String("Artwork 001")

	// Create a new Collection
	t.Run("Should be able to transfer an admin Capability to the receiver account", func(t *testing.T) {
		tx := flow.NewTransaction().
			SetScript(metaya.TransferAdminTransaction(contracts, adminAddress.String())).
			SetGasLimit(100).
			SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
			SetPayer(b.ServiceKey().Address).
			AddAuthorizer(contracts.MetayaAddress)

		test.SignAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, contracts.MetayaAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), contracts.MetayaSigner},
			false,
		)
	})

	// Can create a new play with the new admin
	t.Run("Should be able to create a new Play with the new Admin account", func(t *testing.T) {
		tx := flow.NewTransaction().
			SetScript(metaya.CreatePlayTransaction(contracts)).
			SetGasLimit(100).
			SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
			SetPayer(b.ServiceKey().Address).
			AddAuthorizer(adminAddress)

		metadata := []cadence.KeyValuePair{{Key: title, Value: artwork1}}
		play := cadence.NewDictionary(metadata)
		_ = tx.AddArgument(play)

		test.SignAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, adminAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), adminSigner},
			false,
		)
	})
}

func TestMetayaNFTSetPlaysOwnedByAddress(t *testing.T) {
	b := test.NewBlockchain()

	accountKeys := sdktest.AccountKeyGenerator()

	contracts := metaya.DeployContracts(t, b)

	// Create a new user account
	userAccountKey, userSigner := accountKeys.NewWithSigner()
	userAddress, _ := b.CreateAccount([]*flow.AccountKey{userAccountKey}, nil)

	// Create moment collection
	metaya.SetupAccount(
		t, b,
		userAddress,
		userSigner,
		contracts.MetayaAddress,
	)

	title := cadence.String("Title")
	artwork1 := cadence.String("Artwork 001")
	artwork2 := cadence.String("Artwork 002")
	artwork3 := cadence.String("Artwork 003")

	// Create plays
	artwork1_PlayID := uint32(1)

	tx := flow.NewTransaction().
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

	artwork2_PlayID := uint32(2)

	tx = flow.NewTransaction().
		SetScript(metaya.CreatePlayTransaction(contracts)).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(contracts.MetayaAddress)

	metadata = []cadence.KeyValuePair{{Key: title, Value: artwork2}}
	play = cadence.NewDictionary(metadata)
	_ = tx.AddArgument(play)

	test.SignAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, contracts.MetayaAddress},
		[]crypto.Signer{b.ServiceKey().Signer(), contracts.MetayaSigner},
		false,
	)

	artwork3_PlayID := uint32(3)

	tx = flow.NewTransaction().
		SetScript(metaya.CreatePlayTransaction(contracts)).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(contracts.MetayaAddress)

	metadata = []cadence.KeyValuePair{{Key: title, Value: artwork3}}
	play = cadence.NewDictionary(metadata)
	_ = tx.AddArgument(play)

	test.SignAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, contracts.MetayaAddress},
		[]crypto.Signer{b.ServiceKey().Signer(), contracts.MetayaSigner},
		false,
	)

	// Create Set
	setID := uint32(1)

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
		SetScript(metaya.AddPlaysToSetTransaction(contracts)).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(contracts.MetayaAddress)

	_ = tx.AddArgument(cadence.NewUInt32(setID))
	plays := []cadence.Value{cadence.NewUInt32(artwork1_PlayID), cadence.NewUInt32(artwork2_PlayID), cadence.NewUInt32(artwork3_PlayID)}
	_ = tx.AddArgument(cadence.NewArray(plays))

	test.SignAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, contracts.MetayaAddress},
		[]crypto.Signer{b.ServiceKey().Signer(), contracts.MetayaSigner},
		false,
	)

	// Mint moments to userAddress
	hashs := []string{"FFFFFFFFFFFFFFFFFFFA","FFFFFFFFFFFFFFFFFFFB","FFFFFFFFFFFFFFFFFFFC"}

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
	
	metaya.MintNFT(
		t, b,
		contracts,
		setID, artwork1_PlayID,
		userAddress,
		false,
	)

	metaya.MintNFT(
		t, b,
		contracts,
		setID, artwork2_PlayID,
		userAddress,
	    false,
    )

	// Mint one moment to topshotAddress
	metaya.MintNFT(
		t, b,
		contracts,
		setID, artwork1_PlayID,
		contracts.MetayaAddress,
	    false,
    )

	t.Run("Should return true if the address owns moments corresponding to each SetPlay", func(t *testing.T) {
		setIDs := cadence.NewArray([]cadence.Value{cadence.NewUInt32(setID), cadence.NewUInt32(setID)})
		playIDs := cadence.NewArray([]cadence.Value{cadence.NewUInt32(artwork1_PlayID), cadence.NewUInt32(artwork2_PlayID)})

		result := test.ExecuteScriptAndCheck(
			t, b,
			metaya.GetSetplaysAreOwnedScript(contracts),
			[][]byte{jsoncdc.MustEncode(cadence.Address(userAddress)), jsoncdc.MustEncode(setIDs), jsoncdc.MustEncode(playIDs)},
		)
		assert.EqualValues(t, cadence.NewBool(true), result)
	})

	t.Run("Should return false if the address does not own moments corresponding to each SetPlay", func(t *testing.T) {
		setIDs := cadence.NewArray([]cadence.Value{cadence.NewUInt32(setID), cadence.NewUInt32(setID), cadence.NewUInt32(setID)})
		playIDs := cadence.NewArray([]cadence.Value{cadence.NewUInt32(artwork1_PlayID), cadence.NewUInt32(artwork2_PlayID), cadence.NewUInt32(artwork3_PlayID)})

		result := test.ExecuteScriptAndCheck(
			t, b,
			metaya.GetSetplaysAreOwnedScript(contracts),
			[][]byte{jsoncdc.MustEncode(cadence.Address(userAddress)), jsoncdc.MustEncode(setIDs), jsoncdc.MustEncode(playIDs)},
		)
		assert.EqualValues(t, cadence.NewBool(false), result)
	})
}

func TestMetayaNFTPurchaseNFTs(t *testing.T) {
	b := test.NewBlockchain()

	accountKeys := sdktest.AccountKeyGenerator()

	fungibleTokenAddress, metayaUtilityCoinAddress, metayaUtilityCoinSigner, metayaBeneficiaryCutAddress, metayaBeneficiaryCutSigner := metayaBeneficiaryCut.DeployContracts(t, b)

	// Create copyright owner and commonweal
	copyrightOwner1Address, copyrightOwner1Signer, _ := test.CreateAccount(t, b)
	metayaUtilityCoin.SetupAccount(
		t, b,
		copyrightOwner1Address,
		copyrightOwner1Signer,
		fungibleTokenAddress,
		metayaUtilityCoinAddress,
	)
	metayaBeneficiaryCut.CreateCopyrightOwner(
		t, b,
		"Copyright Owner 001",
		copyrightOwner1Address,
		fungibleTokenAddress, metayaUtilityCoinAddress, metayaBeneficiaryCutAddress,
		metayaBeneficiaryCutSigner,
	)

	copyrightOwner2Address, copyrightOwner2Signer, _ := test.CreateAccount(t, b)
	metayaUtilityCoin.SetupAccount(
		t, b,
		copyrightOwner2Address,
		copyrightOwner2Signer,
		fungibleTokenAddress,
		metayaUtilityCoinAddress,
	)
	metayaBeneficiaryCut.CreateCopyrightOwner(
		t, b,
		"Copyright Owner 002",
		copyrightOwner2Address,
		fungibleTokenAddress, metayaUtilityCoinAddress, metayaBeneficiaryCutAddress,
		metayaBeneficiaryCutSigner,
	)

	commonwealAddress, commonwealSigner, _ := test.CreateAccount(t, b)
	metayaUtilityCoin.SetupAccount(
		t, b,
		commonwealAddress,
		commonwealSigner,
		fungibleTokenAddress,
		metayaUtilityCoinAddress,
	)
	metayaBeneficiaryCut.CreateCommonweal(
		t, b,
		"Commonweal 001",
		"0.002",
		commonwealAddress,
		fungibleTokenAddress, metayaUtilityCoinAddress, metayaBeneficiaryCutAddress,
		metayaBeneficiaryCutSigner,
	)

	tempMetayaAddress, tempMetayaSigner, _ := test.CreateAccount(t, b)
	metayaUtilityCoin.SetupAccount(
		t, b,
		tempMetayaAddress,
		tempMetayaSigner,
		fungibleTokenAddress,
		metayaUtilityCoinAddress,
	)
	metayaBeneficiaryCut.CreateCopyrightOwner(
		t, b,
		"Metaya",
		tempMetayaAddress,
		fungibleTokenAddress, metayaUtilityCoinAddress, metayaBeneficiaryCutAddress,
		metayaBeneficiaryCutSigner,
	)
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

	userAccountKey, userSigner := accountKeys.NewWithSigner()
	userAddress, _ := b.CreateAccount([]*flow.AccountKey{userAccountKey}, nil)
	metayaUtilityCoin.SetupAccount(
		t, b,
		userAddress,
		userSigner,
		fungibleTokenAddress,
		metayaUtilityCoinAddress,
	)

	// Set store cut percentage
	pair1 := cadence.KeyValuePair{Key: cadence.String("Copyright Owner 001"), Value: test.CadenceUFix64("0.5")}
	pair2 := cadence.KeyValuePair{Key: cadence.String("Copyright Owner 002"), Value: test.CadenceUFix64("0.3")}
	pair3 := cadence.KeyValuePair{Key: cadence.String("Metaya"), Value: test.CadenceUFix64("0.2")}
	recipientPairs := make([]cadence.KeyValuePair, 3)
	recipientPairs[0] = pair1
	recipientPairs[1] = pair2
	recipientPairs[2] = pair3

	tx = flow.NewTransaction().
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

	// Set pack cut percentage
	pair1 = cadence.KeyValuePair{Key: cadence.String("Copyright Owner 001"), Value: test.CadenceUFix64("0.2")}
	pair2 = cadence.KeyValuePair{Key: cadence.String("Copyright Owner 002"), Value: test.CadenceUFix64("0.6")}
	pair3 = cadence.KeyValuePair{Key: cadence.String("Metaya"), Value: test.CadenceUFix64("0.2")}
	recipientPairs = make([]cadence.KeyValuePair, 3)
	recipientPairs[0] = pair1
	recipientPairs[1] = pair2
	recipientPairs[2] = pair3

	tx = flow.NewTransaction().
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

	t.Run("Should be able to purchase store by cash", func(t *testing.T) {
		tx := flow.NewTransaction().
			SetScript(metaya.PurchaseStoreByCashTransaction(fungibleTokenAddress.String(), metayaUtilityCoinAddress.String(), metayaBeneficiaryCutAddress.String())).
			SetGasLimit(500).
			SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
			SetPayer(b.ServiceKey().Address).
			AddAuthorizer(metayaUtilityCoinAddress)

		_ = tx.AddArgument(cadence.NewUInt32(1))
		_ = tx.AddArgument(test.CadenceUFix64("100.0"))
		_ = tx.AddArgument(cadence.String("Commonweal 001"))

		test.SignAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, metayaUtilityCoinAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), metayaUtilityCoinSigner},
			false,
		)

		commonwealBalance := test.ExecuteScriptAndCheck(
			t, b,
			metayaUtilityCoin.GetBalanceScript(fungibleTokenAddress.String(), metayaUtilityCoinAddress.String()),
			[][]byte{jsoncdc.MustEncode(cadence.Address(commonwealAddress))},
		)
		assert.EqualValues(t, test.CadenceUFix64("0.2"), commonwealBalance)

		copyrightOwner1Balance := test.ExecuteScriptAndCheck(
			t, b,
			metayaUtilityCoin.GetBalanceScript(fungibleTokenAddress.String(), metayaUtilityCoinAddress.String()),
			[][]byte{jsoncdc.MustEncode(cadence.Address(copyrightOwner1Address))},
		)
		assert.EqualValues(t, test.CadenceUFix64("49.9"), copyrightOwner1Balance)

		copyrightOwner2Balance := test.ExecuteScriptAndCheck(
			t, b,
			metayaUtilityCoin.GetBalanceScript(fungibleTokenAddress.String(), metayaUtilityCoinAddress.String()),
			[][]byte{jsoncdc.MustEncode(cadence.Address(copyrightOwner2Address))},
		)
		assert.EqualValues(t, test.CadenceUFix64("29.94"), copyrightOwner2Balance)

		tempMetayaBalance := test.ExecuteScriptAndCheck(
			t, b,
			metayaUtilityCoin.GetBalanceScript(fungibleTokenAddress.String(), metayaUtilityCoinAddress.String()),
			[][]byte{jsoncdc.MustEncode(cadence.Address(tempMetayaAddress))},
		)
		assert.EqualValues(t, test.CadenceUFix64("19.96"), tempMetayaBalance)

		supply := test.ExecuteScriptAndCheck(
			t, b,
			metayaUtilityCoin.GetSupplyScript(fungibleTokenAddress.String(), metayaUtilityCoinAddress.String()),
			nil,
		)
		assert.EqualValues(t, test.CadenceUFix64("100.0"), supply)
	})

	t.Run("Should be able to purchase store by IUC", func(t *testing.T) {
		metayaUtilityCoin.MintTokens(
			t, b,
			fungibleTokenAddress,
			metayaUtilityCoinAddress,
			metayaUtilityCoinSigner,
			userAddress,
			"100.0",
			false,
		)

		tx := flow.NewTransaction().
			SetScript(metaya.PurchaseStoreByIUCTransaction(fungibleTokenAddress.String(), metayaUtilityCoinAddress.String(), metayaBeneficiaryCutAddress.String())).
			SetGasLimit(500).
			SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
			SetPayer(b.ServiceKey().Address).
			AddAuthorizer(userAddress)

		_ = tx.AddArgument(cadence.NewUInt32(1))
		_ = tx.AddArgument(test.CadenceUFix64("100.0"))
		_ = tx.AddArgument(cadence.String("Commonweal 001"))

		test.SignAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, userAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), userSigner},
			false,
		)

		commonwealBalance := test.ExecuteScriptAndCheck(
			t, b,
			metayaUtilityCoin.GetBalanceScript(fungibleTokenAddress.String(), metayaUtilityCoinAddress.String()),
			[][]byte{jsoncdc.MustEncode(cadence.Address(commonwealAddress))},
		)
		assert.EqualValues(t, test.CadenceUFix64("0.4"), commonwealBalance)

		copyrightOwner1Balance := test.ExecuteScriptAndCheck(
			t, b,
			metayaUtilityCoin.GetBalanceScript(fungibleTokenAddress.String(), metayaUtilityCoinAddress.String()),
			[][]byte{jsoncdc.MustEncode(cadence.Address(copyrightOwner1Address))},
		)
		assert.EqualValues(t, test.CadenceUFix64("99.8"), copyrightOwner1Balance)

		copyrightOwner2Balance := test.ExecuteScriptAndCheck(
			t, b,
			metayaUtilityCoin.GetBalanceScript(fungibleTokenAddress.String(), metayaUtilityCoinAddress.String()),
			[][]byte{jsoncdc.MustEncode(cadence.Address(copyrightOwner2Address))},
		)
		assert.EqualValues(t, test.CadenceUFix64("59.88"), copyrightOwner2Balance)

		tempMetayaBalance := test.ExecuteScriptAndCheck(
			t, b,
			metayaUtilityCoin.GetBalanceScript(fungibleTokenAddress.String(), metayaUtilityCoinAddress.String()),
			[][]byte{jsoncdc.MustEncode(cadence.Address(tempMetayaAddress))},
		)
		assert.EqualValues(t, test.CadenceUFix64("39.92"), tempMetayaBalance)

		userBalance := test.ExecuteScriptAndCheck(
			t, b,
			metayaUtilityCoin.GetBalanceScript(fungibleTokenAddress.String(), metayaUtilityCoinAddress.String()),
			[][]byte{jsoncdc.MustEncode(cadence.Address(userAddress))},
		)
		assert.EqualValues(t, test.CadenceUFix64("0.0"), userBalance)

		supply := test.ExecuteScriptAndCheck(
			t, b,
			metayaUtilityCoin.GetSupplyScript(fungibleTokenAddress.String(), metayaUtilityCoinAddress.String()),
			nil,
		)
		assert.EqualValues(t, test.CadenceUFix64("200.0"), supply)
	})

	t.Run("Should be able to purchase pack by cash", func(t *testing.T) {
		tx := flow.NewTransaction().
			SetScript(metaya.PurchasePackByCashTransaction(fungibleTokenAddress.String(), metayaUtilityCoinAddress.String(), metayaBeneficiaryCutAddress.String())).
			SetGasLimit(500).
			SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
			SetPayer(b.ServiceKey().Address).
			AddAuthorizer(metayaUtilityCoinAddress)

		_ = tx.AddArgument(cadence.NewUInt32(1))
		_ = tx.AddArgument(test.CadenceUFix64("100.0"))
		_ = tx.AddArgument(cadence.String("Commonweal 001"))

		test.SignAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, metayaUtilityCoinAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), metayaUtilityCoinSigner},
			false,
		)

		commonwealBalance := test.ExecuteScriptAndCheck(
			t, b,
			metayaUtilityCoin.GetBalanceScript(fungibleTokenAddress.String(), metayaUtilityCoinAddress.String()),
			[][]byte{jsoncdc.MustEncode(cadence.Address(commonwealAddress))},
		)
		assert.EqualValues(t, test.CadenceUFix64("0.6"), commonwealBalance)

		copyrightOwner1Balance := test.ExecuteScriptAndCheck(
			t, b,
			metayaUtilityCoin.GetBalanceScript(fungibleTokenAddress.String(), metayaUtilityCoinAddress.String()),
			[][]byte{jsoncdc.MustEncode(cadence.Address(copyrightOwner1Address))},
		)
		assert.EqualValues(t, test.CadenceUFix64("119.76"), copyrightOwner1Balance)

		copyrightOwner2Balance := test.ExecuteScriptAndCheck(
			t, b,
			metayaUtilityCoin.GetBalanceScript(fungibleTokenAddress.String(), metayaUtilityCoinAddress.String()),
			[][]byte{jsoncdc.MustEncode(cadence.Address(copyrightOwner2Address))},
		)
		assert.EqualValues(t, test.CadenceUFix64("119.76"), copyrightOwner2Balance)

		tempMetayaBalance := test.ExecuteScriptAndCheck(
			t, b,
			metayaUtilityCoin.GetBalanceScript(fungibleTokenAddress.String(), metayaUtilityCoinAddress.String()),
			[][]byte{jsoncdc.MustEncode(cadence.Address(tempMetayaAddress))},
		)
		assert.EqualValues(t, test.CadenceUFix64("59.88"), tempMetayaBalance)

		supply := test.ExecuteScriptAndCheck(
			t, b,
			metayaUtilityCoin.GetSupplyScript(fungibleTokenAddress.String(), metayaUtilityCoinAddress.String()),
			nil,
		)
		assert.EqualValues(t, test.CadenceUFix64("300.0"), supply)
	})

	t.Run("Should be able to purchase pack by IUC", func(t *testing.T) {
		metayaUtilityCoin.MintTokens(
			t, b,
			fungibleTokenAddress,
			metayaUtilityCoinAddress,
			metayaUtilityCoinSigner,
			userAddress,
			"100.0",
			false,
		)

		tx := flow.NewTransaction().
			SetScript(metaya.PurchasePackByIUCTransaction(fungibleTokenAddress.String(), metayaUtilityCoinAddress.String(), metayaBeneficiaryCutAddress.String())).
			SetGasLimit(500).
			SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
			SetPayer(b.ServiceKey().Address).
			AddAuthorizer(userAddress)

		_ = tx.AddArgument(cadence.NewUInt32(1))
		_ = tx.AddArgument(test.CadenceUFix64("100.0"))
		_ = tx.AddArgument(cadence.String("Commonweal 001"))

		test.SignAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, userAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), userSigner},
			false,
		)

		commonwealBalance := test.ExecuteScriptAndCheck(
			t, b,
			metayaUtilityCoin.GetBalanceScript(fungibleTokenAddress.String(), metayaUtilityCoinAddress.String()),
			[][]byte{jsoncdc.MustEncode(cadence.Address(commonwealAddress))},
		)
		assert.EqualValues(t, test.CadenceUFix64("0.8"), commonwealBalance)

		copyrightOwner1Balance := test.ExecuteScriptAndCheck(
			t, b,
			metayaUtilityCoin.GetBalanceScript(fungibleTokenAddress.String(), metayaUtilityCoinAddress.String()),
			[][]byte{jsoncdc.MustEncode(cadence.Address(copyrightOwner1Address))},
		)
		assert.EqualValues(t, test.CadenceUFix64("139.72"), copyrightOwner1Balance)

		copyrightOwner2Balance := test.ExecuteScriptAndCheck(
			t, b,
			metayaUtilityCoin.GetBalanceScript(fungibleTokenAddress.String(), metayaUtilityCoinAddress.String()),
			[][]byte{jsoncdc.MustEncode(cadence.Address(copyrightOwner2Address))},
		)
		assert.EqualValues(t, test.CadenceUFix64("179.64"), copyrightOwner2Balance)

		tempMetayaBalance := test.ExecuteScriptAndCheck(
			t, b,
			metayaUtilityCoin.GetBalanceScript(fungibleTokenAddress.String(), metayaUtilityCoinAddress.String()),
			[][]byte{jsoncdc.MustEncode(cadence.Address(tempMetayaAddress))},
		)
		assert.EqualValues(t, test.CadenceUFix64("79.84"), tempMetayaBalance)

		userBalance := test.ExecuteScriptAndCheck(
			t, b,
			metayaUtilityCoin.GetBalanceScript(fungibleTokenAddress.String(), metayaUtilityCoinAddress.String()),
			[][]byte{jsoncdc.MustEncode(cadence.Address(userAddress))},
		)
		assert.EqualValues(t, test.CadenceUFix64("0.0"), userBalance)

		supply := test.ExecuteScriptAndCheck(
			t, b,
			metayaUtilityCoin.GetSupplyScript(fungibleTokenAddress.String(), metayaUtilityCoinAddress.String()),
			nil,
		)
		assert.EqualValues(t, test.CadenceUFix64("400.0"), supply)
	})
}
