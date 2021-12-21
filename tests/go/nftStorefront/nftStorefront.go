package nftStorefront

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/onflow/cadence"
	emulator "github.com/onflow/flow-emulator"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	sdktemplates "github.com/onflow/flow-go-sdk/templates"
	sdktest "github.com/onflow/flow-go-sdk/test"
	"github.com/stretchr/testify/assert"
	flowStorageFees "github.com/onflow/flow-core-contracts/lib/go/contracts"

	"github.com/0xRaito/metaya-smart-contracts/tests/go/test"
	"github.com/0xRaito/metaya-smart-contracts/tests/go/metayaUtilityCoin"
	"github.com/0xRaito/metaya-smart-contracts/tests/go/metayaNFT"
	"github.com/0xRaito/metaya-smart-contracts/tests/go/metayaBeneficiaryCut"
)

const (
	nftStorefrontContractPath          = "../../contracts/NFTStorefront.cdc"

	nftStorefrontTransactionRootPath   = "../../transactions/NFTStorefront"
	nftStorefrontScriptRootPath        = "../../scripts/NFTStorefront"

	nftStorefrontSetupAccountPath      = nftStorefrontTransactionRootPath + "/setup_account.cdc"
	nftStorefrontSellItemByIUCPath     = nftStorefrontTransactionRootPath + "/sell_item_by_IUC.cdc"
	nftStorefrontBuyItemByIUCPath      = nftStorefrontTransactionRootPath + "/buy_item_by_IUC.cdc"
	nftStorefrontBuyItemByCashPath     = nftStorefrontTransactionRootPath + "/buy_item_by_cash.cdc"
	nftStorefrontRemoveItemPath        = nftStorefrontTransactionRootPath + "/remove_item.cdc"
	nftStorefrontCleanupItemPath       = nftStorefrontTransactionRootPath + "/cleanup_item.cdc"

	nftStorefrontGetListingIdsPath     = nftStorefrontScriptRootPath + "/get_listing_ids.cdc"
	nftStorefrontGetListingDetailsPath = nftStorefrontScriptRootPath + "/get_listing_details.cdc"
)

func DeployContracts(
	t *testing.T,
	b *emulator.Blockchain,
) test.Contracts {
	accountKeys := sdktest.AccountKeyGenerator()

	ftAddress, metayaUtilityCoinAddress, metayaUtilityCoinSigner, metayaBeneficiaryCutAddress, metayaBeneficiaryCutSigner := metayaBeneficiaryCut.DeployContracts(t, b)

	contracts := metaya.DeployContracts(t, b)

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

	// Should be able to deploy a contract as a new account with one key
	nftStorefrontAccountKey, nftStorefrontSigner := accountKeys.NewWithSigner()
	nftStorefrontCode := LoadNFTStorefront(
		ftAddress,
		contracts.NonFungibleTokenAddress,
	)

	nftStorefrontAddress, err := b.CreateAccount(
		[]*flow.AccountKey{nftStorefrontAccountKey},
		[]sdktemplates.Contract{
			{
				Name:   "NFTStorefront",
				Source: string(nftStorefrontCode),
			},
		},
	)
	assert.NoError(t, err)

	_, err = b.CommitBlock()
	assert.NoError(t, err)

	// Simplify the workflow by having the contract address also be our initial test collection
	SetupAccount(t, b, nftStorefrontAddress, nftStorefrontSigner, nftStorefrontAddress)

	return test.Contracts{
		FungibleTokenAddress:            ftAddress,
		NonFungibleTokenAddress:         contracts.NonFungibleTokenAddress,
		MetayaUtilityCoinAddress:        metayaUtilityCoinAddress,
		MetayaUtilityCoinSigner:         metayaUtilityCoinSigner,
		SoundlinksDIDAddress:            contracts.SoundlinksDIDAddress,
		SoundlinksDIDSigner:             contracts.SoundlinksDIDSigner,
		MetayaAddress:                   contracts.MetayaAddress,
		MetayaSigner:                    contracts.MetayaSigner,
		MetayaBeneficiaryCutAddress:     metayaBeneficiaryCutAddress,
		MetayaBeneficiaryCutSigner:      metayaBeneficiaryCutSigner,
		NFTStorefrontAddress:            nftStorefrontAddress,
		NFTStorefrontSigner:             nftStorefrontSigner,
		FlowTokenAddress:                test.FlowTokenAddress,
		FlowStorageFeesAddress:          flowStorageFeesAddress,
	}
}

