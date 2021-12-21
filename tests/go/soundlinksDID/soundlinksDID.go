package soundlinksDID

import (
	"regexp"
	"testing"

	"github.com/onflow/cadence"
	emulator "github.com/onflow/flow-emulator"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	sdktemplates "github.com/onflow/flow-go-sdk/templates"
	sdktest "github.com/onflow/flow-go-sdk/test"
	nftcontracts "github.com/onflow/flow-nft/lib/go/contracts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/0xRaito/metaya-smart-contracts/tests/go/test"
)

const (
	soundlinksDIDTransactionsRootPath    = "../../transactions/SoundlinksDID"
	soundlinksDIDScriptsRootPath         = "../../scripts/SoundlinksDID"

	soundlinksDIDContractPath            = "../../contracts/SoundlinksDID.cdc"

	soundlinksDIDSetupAccountPath        = soundlinksDIDTransactionsRootPath + "/setup_account.cdc"
	soundlinksDIDMintDIDsPath            = soundlinksDIDTransactionsRootPath + "/mint_DIDs.cdc"
	soundlinksDIDPurchaseDIDsPath        = soundlinksDIDTransactionsRootPath + "/purchase_DIDs.cdc"
	
	soundlinksDIDGetSupplyPath           = soundlinksDIDScriptsRootPath + "/get_supply.cdc"
	soundlinksDIDGetAmountPath           = soundlinksDIDScriptsRootPath + "/get_amount.cdc"
)

func DeployContracts(
	t *testing.T,
	b *emulator.Blockchain,
) (flow.Address, flow.Address, crypto.Signer) {
	accountKeys := sdktest.AccountKeyGenerator()

	// Should be able to deploy a contract as a new account with no keys
	nftCode := nftcontracts.NonFungibleToken()
	nftAddress, err := b.CreateAccount(
		nil,
		[]sdktemplates.Contract{
			{
				Name:   "NonFungibleToken",
				Source: string(nftCode),
			},
		},
	)
	require.NoError(t, err)

	_, err = b.CommitBlock()
	assert.NoError(t, err)

	// Should be able to deploy a contract as a new account with one key
	soundlinksDIDAccountKey, soundlinksDIDSigner := accountKeys.NewWithSigner()
	soundlinksDIDCode := loadSoundlinksDID(nftAddress.String())
	soundlinksDIDAddr, err := b.CreateAccount(
		[]*flow.AccountKey{soundlinksDIDAccountKey},
		[]sdktemplates.Contract{
			{
				Name:   "SoundlinksDID",
				Source: string(soundlinksDIDCode),
			},
		},
	)
	assert.NoError(t, err)

	_, err = b.CommitBlock()
	assert.NoError(t, err)

	// Simplify the workflow by having the contract address also be our initial test collection
	SetupAccount(t, b, soundlinksDIDAddr, soundlinksDIDSigner, nftAddress, soundlinksDIDAddr)

	return nftAddress, soundlinksDIDAddr, soundlinksDIDSigner
}

func SetupAccount(
	t *testing.T,
	b *emulator.Blockchain,
	userAddress flow.Address,
	userSigner crypto.Signer,
	nftAddress flow.Address,
	soundlinksDIDAddress flow.Address,
) {
	tx := flow.NewTransaction().
		SetScript(SetupAccountTransaction(nftAddress.String(), soundlinksDIDAddress.String())).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(userAddress)

	test.SignAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, userAddress},
		[]crypto.Signer{b.ServiceKey().Signer(), userSigner},
		false,
	)
}

func MintDIDs(
	t *testing.T, b *emulator.Blockchain,
	nftAddress, soundlinksDIDAddr, recipientAddress flow.Address,
	soundlinksDIDSigner crypto.Signer, hashs []string,
) {
	tx := flow.NewTransaction().
		SetScript(MintDIDsTransaction(nftAddress.String(), soundlinksDIDAddr.String())).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(soundlinksDIDAddr)

	_ = tx.AddArgument(cadence.NewAddress(recipientAddress))

	hashsArray := make([]cadence.Value, len(hashs))
	for i := 0; i < len(hashs); i++ {
		hashsArray[i] = cadence.String(hashs[i])
	}

	_ = tx.AddArgument(cadence.NewArray(hashsArray))
	
	test.SignAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, soundlinksDIDAddr},
		[]crypto.Signer{b.ServiceKey().Signer(), soundlinksDIDSigner},
		false,
	)
}

func PurchaseDIDs(
	t *testing.T, b *emulator.Blockchain,
	nftAddress, soundlinksDIDAddr, recipientAddr flow.Address,
	soundlinksDIDSigner, recipientSigner crypto.Signer,
	purchaseAmount uint32, hashs []string, purchaseUnitPrice string, shouldFail bool,
) {

	tx := flow.NewTransaction().
		SetScript(PurchaseDIDsTransaction(nftAddress.String(), soundlinksDIDAddr.String())).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(soundlinksDIDAddr).
		AddAuthorizer(recipientAddr)

	_ = tx.AddArgument(cadence.NewUInt32(purchaseAmount))
	_ = tx.AddArgument(cadence.NewArray([]cadence.Value{cadence.String(hashs[0]),cadence.String(hashs[1])}))
	_ = tx.AddArgument(test.CadenceUFix64(purchaseUnitPrice))

	test.SignAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, soundlinksDIDAddr, recipientAddr},
		[]crypto.Signer{b.ServiceKey().Signer(), soundlinksDIDSigner, recipientSigner},
		shouldFail,
	)
}

func replaceAddressPlaceholders(code, nftAddress, soundlinksDIDAddress string) []byte {
	return []byte(test.ReplaceImports(
		code,
		map[string]*regexp.Regexp{
			nftAddress:           test.NonFungibleTokenAddressPlaceholder,
			soundlinksDIDAddress: test.SoundlinksDIDAddressPlaceHolder,
		},
	))
}

func loadSoundlinksDID(nftAddress string) []byte {
	return []byte(test.ReplaceImports(
		string(test.ReadFile(soundlinksDIDContractPath)),
		map[string]*regexp.Regexp{
			nftAddress: test.NonFungibleTokenAddressPlaceholder,
		},
	))
}

func SetupAccountTransaction(nftAddress, soundlinksDIDAddress string) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(soundlinksDIDSetupAccountPath)),
		nftAddress,
		soundlinksDIDAddress,
	)
}

func MintDIDsTransaction(nftAddress, soundlinksDIDAddress string) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(soundlinksDIDMintDIDsPath)),
		nftAddress,
		soundlinksDIDAddress,
	)
}

func PurchaseDIDsTransaction(nftAddress, soundlinksDIDAddress string) []byte {
	code := string(test.ReadFile(soundlinksDIDPurchaseDIDsPath))

	code = test.FungibleTokenAddressPlaceholder.ReplaceAllString(code, "0x"+test.FTAddress.String())
	code = test.FlowTokenAddressPlaceHolder.ReplaceAllString(code, "0x"+test.FlowTokenAddress.String())

	return replaceAddressPlaceholders(
		code,
		nftAddress,
		soundlinksDIDAddress,
	)
}

func GetSoundlinksDIDSupplyScript(nftAddress, soundlinksDIDAddress string) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(soundlinksDIDGetSupplyPath)),
		nftAddress,
		soundlinksDIDAddress,
	)
}

func GetCollectionLengthScript(nftAddress, soundlinksDIDAddress string) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(soundlinksDIDGetAmountPath)),
		nftAddress,
		soundlinksDIDAddress,
	)
}
