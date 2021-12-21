package metayaBeneficiaryCut

import (
	"regexp"
	"testing"

	"github.com/onflow/cadence"
	emulator "github.com/onflow/flow-emulator"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	sdktemplates "github.com/onflow/flow-go-sdk/templates"
	sdktest "github.com/onflow/flow-go-sdk/test"
	"github.com/stretchr/testify/assert"

	"github.com/0xRaito/metaya-smart-contracts/tests/go/test"
	"github.com/0xRaito/metaya-smart-contracts/tests/go/metayaUtilityCoin"
)

const (
	metayaBeneficiaryCutTransactionsRootPath                  = "../../transactions/MetayaBeneficiaryCut"
	metayaBeneficiaryCutScriptsRootPath                       = "../../scripts/MetayaBeneficiaryCut"

	metayaBeneficiaryCutContractPath                          = "../../contracts/MetayaBeneficiaryCut.cdc"

	metayaBeneficiaryCutSetCopyrightOwnerPath                 = metayaBeneficiaryCutTransactionsRootPath + "/set_copyrightOwner.cdc"
	metayaBeneficiaryCutSetCommonwealPath                     = metayaBeneficiaryCutTransactionsRootPath + "/set_commonweal.cdc"
	metayaBeneficiaryCutSetMetayaCapabilityPath               = metayaBeneficiaryCutTransactionsRootPath + "/set_metaya_capability.cdc"
	metayaBeneficiaryCutSetMetayaMarketCutPercentagePath      = metayaBeneficiaryCutTransactionsRootPath + "/set_metaya_marketCutPercentage.cdc"
	metayaBeneficiaryCutSetStoreCutPercentagePath             = metayaBeneficiaryCutTransactionsRootPath + "/set_storeCutPercentage.cdc"
	metayaBeneficiaryCutSetPackCutPercentagePath              = metayaBeneficiaryCutTransactionsRootPath + "/set_packCutPercentage.cdc"
	metayaBeneficiaryCutSetMarketCutPercentagePath            = metayaBeneficiaryCutTransactionsRootPath + "/set_marketCutPercentage.cdc"
	metayaBeneficiaryCutDelCopyrightOwnerPath                 = metayaBeneficiaryCutTransactionsRootPath + "/del_copyrightOwner.cdc"
	metayaBeneficiaryCutDelCommonwealPath                     = metayaBeneficiaryCutTransactionsRootPath + "/del_commonweal.cdc"
	metayaBeneficiaryCutDelStoreCutPercentagePath             = metayaBeneficiaryCutTransactionsRootPath + "/del_storeCutPercentage.cdc"
	metayaBeneficiaryCutDelPackCutPercentagePath              = metayaBeneficiaryCutTransactionsRootPath + "/del_packCutPercentage.cdc"
	metayaBeneficiaryCutDelMarketCutPercentagePath            = metayaBeneficiaryCutTransactionsRootPath + "/del_marketCutPercentage.cdc"

	metayaBeneficiaryCutGetCopyrightOwnerNamesPath            = metayaBeneficiaryCutScriptsRootPath + "/get_copyrightOwner_names.cdc"
	metayaBeneficiaryCutGetCopyrightOwnerAmountPath           = metayaBeneficiaryCutScriptsRootPath + "/get_copyrightOwner_amount.cdc"
	metayaBeneficiaryCutGetCopyrightOwnerContainPath          = metayaBeneficiaryCutScriptsRootPath + "/get_copyrightOwner_contain.cdc"
	metayaBeneficiaryCutGetCopyrightOwnerAddressByNamePath    = metayaBeneficiaryCutScriptsRootPath + "/get_copyrightOwner_address_by_name.cdc"
	metayaBeneficiaryCutGetCommonwealNamesPath                = metayaBeneficiaryCutScriptsRootPath + "/get_commonweal_names.cdc"
	metayaBeneficiaryCutGetCommonwealCutPercentageByNamePath  = metayaBeneficiaryCutScriptsRootPath + "/get_commonwealCutPercentage_by_name.cdc"
	metayaBeneficiaryCutGetMetayaAddressPath                  = metayaBeneficiaryCutScriptsRootPath + "/get_metaya_address.cdc"
	metayaBeneficiaryCutGetMetayaMarketCutPercentagePath      = metayaBeneficiaryCutScriptsRootPath + "/get_metaya_marketCutPercentage.cdc"
	metayaBeneficiaryCutGetStoreCutPercentagesAmountPath      = metayaBeneficiaryCutScriptsRootPath + "/get_storeCutPercentages_amount.cdc"
	metayaBeneficiaryCutGetStoreCutPercentageByNamePath       = metayaBeneficiaryCutScriptsRootPath + "/get_storeCutPercentage_by_name.cdc"
	metayaBeneficiaryCutGetPackCutPercentagesAmountPath       = metayaBeneficiaryCutScriptsRootPath + "/get_packCutPercentages_amount.cdc"
	metayaBeneficiaryCutGetPackCutPercentageByNamePath        = metayaBeneficiaryCutScriptsRootPath + "/get_packCutPercentage_by_name.cdc"
	metayaBeneficiaryCutGetMarketCutPercentagesAmountPath     = metayaBeneficiaryCutScriptsRootPath + "/get_marketCutPercentages_amount.cdc"
	metayaBeneficiaryCutGetMarketCutPercentageByNamePath      = metayaBeneficiaryCutScriptsRootPath + "/get_marketCutPercentage_by_name.cdc"
)

