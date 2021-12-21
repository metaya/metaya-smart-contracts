package metaya

import (
	"regexp"
	"testing"

	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
	emulator "github.com/onflow/flow-emulator"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	sdktemplates "github.com/onflow/flow-go-sdk/templates"
	sdktest "github.com/onflow/flow-go-sdk/test"
	"github.com/stretchr/testify/assert"

	"github.com/0xRaito/metaya-smart-contracts/tests/go/test"
	"github.com/0xRaito/metaya-smart-contracts/tests/go/soundlinksDID"
)

const (
	metayaTransactionsRootPath          = "../../transactions/MetayaNFT"
	metayaScriptsRootPath               = "../../scripts/MetayaNFT"

	metayaContractPath                  = "../../contracts/Metaya.cdc"
	metayaShardedCollectionContractPath = "../../contracts/MetayaShardedCollection.cdc"
	metayaAdminReceiverContractPath     = "../../contracts/MetayaAdminReceiver.cdc"
	
	metayaSetupAccountPath              = metayaTransactionsRootPath + "/user/setup_account.cdc"
	metayaTransferMomentPath            = metayaTransactionsRootPath + "/user/transfer_moment.cdc"
	metayaBatchTransferMomentsPath      = metayaTransactionsRootPath + "/user/batch_transfer_moments.cdc"

	metayaCreatePlayPath                = metayaTransactionsRootPath + "/admin/create_play.cdc"
	metayaCreateSetPath                 = metayaTransactionsRootPath + "/admin/create_set.cdc"
	metayaAddPlayToSetPath              = metayaTransactionsRootPath + "/admin/add_play_to_set.cdc"
	metayaAddPlaysToSetPath             = metayaTransactionsRootPath + "/admin/add_plays_to_set.cdc"
	metayaLockSetPath                   = metayaTransactionsRootPath + "/admin/lock_set.cdc"
	metayaRetirePlayFromSetPath         = metayaTransactionsRootPath + "/admin/retire_play_from_set.cdc"
	metayaRetireAllPlaysFromSetPath     = metayaTransactionsRootPath + "/admin/retire_allPlays_from_set.cdc"
	metayaStartNewSeriesPath            = metayaTransactionsRootPath + "/admin/start_new_series.cdc"
	metayaMintMomentPath                = metayaTransactionsRootPath + "/admin/mint_moment.cdc"
	metayaBatchMintMomentsPath          = metayaTransactionsRootPath + "/admin/batch_mint_moments.cdc"
	metayaProvideMomentPath             = metayaTransactionsRootPath + "/admin/provide_moment.cdc"
	metayaPurchasePackByIUCPath         = metayaTransactionsRootPath + "/admin/purchase_pack_by_IUC.cdc"
	metayaPurchasePackByCashPath        = metayaTransactionsRootPath + "/admin/purchase_pack_by_cash.cdc"
	metayaPurchaseStoreByIUCPath        = metayaTransactionsRootPath + "/admin/purchase_store_by_IUC.cdc"
	metayaPurchaseStoreByCashPath       = metayaTransactionsRootPath + "/admin/purchase_store_by_cash.cdc"
	metayaTransferAdminPath             = metayaTransactionsRootPath + "/admin/transfer_admin.cdc"
	metayaCreatePlayStructPath          = metayaTransactionsRootPath + "/admin/create_play_struct.cdc"

	metayaSetupShardedCollectionPath    = metayaTransactionsRootPath + "/shardedCollection/setup_sharded_collection.cdc"
	metayaTransferFromShardedPath       = metayaTransactionsRootPath + "/shardedCollection/transfer_from_sharded.cdc"
	metayaBatchTransferFromShardedPath  = metayaTransactionsRootPath + "/shardedCollection/batch_transfer_from_sharded.cdc"

	metayaGetCurrentSeriesPath          = metayaScriptsRootPath + "/get_currentSeries.cdc"
	metayaGetTotalSupplyPath            = metayaScriptsRootPath + "/get_totalSupply.cdc"

	metayaGetAllPlaysPath               = metayaScriptsRootPath + "/plays/get_all_plays.cdc"
	metayaGetNextPlayIDPath             = metayaScriptsRootPath + "/plays/get_nextPlayID.cdc"
	metayaGetPlayMetadataPath           = metayaScriptsRootPath + "/plays/get_play_metadata.cdc"
	metayaGetPlayMetadataFieldPath      = metayaScriptsRootPath + "/plays/get_play_metadata_field.cdc"

	metayaGetEditionRetiredPath         = metayaScriptsRootPath + "/sets/get_edition_retired.cdc"
	metayaGetNumMomentsInEditionPath    = metayaScriptsRootPath + "/sets/get_numMoments_in_edition.cdc"
	metayaGetSetIDsByNamePath           = metayaScriptsRootPath + "/sets/get_setIDs_by_name.cdc"
	metayaGetSetSeriesPath              = metayaScriptsRootPath + "/sets/get_setSeries.cdc"
	metayaGetNextSetIDPath              = metayaScriptsRootPath + "/sets/get_nextSetID.cdc"
	metayaGetPlaysInSetPath             = metayaScriptsRootPath + "/sets/get_plays_in_set.cdc"
	metayaGetSetNamePath                = metayaScriptsRootPath + "/sets/get_setName.cdc"
	metayaGetSetLockedPath              = metayaScriptsRootPath + "/sets/get_set_locked.cdc"
	metayaGetSetDataPath                = metayaScriptsRootPath + "/sets/get_set_data.cdc"

	metayaGetCollectionIdsPath          = metayaScriptsRootPath + "/collections/get_collection_ids.cdc"
	metayaGetIdInCollectionPath         = metayaScriptsRootPath + "/collections/get_id_in_collection.cdc"
	metayaGetMetadataPath               = metayaScriptsRootPath + "/collections/get_metadata.cdc"
	metayaGetMetadataFieldPath          = metayaScriptsRootPath + "/collections/get_metadata_field.cdc"
	metayaGetMomentSeriesPath           = metayaScriptsRootPath + "/collections/get_moment_series.cdc"
	metayaGetMomentPlayIDPath           = metayaScriptsRootPath + "/collections/get_moment_playID.cdc"
	metayaGetMomentSetIDPath            = metayaScriptsRootPath + "/collections/get_moment_setID.cdc"
	metayaGetMomentSerialNumPath        = metayaScriptsRootPath + "/collections/get_moment_serialNum.cdc"
	metayaGetMomentSetNamePath          = metayaScriptsRootPath + "/collections/get_moment_setName.cdc"
	metayaGetSetplaysAreOwnedPath       = metayaScriptsRootPath + "/collections/get_setplays_are_owned.cdc"
)