func SetupAccount(
	t *testing.T,
	b *emulator.Blockchain,
	userAddress flow.Address,
	userSigner crypto.Signer,
	nftStorefrontAddress flow.Address,
) {
	tx := flow.NewTransaction().
		SetScript(NFTStorefrontSetupAccountTransaction(nftStorefrontAddress.String())).
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

// Create a new account with the MetayaUtilityCoin and Metaya resources set up BUT no NFTStorefront resource.
func CreatePurchaserAccount(
	t *testing.T,
	b *emulator.Blockchain,
	contracts test.Contracts,
) (flow.Address, crypto.Signer) {
	userAddress, userSigner, _ := test.CreateAccount(t, b)
	metayaUtilityCoin.SetupAccount(t, b, userAddress, userSigner, contracts.FungibleTokenAddress, contracts.MetayaUtilityCoinAddress)
	metaya.SetupAccount(t, b, userAddress, userSigner, contracts.MetayaAddress)
	return userAddress, userSigner
}

func CreateAccount(
	t *testing.T,
	b *emulator.Blockchain,
	contracts test.Contracts,
) (flow.Address, crypto.Signer) {
	userAddress, userSigner := CreatePurchaserAccount(t, b, contracts)
	SetupAccount(t, b, userAddress, userSigner, contracts.NFTStorefrontAddress)
	return userAddress, userSigner
}

func ListItem(
	t *testing.T,
	b *emulator.Blockchain,
	contracts test.Contracts,
	userAddress flow.Address,
	userSigner crypto.Signer,
	tokenID uint64,
	price string,
	shouldFail bool,
) (listingResourceID uint64) {
	tx := flow.NewTransaction().
		SetScript(NFTStorefrontSellItemByIUCTransaction(contracts)).
		SetGasLimit(500).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(userAddress)

	_ = tx.AddArgument(cadence.NewUInt64(tokenID))
	_ = tx.AddArgument(test.CadenceUFix64(price))

	result := test.SignAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, userAddress},
		[]crypto.Signer{b.ServiceKey().Signer(), userSigner},
		shouldFail,
	)

	listingAvailableEventType := fmt.Sprintf(
		"A.%s.NFTStorefront.ListingAvailable",
		contracts.NFTStorefrontAddress,
	)

	for _, event := range result.Events {
		if event.Type == listingAvailableEventType {
			return event.Value.Fields[1].ToGoValue().(uint64)
		}
	}

	return 0
}

func PurchaseItemByIUC(
	t *testing.T,
	b *emulator.Blockchain,
	contracts test.Contracts,
	userAddress flow.Address,
	userSigner crypto.Signer,
	sellerAddress flow.Address,
	listingResourceID uint64,
	shouldFail bool,
) {
	tx := flow.NewTransaction().
		SetScript(NFTStorefrontBuyItemByIUCTransaction(contracts, contracts.FlowStorageFeesAddress.String())).
		SetGasLimit(500).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(userAddress).
		AddAuthorizer(contracts.MetayaAddress)

	_ = tx.AddArgument(cadence.NewUInt64(listingResourceID))
	_ = tx.AddArgument(cadence.NewAddress(sellerAddress))

	test.SignAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, userAddress, contracts.MetayaAddress},
		[]crypto.Signer{b.ServiceKey().Signer(), userSigner, contracts.MetayaSigner},
		shouldFail,
	)
}

func PurchaseItemByCash(
	t *testing.T,
	b *emulator.Blockchain,
	contracts test.Contracts,
	sellerAddress flow.Address,
	buyerAddress flow.Address,
	listingResourceID uint64,
	shouldFail bool,
) {
	tx := flow.NewTransaction().
		SetScript(NFTStorefrontBuyItemByCashTransaction(contracts, contracts.FlowStorageFeesAddress.String())).
		SetGasLimit(500).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(contracts.MetayaUtilityCoinAddress)

	_ = tx.AddArgument(cadence.NewUInt64(listingResourceID))
	_ = tx.AddArgument(cadence.NewAddress(sellerAddress))
	_ = tx.AddArgument(cadence.NewAddress(buyerAddress))

	test.SignAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, contracts.MetayaUtilityCoinAddress},
		[]crypto.Signer{b.ServiceKey().Signer(), contracts.MetayaUtilityCoinSigner},
		shouldFail,
	)
}