func DeployContracts(
	t *testing.T,
	b *emulator.Blockchain,
) (flow.Address, flow.Address, crypto.Signer, flow.Address, crypto.Signer) {
	accountKeys := sdktest.AccountKeyGenerator()

	metayaBeneficiaryCutAccountKey, metayaBeneficiaryCutSigner := accountKeys.NewWithSigner()
	metayaBeneficiaryCutCode := loadMetayaBeneficiaryCut(test.FTAddress.String())

	metayaBeneficiaryCutAddress, err := b.CreateAccount(
		[]*flow.AccountKey{metayaBeneficiaryCutAccountKey},
		[]sdktemplates.Contract{
			{
				Name:   "MetayaBeneficiaryCut",
				Source: string(metayaBeneficiaryCutCode),
			},
		},
	)
	assert.NoError(t, err)

	_, err = b.CommitBlock()
	assert.NoError(t, err)

	metayaUtilityCoinAccountKey, metayaUtilityCoinSigner := accountKeys.NewWithSigner()
	metayaUtilityCoinCode := metayaUtilityCoin.LoadMetayaUtilityCoin(test.FTAddress.String())

	metayaUtilityCoinAddress, err := b.CreateAccount(
		[]*flow.AccountKey{metayaUtilityCoinAccountKey},
		[]sdktemplates.Contract{
			{
				Name:   "MetayaUtilityCoin",
				Source: string(metayaUtilityCoinCode),
			},
		},
	)
	assert.NoError(t, err)

	_, err = b.CommitBlock()
	assert.NoError(t, err)

	return test.FTAddress, metayaUtilityCoinAddress, metayaUtilityCoinSigner, metayaBeneficiaryCutAddress, metayaBeneficiaryCutSigner
}

func CreateCopyrightOwner(
	t *testing.T,
	b *emulator.Blockchain,
	copyrightOwnerName string,
	copyrightOwnerAddress flow.Address,
	fungibleTokenAddress, metayaUtilityCoinAddress, metayaBeneficiaryCutAddress flow.Address,
	metayaBeneficiaryCutSigner crypto.Signer,
) {
	tx := flow.NewTransaction().
	SetScript(SetCopyrightOwnerTransaction(fungibleTokenAddress.String(), metayaUtilityCoinAddress.String(), metayaBeneficiaryCutAddress.String())).
	SetGasLimit(100).
	SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
	SetPayer(b.ServiceKey().Address).
	AddAuthorizer(metayaBeneficiaryCutAddress)

	_ = tx.AddArgument(cadence.String(copyrightOwnerName))
	_ = tx.AddArgument(cadence.Address(copyrightOwnerAddress))

	test.SignAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, metayaBeneficiaryCutAddress},
		[]crypto.Signer{b.ServiceKey().Signer(), metayaBeneficiaryCutSigner},
		false,
	)
}