func DeployContracts(
	t *testing.T,
	b *emulator.Blockchain,
) (test.Contracts){
	accountKeys := sdktest.AccountKeyGenerator()

	nftAddress, soundlinksDIDAddress, soundlinksDIDSigner := soundlinksDID.DeployContracts(t, b)

	// Should be able to deploy a contract as a new account with one key
	metayaAccountKey, metayaSigner := accountKeys.NewWithSigner()
	metayaCode := loadMetaya(
		nftAddress.String(),
		soundlinksDIDAddress.String(),
	)

	metayaAddress, err := b.CreateAccount(
		[]*flow.AccountKey{metayaAccountKey},
		[]sdktemplates.Contract{
			{
				Name:   "Metaya",
				Source: string(metayaCode),
			},
		},
	)
	assert.NoError(t, err)

	_, err = b.CommitBlock()
	assert.NoError(t, err)

	// Simplify the workflow by having the contract address also be our initial test collection
	SetupAccount(t, b, metayaAddress, metayaSigner, metayaAddress)

	//return
	return test.Contracts{
		NonFungibleTokenAddress:         nftAddress,
		SoundlinksDIDAddress:            soundlinksDIDAddress,
		SoundlinksDIDSigner:             soundlinksDIDSigner,
		MetayaAddress:                   metayaAddress,
		MetayaSigner:                    metayaSigner,
	}
}