func RemoveItem(
	t *testing.T,
	b *emulator.Blockchain,
	contracts test.Contracts,
	userAddress flow.Address,
	userSigner crypto.Signer,
	tokenID uint64,
	shouldFail bool,
) {
	tx := flow.NewTransaction().
		SetScript(NFTStorefrontRemoveItemTransaction(contracts)).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(userAddress)

	_ = tx.AddArgument(cadence.NewUInt64(tokenID))

	test.SignAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, userAddress},
		[]crypto.Signer{b.ServiceKey().Signer(), userSigner},
		shouldFail,
	)
}

func replaceAddressPlaceholders(code string, contracts test.Contracts) []byte {
	return []byte(test.ReplaceImports(
		code,
		map[string]*regexp.Regexp{
			contracts.FungibleTokenAddress.String():            test.FungibleTokenAddressPlaceholder,
			contracts.NonFungibleTokenAddress.String():         test.NonFungibleTokenAddressPlaceholder,
			contracts.MetayaUtilityCoinAddress.String():        test.MetayaUtilityCoinAddressPlaceHolder,
			contracts.MetayaAddress.String():                   test.MetayaAddressPlaceHolder,
			contracts.MetayaBeneficiaryCutAddress.String():     test.MetayaBeneficiaryCutAddressPlaceHolder,
			contracts.NFTStorefrontAddress.String():            test.NFTStorefrontAddressPlaceholder,
		},
	))
}

func LoadNFTStorefront(
	ftAddress flow.Address,
	nftAddress flow.Address,
) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(nftStorefrontContractPath)),
		test.Contracts{
			FungibleTokenAddress:    ftAddress,
			NonFungibleTokenAddress: nftAddress,
		},
	)
}

func NFTStorefrontSetupAccountTransaction(nftStorefrontAddr string) []byte {
	return []byte(test.ReplaceImports(
		string(test.ReadFile(nftStorefrontSetupAccountPath)),
		map[string]*regexp.Regexp{
			nftStorefrontAddr: test.NFTStorefrontAddressPlaceholder,
		},
	))
}

func NFTStorefrontSellItemByIUCTransaction(contracts test.Contracts) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(nftStorefrontSellItemByIUCPath)),
		contracts,
	)
}

func NFTStorefrontBuyItemByIUCTransaction(contracts test.Contracts, flowStorageFeesAddress string) []byte {
	code := string(test.ReadFile(nftStorefrontBuyItemByIUCPath))

	code = test.FlowTokenAddressPlaceHolder.ReplaceAllString(code, "0x"+test.FlowTokenAddress.String())
	code = test.FlowStorageFeesAddressPlaceHolder.ReplaceAllString(code, "0x"+flowStorageFeesAddress)
	
	return replaceAddressPlaceholders(
		code,
		contracts,
	)
}

func NFTStorefrontBuyItemByCashTransaction(contracts test.Contracts, flowStorageFeesAddress string) []byte {
	code := string(test.ReadFile(nftStorefrontBuyItemByCashPath))

	code = test.FlowTokenAddressPlaceHolder.ReplaceAllString(code, "0x"+test.FlowTokenAddress.String())
	code = test.FlowStorageFeesAddressPlaceHolder.ReplaceAllString(code, "0x"+flowStorageFeesAddress)
	
	return replaceAddressPlaceholders(
		code,
		contracts,
	)
}

func NFTStorefrontRemoveItemTransaction(contracts test.Contracts) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(nftStorefrontRemoveItemPath)),
		contracts,
	)
}

func NFTStorefrontCleanupItemTransaction(contracts test.Contracts) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(nftStorefrontCleanupItemPath)),
		contracts,
	)
}

func NFTStorefrontGetListingIdsScript(contracts test.Contracts) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(nftStorefrontGetListingIdsPath)),
		contracts,
	)
}

func NFTStorefrontGetListingDetailsScript(contracts test.Contracts) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(nftStorefrontGetListingDetailsPath)),
		contracts,
	)
}