func CreateCommonweal(
	t *testing.T,
	b *emulator.Blockchain,
	commonwealName string,
	commonwealCutPercentage string,
	commonwealAddress flow.Address,
	fungibleTokenAddress, metayaUtilityCoinAddress, metayaBeneficiaryCutAddress flow.Address,
	metayaBeneficiaryCutSigner crypto.Signer,
) {
	tx := flow.NewTransaction().
	SetScript(SetCommonwealTransaction(fungibleTokenAddress.String(), metayaUtilityCoinAddress.String(), metayaBeneficiaryCutAddress.String())).
	SetGasLimit(100).
	SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
	SetPayer(b.ServiceKey().Address).
	AddAuthorizer(metayaBeneficiaryCutAddress)

	_ = tx.AddArgument(cadence.String(commonwealName))
	_ = tx.AddArgument(cadence.Address(commonwealAddress))
	_ = tx.AddArgument(test.CadenceUFix64(commonwealCutPercentage))

	test.SignAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, metayaBeneficiaryCutAddress},
		[]crypto.Signer{b.ServiceKey().Signer(), metayaBeneficiaryCutSigner},
		false,
	)
}

func replaceAddressPlaceholders(code, metayaBeneficiaryCutAddress string) []byte {
	return []byte(test.ReplaceImports(
		code,
		map[string]*regexp.Regexp{
			metayaBeneficiaryCutAddress:       test.MetayaBeneficiaryCutAddressPlaceHolder,
		},
	))
}

func loadMetayaBeneficiaryCut(fungibleTokenAddress string) []byte {
	return []byte(test.ReplaceImports(
		string(test.ReadFile(metayaBeneficiaryCutContractPath)),
		map[string]*regexp.Regexp{
			fungibleTokenAddress:              test.FungibleTokenAddressPlaceholder,
		},
	))
}

func SetCopyrightOwnerTransaction(fungibleTokenAddress, metayaUtilityCoinAddress, metayaBeneficiaryCutAddress string) []byte {
	code := string(test.ReadFile(metayaBeneficiaryCutSetCopyrightOwnerPath))

	code = test.FungibleTokenAddressPlaceholder.ReplaceAllString(code, "0x"+fungibleTokenAddress)
	code = test.MetayaUtilityCoinAddressPlaceHolder.ReplaceAllString(code, "0x"+metayaUtilityCoinAddress)

	return replaceAddressPlaceholders(
		code,
		metayaBeneficiaryCutAddress,
	)
}

func SetCommonwealTransaction(fungibleTokenAddress, metayaUtilityCoinAddress, metayaBeneficiaryCutAddress string) []byte {
	code := string(test.ReadFile(metayaBeneficiaryCutSetCommonwealPath))

	code = test.FungibleTokenAddressPlaceholder.ReplaceAllString(code, "0x"+fungibleTokenAddress)
	code = test.MetayaUtilityCoinAddressPlaceHolder.ReplaceAllString(code, "0x"+metayaUtilityCoinAddress)

	return replaceAddressPlaceholders(
		code,
		metayaBeneficiaryCutAddress,
	)
}

func SetMetayaCapabilityTransaction(fungibleTokenAddress, metayaUtilityCoinAddress, metayaBeneficiaryCutAddress string) []byte {
	code := string(test.ReadFile(metayaBeneficiaryCutSetMetayaCapabilityPath))

	code = test.FungibleTokenAddressPlaceholder.ReplaceAllString(code, "0x"+fungibleTokenAddress)
	code = test.MetayaUtilityCoinAddressPlaceHolder.ReplaceAllString(code, "0x"+metayaUtilityCoinAddress)

	return replaceAddressPlaceholders(
		code,
		metayaBeneficiaryCutAddress,
	)
}

func SetMetayaMarketCutPercentageTransaction(metayaBeneficiaryCutAddress string) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaBeneficiaryCutSetMetayaMarketCutPercentagePath)),
		metayaBeneficiaryCutAddress,
	)
}

func SetStoreCutPercentageTransaction(metayaBeneficiaryCutAddress string) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaBeneficiaryCutSetStoreCutPercentagePath)),
		metayaBeneficiaryCutAddress,
	)
}

func SetPackCutPercentageTransaction(metayaBeneficiaryCutAddress string) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaBeneficiaryCutSetPackCutPercentagePath)),
		metayaBeneficiaryCutAddress,
	)
}