func SetupAccount(
	t *testing.T,
	b *emulator.Blockchain,
	userAddress flow.Address,
	userSigner crypto.Signer,
	metayaAddress flow.Address,
) {
	tx := flow.NewTransaction().
		SetScript(SetupAccountTransaction(metayaAddress.String())).
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

func MintNFT(
	t *testing.T, b *emulator.Blockchain,
	contracts test.Contracts,
	setID, playID uint32,
	recipientAddress flow.Address,
	shouldFail bool,
) {
	tx := flow.NewTransaction().
		SetScript(MintMomentTransaction(contracts)).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(contracts.MetayaAddress)

	_ = tx.AddArgument(cadence.NewUInt32(setID))
	_ = tx.AddArgument(cadence.NewUInt32(playID))
	_ = tx.AddArgument(cadence.NewAddress(recipientAddress))

	test.SignAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, contracts.MetayaAddress},
		[]crypto.Signer{b.ServiceKey().Signer(), contracts.MetayaSigner},
		shouldFail,
	)
}

func MintNFTs(
	t *testing.T, b *emulator.Blockchain,
	contracts test.Contracts,
	setID, playID uint32,
	quantity uint32,
	recipientAddress flow.Address,
	shouldFail bool,
) {
	tx := flow.NewTransaction().
		SetScript(BatchMintMomentsTransaction(contracts)).
		SetGasLimit(500).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(contracts.MetayaAddress)

	_ = tx.AddArgument(cadence.NewUInt32(setID))
	_ = tx.AddArgument(cadence.NewUInt32(playID))
	_ = tx.AddArgument(cadence.NewUInt32(quantity))
	_ = tx.AddArgument(cadence.NewAddress(recipientAddress))

	test.SignAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, contracts.MetayaAddress},
		[]crypto.Signer{b.ServiceKey().Signer(), contracts.MetayaSigner},
		shouldFail,
	)
}

func TransferNFT(
	t *testing.T, b *emulator.Blockchain,
	contracts test.Contracts,
	userAddress flow.Address,
	userSigner crypto.Signer,
	nftStorefront flow.Address,
	nftID uint64,
	recipientAddress flow.Address,
	shouldFail bool,
) {
	tx := flow.NewTransaction().
		SetScript(TransferMomentTransaction(contracts, nftStorefront.String())).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(userAddress)

	_ = tx.AddArgument(cadence.NewAddress(recipientAddress))
	_ = tx.AddArgument(cadence.NewUInt64(nftID))

	test.SignAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, userAddress},
		[]crypto.Signer{b.ServiceKey().Signer(), userSigner},
		shouldFail,
	)
}

func TransferNFTs(
	t *testing.T, b *emulator.Blockchain,
	contracts test.Contracts,
	userAddress flow.Address,
	userSigner crypto.Signer,
	nftStorefront flow.Address,
	nftID []uint64,
	recipientAddress flow.Address,
	shouldFail bool,
) {
	tx := flow.NewTransaction().
		SetScript(BatchTransferMomentsTransaction(contracts, nftStorefront.String())).
		SetGasLimit(900).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(userAddress)

	_ = tx.AddArgument(cadence.NewAddress(recipientAddress))

	nftArray := make([]cadence.Value, len(nftID))
	for i := 0; i < len(nftID); i++ {
		nftArray[i] = cadence.NewUInt64(nftID[i])
	}
	_ = tx.AddArgument(cadence.NewArray(nftArray))

	test.SignAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, userAddress},
		[]crypto.Signer{b.ServiceKey().Signer(), userSigner},
		shouldFail,
	)
}

// Used to verify set metadata in tests
type SetMetadata struct {
	SetID  uint32
	Name   string
	Series uint32
	Plays  []uint32
	//Retired {UInt32: Bool},
	Locked bool
	//NumberMintedPerPlay {UInt32: UInt32},
}

// Verifies that the epoch metadata is equal to the provided expected values
func VerifyQuerySetMetadata(
	t *testing.T,
	b *emulator.Blockchain,
	contracts test.Contracts,
	expectedMetadata SetMetadata) {

	result := test.ExecuteScriptAndCheck(
		t, b,
		GetSetDataScript(contracts),
		[][]byte{jsoncdc.MustEncode(cadence.UInt32(expectedMetadata.SetID))},
	)
	metadataFields := result.(cadence.Struct).Fields

	setID := metadataFields[0]
	assert.EqualValues(t, cadence.NewUInt32(expectedMetadata.SetID), setID)

	name := metadataFields[1]
	assert.EqualValues(t, cadence.String(expectedMetadata.Name), name)

	series := metadataFields[2]
	assert.EqualValues(t, cadence.NewUInt32(expectedMetadata.Series), series)

	if len(expectedMetadata.Plays) != 0 {
		plays := metadataFields[3].(cadence.Array).Values

		for i, play := range plays {
			expectedPlayID := cadence.NewUInt32(expectedMetadata.Plays[i])
			assert.EqualValues(t, expectedPlayID, play)
		}
	}

	locked := metadataFields[5]
	assert.EqualValues(t, cadence.NewBool(expectedMetadata.Locked), locked)
}

