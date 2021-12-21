package test

import (
	"testing"

	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
	"github.com/stretchr/testify/assert"

	"github.com/0xRaito/metaya-smart-contracts/tests/go/soundlinksDID"
	"github.com/0xRaito/metaya-smart-contracts/tests/go/test"
)

var hashs = []string{"FFFFFFFFFFFFFFFFFFF1","FFFFFFFFFFFFFFFFFFF2","FFFFFFFFFFFFFFFFFFF3"}

func TestSoundlinksDIDDeployContracts(t *testing.T) {
	b := test.NewBlockchain()
	soundlinksDID.DeployContracts(t, b)
}

func TestSoundlinksDIDMintDIDs(t *testing.T) {
	b := test.NewBlockchain()

	nftAddress, soundlinksDIDAddr, soundlinksDIDSigner := soundlinksDID.DeployContracts(t, b)

	supply := test.ExecuteScriptAndCheck(
		t, b,
		soundlinksDID.GetSoundlinksDIDSupplyScript(nftAddress.String(), soundlinksDIDAddr.String()),
		nil,
	)
	assert.EqualValues(t, cadence.NewUInt64(0), supply)

	// Assert that the account collection is empty
	length := test.ExecuteScriptAndCheck(
		t, b,
		soundlinksDID.GetCollectionLengthScript(nftAddress.String(), soundlinksDIDAddr.String()),
		[][]byte{jsoncdc.MustEncode(cadence.NewAddress(soundlinksDIDAddr))},
	)
	assert.EqualValues(t, cadence.NewInt32(0), length)

	t.Run("Should be able to mint Soundlinks DIDs", func(t *testing.T) {
		soundlinksDID.MintDIDs(t, b, nftAddress, soundlinksDIDAddr, soundlinksDIDAddr, soundlinksDIDSigner, hashs)

		// Assert that the account collection is correct length
		length = test.ExecuteScriptAndCheck(
			t, b,
			soundlinksDID.GetCollectionLengthScript(nftAddress.String(), soundlinksDIDAddr.String()),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(soundlinksDIDAddr))},
		)
		assert.EqualValues(t, cadence.NewInt32(3), length)
	})
}

func TestSoundlinksDIDPurchaseDIDs(t *testing.T) {
	b := test.NewBlockchain()

	nftAddress, soundlinksDIDAddr, soundlinksDIDSigner := soundlinksDID.DeployContracts(t, b)

	userAddress, userSigner, _ := test.CreateAccount(t, b)

	test.FundAccount(t, b, userAddress, test.DefaultAccountFunding)

	// Create a new Collection
	t.Run("Should be able to create a new empty Soundlinks DID collection", func(t *testing.T) {
		soundlinksDID.SetupAccount(t, b, userAddress, userSigner, nftAddress, soundlinksDIDAddr)

		length := test.ExecuteScriptAndCheck(
			t, b,
			soundlinksDID.GetCollectionLengthScript(nftAddress.String(), soundlinksDIDAddr.String()),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(userAddress))},
		)
		assert.EqualValues(t, cadence.NewInt32(0), length)
	})

	// Purchase DIDs
	t.Run("Should be able to purchase DIDs and deposit to the purchase accounts collection", func(t *testing.T) {
		purchaseAmount := uint32(2)
		purchaseUnitPrice := "1.0"

		soundlinksDID.PurchaseDIDs(
			t, b,
			nftAddress, soundlinksDIDAddr, userAddress,
			soundlinksDIDSigner, userSigner,
			purchaseAmount, hashs, purchaseUnitPrice,
			false,
		)

		length := test.ExecuteScriptAndCheck(
			t, b,
			soundlinksDID.GetCollectionLengthScript(nftAddress.String(), soundlinksDIDAddr.String()),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(userAddress))},
		)
		assert.EqualValues(t, cadence.NewInt32(2), length)
	})
}