func SetMarketCutPercentageTransaction(metayaBeneficiaryCutAddress string) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaBeneficiaryCutSetMarketCutPercentagePath)),
		metayaBeneficiaryCutAddress,
	)
}

func DelCopyrightOwnerTransaction(metayaBeneficiaryCutAddress string) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaBeneficiaryCutDelCopyrightOwnerPath)),
		metayaBeneficiaryCutAddress,
	)
}

func DelCommonwealTransaction(metayaBeneficiaryCutAddress string) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaBeneficiaryCutDelCommonwealPath)),
		metayaBeneficiaryCutAddress,
	)
}

func DelStoreCutPercentageTransaction(metayaBeneficiaryCutAddress string) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaBeneficiaryCutDelStoreCutPercentagePath)),
		metayaBeneficiaryCutAddress,
	)
}

func DelPackCutPercentageTransaction(metayaBeneficiaryCutAddress string) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaBeneficiaryCutDelPackCutPercentagePath)),
		metayaBeneficiaryCutAddress,
	)
}

func DelMarketCutPercentageTransaction(metayaBeneficiaryCutAddress string) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaBeneficiaryCutDelMarketCutPercentagePath)),
		metayaBeneficiaryCutAddress,
	)
}

func GetCopyrightOwnerNamesScript(metayaBeneficiaryCutAddress string) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaBeneficiaryCutGetCopyrightOwnerNamesPath)),
		metayaBeneficiaryCutAddress,
	)
}

func GetCopyrightOwnerAmountScript(metayaBeneficiaryCutAddress string) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaBeneficiaryCutGetCopyrightOwnerAmountPath)),
		metayaBeneficiaryCutAddress,
	)
}

func GetCopyrightOwnerContainScript(metayaBeneficiaryCutAddress string) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaBeneficiaryCutGetCopyrightOwnerContainPath)),
		metayaBeneficiaryCutAddress,
	)
}

func GetCopyrightOwnerAddressByNameScript(metayaBeneficiaryCutAddress string) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaBeneficiaryCutGetCopyrightOwnerAddressByNamePath)),
		metayaBeneficiaryCutAddress,
	)
}

func GetCommonwealNamesScript(metayaBeneficiaryCutAddress string) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaBeneficiaryCutGetCommonwealNamesPath)),
		metayaBeneficiaryCutAddress,
	)
}

func GetCommonwealCutPercentageByNameScript(metayaBeneficiaryCutAddress string) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaBeneficiaryCutGetCommonwealCutPercentageByNamePath)),
		metayaBeneficiaryCutAddress,
	)
}

func GetMetayaAddressScript(metayaBeneficiaryCutAddress string) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaBeneficiaryCutGetMetayaAddressPath)),
		metayaBeneficiaryCutAddress,
	)
}

func GetMetayaMarketCutPercentageScript(metayaBeneficiaryCutAddress string) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaBeneficiaryCutGetMetayaMarketCutPercentagePath)),
		metayaBeneficiaryCutAddress,
	)
}

func GetStoreCutPercentagesAmountScript(metayaBeneficiaryCutAddress string) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaBeneficiaryCutGetStoreCutPercentagesAmountPath)),
		metayaBeneficiaryCutAddress,
	)
}

func GetStoreCutPercentageByNameScript(metayaBeneficiaryCutAddress string) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaBeneficiaryCutGetStoreCutPercentageByNamePath)),
		metayaBeneficiaryCutAddress,
	)
}

func GetPackCutPercentagesAmountScript(metayaBeneficiaryCutAddress string) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaBeneficiaryCutGetPackCutPercentagesAmountPath)),
		metayaBeneficiaryCutAddress,
	)
}

func GetPackCutPercentageByNameScript(metayaBeneficiaryCutAddress string) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaBeneficiaryCutGetPackCutPercentageByNamePath)),
		metayaBeneficiaryCutAddress,
	)
}

func GetMarketCutPercentagesAmountScript(metayaBeneficiaryCutAddress string) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaBeneficiaryCutGetMarketCutPercentagesAmountPath)),
		metayaBeneficiaryCutAddress,
	)
}

func GetMarketCutPercentageByNameScript(metayaBeneficiaryCutAddress string) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaBeneficiaryCutGetMarketCutPercentageByNamePath)),
		metayaBeneficiaryCutAddress,
	)
}