func replaceAddressPlaceholders(code string, contracts test.Contracts) []byte {
	return []byte(test.ReplaceImports(
		code,
		map[string]*regexp.Regexp{
			contracts.NonFungibleTokenAddress.String():          test.NonFungibleTokenAddressPlaceholder,
			contracts.SoundlinksDIDAddress.String():             test.SoundlinksDIDAddressPlaceHolder,
			contracts.MetayaAddress.String():                    test.MetayaAddressPlaceHolder,
		},
	))
}

func loadMetaya(nftAddress, soundlinksDIDAddress string) []byte {
	return []byte(test.ReplaceImports(
		string(test.ReadFile(metayaContractPath)),
		map[string]*regexp.Regexp{
			nftAddress:           test.NonFungibleTokenAddressPlaceholder,
			soundlinksDIDAddress: test.SoundlinksDIDAddressPlaceHolder,
		},
	))
}

func LoadMetayaShardedCollection(nftAddress, metayaAddress string) []byte {
	return []byte(test.ReplaceImports(
		string(test.ReadFile(metayaShardedCollectionContractPath)),
		map[string]*regexp.Regexp{
			nftAddress:           test.NonFungibleTokenAddressPlaceholder,
			metayaAddress:        test.MetayaAddressPlaceHolder,
		},
	))
}

func LoadMetayaAdminReceiver(metayaAddress, metayaShardedCollectionAddress string) []byte {
	return []byte(test.ReplaceImports(
		string(test.ReadFile(metayaAdminReceiverContractPath)),
		map[string]*regexp.Regexp{
			metayaAddress:                   test.MetayaAddressPlaceHolder,
			metayaShardedCollectionAddress:  test.MetayaShardedCollectionAddressPlaceHolder,
		},
	))
}

func SetupAccountTransaction(metayaAddress string) []byte {
	return []byte(test.ReplaceImports(
		string(test.ReadFile(metayaSetupAccountPath)),
		map[string]*regexp.Regexp{
			metayaAddress: test.MetayaAddressPlaceHolder,
		},
	))
}

func TransferMomentTransaction(contracts test.Contracts, nftStorefront string) []byte {
	code := string(test.ReadFile(metayaTransferMomentPath))
	code = test.NFTStorefrontAddressPlaceholder.ReplaceAllString(code, "0x"+nftStorefront)

	return replaceAddressPlaceholders(
		code,
		contracts,
	)
}

func BatchTransferMomentsTransaction(contracts test.Contracts, nftStorefront string) []byte {
	code := string(test.ReadFile(metayaBatchTransferMomentsPath))
	code = test.NFTStorefrontAddressPlaceholder.ReplaceAllString(code, "0x"+nftStorefront)

	return replaceAddressPlaceholders(
		code,
		contracts,
	)
}

func CreatePlayTransaction(contracts test.Contracts) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaCreatePlayPath)),
		contracts,
	)
}

func CreateSetTransaction(contracts test.Contracts) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaCreateSetPath)),
		contracts,
	)
}

func AddPlayToSetTransaction(contracts test.Contracts) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaAddPlayToSetPath)),
		contracts,
	)
}

func AddPlaysToSetTransaction(contracts test.Contracts) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaAddPlaysToSetPath)),
		contracts,
	)
}

func LockSetTransaction(contracts test.Contracts) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaLockSetPath)),
		contracts,
	)
}

func RetirePlayFromSetTransaction(contracts test.Contracts) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaRetirePlayFromSetPath)),
		contracts,
	)
}

func RetireAllPlaysFromSetTransaction(contracts test.Contracts) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaRetireAllPlaysFromSetPath)),
		contracts,
	)
}

func StartNewSeriesTransaction(contracts test.Contracts) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaStartNewSeriesPath)),
		contracts,
	)
}

func MintMomentTransaction(contracts test.Contracts) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaMintMomentPath)),
		contracts,
	)
}

func BatchMintMomentsTransaction(contracts test.Contracts) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaBatchMintMomentsPath)),
		contracts,
	)
}

func ProvideMomentTransaction(contracts test.Contracts, metayaShardedCollectionAddress, flowStorageFeesAddress string) []byte {
	code := string(test.ReadFile(metayaProvideMomentPath))

	code = test.MetayaShardedCollectionAddressPlaceHolder.ReplaceAllString(code, "0x"+metayaShardedCollectionAddress)
	code = test.FlowStorageFeesAddressPlaceHolder.ReplaceAllString(code, "0x"+flowStorageFeesAddress)
	code = test.FungibleTokenAddressPlaceholder.ReplaceAllString(code, "0x"+test.FTAddress.String())
	code = test.FlowTokenAddressPlaceHolder.ReplaceAllString(code, "0x"+test.FlowTokenAddress.String())

	return replaceAddressPlaceholders(
		code,
		contracts,
	)
}

func PurchasePackByIUCTransaction(ftAddress, metayaUtilityCoinAddress, metayaBeneficiaryCutAddress string) []byte {
	return []byte(test.ReplaceImports(
		string(test.ReadFile(metayaPurchasePackByIUCPath)),
		map[string]*regexp.Regexp{
			ftAddress:                        test.FungibleTokenAddressPlaceholder,
			metayaUtilityCoinAddress:         test.MetayaUtilityCoinAddressPlaceHolder,
			metayaBeneficiaryCutAddress:      test.MetayaBeneficiaryCutAddressPlaceHolder,
		},
	))
}

func PurchasePackByCashTransaction(ftAddress, metayaUtilityCoinAddress, metayaBeneficiaryCutAddress string) []byte {
	return []byte(test.ReplaceImports(
		string(test.ReadFile(metayaPurchasePackByCashPath)),
		map[string]*regexp.Regexp{
			ftAddress:                        test.FungibleTokenAddressPlaceholder,
			metayaUtilityCoinAddress:         test.MetayaUtilityCoinAddressPlaceHolder,
			metayaBeneficiaryCutAddress:      test.MetayaBeneficiaryCutAddressPlaceHolder,
		},
	))
}

func PurchaseStoreByIUCTransaction(ftAddress, metayaUtilityCoinAddress, metayaBeneficiaryCutAddress string) []byte {
	return []byte(test.ReplaceImports(
		string(test.ReadFile(metayaPurchaseStoreByIUCPath)),
		map[string]*regexp.Regexp{
			ftAddress:                        test.FungibleTokenAddressPlaceholder,
			metayaUtilityCoinAddress:         test.MetayaUtilityCoinAddressPlaceHolder,
			metayaBeneficiaryCutAddress:      test.MetayaBeneficiaryCutAddressPlaceHolder,
		},
	))
}

func PurchaseStoreByCashTransaction(ftAddress, metayaUtilityCoinAddress, metayaBeneficiaryCutAddress string) []byte {
	return []byte(test.ReplaceImports(
		string(test.ReadFile(metayaPurchaseStoreByCashPath)),
		map[string]*regexp.Regexp{
			ftAddress:                        test.FungibleTokenAddressPlaceholder,
			metayaUtilityCoinAddress:         test.MetayaUtilityCoinAddressPlaceHolder,
			metayaBeneficiaryCutAddress:      test.MetayaBeneficiaryCutAddressPlaceHolder,
		},
	))
}

func TransferAdminTransaction(contracts test.Contracts, metayaAdminReceiverAdress string) []byte {
	code := string(test.ReadFile(metayaTransferAdminPath))

	code = test.MetayaAdminReceiverAddressPlaceHolder.ReplaceAllString(code, "0x"+metayaAdminReceiverAdress)
	
	return replaceAddressPlaceholders(
		code,
		contracts,
	)
}

func CreatePlayStructTransaction(contracts test.Contracts) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaCreatePlayStructPath)),
		contracts,
	)
}

func SetupShardedCollectionTransaction(contracts test.Contracts, metayaShardedCollectionAddress string) []byte {
	code := string(test.ReadFile(metayaSetupShardedCollectionPath))

	code = test.MetayaShardedCollectionAddressPlaceHolder.ReplaceAllString(code, "0x"+metayaShardedCollectionAddress)

	return replaceAddressPlaceholders(
		code,
		contracts,
	)
}

func TransferFromShardedTransaction(contracts test.Contracts, metayaShardedCollectionAddress string) []byte {
	code := string(test.ReadFile(metayaTransferFromShardedPath))
	code = test.MetayaShardedCollectionAddressPlaceHolder.ReplaceAllString(code, "0x"+metayaShardedCollectionAddress)

	return replaceAddressPlaceholders(
		code,
		contracts,
	)
}

func BatchTransferFromShardedTransaction(contracts test.Contracts, metayaShardedCollectionAddress string) []byte {
	code := string(test.ReadFile(metayaBatchTransferFromShardedPath))
	code = test.MetayaShardedCollectionAddressPlaceHolder.ReplaceAllString(code, "0x"+metayaShardedCollectionAddress)

	return replaceAddressPlaceholders(
		code,
		contracts,
	)
}

func GetCurrentSeriesScript(contracts test.Contracts) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaGetCurrentSeriesPath)),
		contracts,
	)
}

func GetTotalSupplyScript(contracts test.Contracts) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaGetTotalSupplyPath)),
		contracts,
	)
}

func GetAllPlaysScript(contracts test.Contracts) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaGetAllPlaysPath)),
		contracts,
	)
}

func GetNextPlayIDScript(contracts test.Contracts) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaGetNextPlayIDPath)),
		contracts,
	)
}

func GetPlayMetadataScript(contracts test.Contracts) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaGetPlayMetadataPath)),
		contracts,
	)
}

func GetPlayMetadataFieldScript(contracts test.Contracts) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaGetPlayMetadataFieldPath)),
		contracts,
	)
}

func GetEditionRetiredScript(contracts test.Contracts) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaGetEditionRetiredPath)),
		contracts,
	)
}

func GetNumMomentsInEditionScript(contracts test.Contracts) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaGetNumMomentsInEditionPath)),
		contracts,
	)
}

func GetSetIDsByNameScript(contracts test.Contracts) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaGetSetIDsByNamePath)),
		contracts,
	)
}

func GetSetSeriesScript(contracts test.Contracts) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaGetSetSeriesPath)),
		contracts,
	)
}

func GetNextSetIDScript(contracts test.Contracts) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaGetNextSetIDPath)),
		contracts,
	)
}

func GetPlaysInSetScript(contracts test.Contracts) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaGetPlaysInSetPath)),
		contracts,
	)
}

func GetSetNameScript(contracts test.Contracts) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaGetSetNamePath)),
		contracts,
	)
}

func GetSetLockedScript(contracts test.Contracts) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaGetSetLockedPath)),
		contracts,
	)
}

func GetSetDataScript(contracts test.Contracts) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaGetSetDataPath)),
		contracts,
	)
}

func GetCollectionIdsScript(contracts test.Contracts) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaGetCollectionIdsPath)),
		contracts,
	)
}

func GetIdInCollectionScript(contracts test.Contracts) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaGetIdInCollectionPath)),
		contracts,
	)
}

func GetMetadataScript(contracts test.Contracts) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaGetMetadataPath)),
		contracts,
	)
}

func GetMetadataFieldScript(contracts test.Contracts) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaGetMetadataFieldPath)),
		contracts,
	)
}

func GetMomentSeriesScript(contracts test.Contracts) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaGetMomentSeriesPath)),
		contracts,
	)
}

func GetMomentPlayIDScript(contracts test.Contracts) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaGetMomentPlayIDPath)),
		contracts,
	)
}

func GetMomentSetIDScript(contracts test.Contracts) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaGetMomentSetIDPath)),
		contracts,
	)
}

func GetMomentSerialNumScript(contracts test.Contracts) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaGetMomentSerialNumPath)),
		contracts,
	)
}

func GetMomentSetNameScript(contracts test.Contracts) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaGetMomentSetNamePath)),
		contracts,
	)
}

func GetSetplaysAreOwnedScript(contracts test.Contracts) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaGetSetplaysAreOwnedPath)),
		contracts,
	)
}
