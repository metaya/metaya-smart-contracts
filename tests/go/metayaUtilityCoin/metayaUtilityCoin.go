package metayaUtilityCoin

import (
	"regexp"
	"testing"

	"github.com/onflow/cadence"
	emulator "github.com/onflow/flow-emulator"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	sdktemplates "github.com/onflow/flow-go-sdk/templates"
	sdktest "github.com/onflow/flow-go-sdk/test"
	ftcontracts "github.com/onflow/flow-ft/lib/go/contracts"
	"github.com/stretchr/testify/assert"

	"github.com/0xRaito/metaya-smart-contracts/tests/go/test"
)

const (
	metayaUtilityCoinTransactionsRootPath       = "../../transactions/MetayaUtilityCoin"
	metayaUtilityCoinScriptsRootPath            = "../../scripts/MetayaUtilityCoin"

	metayaUtilityCoinContractPath               = "../../contracts/MetayaUtilityCoin.cdc"

	metayaUtilityCoinSetupAccountPath           = metayaUtilityCoinTransactionsRootPath + "/setup_account.cdc"
	metayaUtilityCoinTransferTokensPath         = metayaUtilityCoinTransactionsRootPath + "/transfer_tokens.cdc"
	metayaUtilityCoinTransferManyAccountsPath   = metayaUtilityCoinTransactionsRootPath + "/transfer_many_accounts.cdc"
	metayaUtilityCoinMintTokensPath             = metayaUtilityCoinTransactionsRootPath + "/mint_tokens.cdc"
	metayaUtilityCoinBurnTokensByAdminPath      = metayaUtilityCoinTransactionsRootPath + "/burn_tokens_by_admin.cdc"
	metayaUtilityCoinBurnTokensByUserPath       = metayaUtilityCoinTransactionsRootPath + "/burn_tokens_by_user.cdc"
	
	metayaUtilityCoinGetSupplyPath              = metayaUtilityCoinScriptsRootPath + "/get_supply.cdc"
	metayaUtilityCoinGetBalancePath             = metayaUtilityCoinScriptsRootPath + "/get_balance.cdc"
)

func DeployContracts(
	t *testing.T,
	b *emulator.Blockchain,
) (flow.Address, flow.Address, crypto.Signer) {
	accountKeys := sdktest.AccountKeyGenerator()

	// Should be able to deploy a contract as a new account with no keys.
	fungibleTokenCode := ftcontracts.FungibleToken()
	fungibleTokenAddress, err := b.CreateAccount(
		[]*flow.AccountKey{},
		[]sdktemplates.Contract{
			{
				Name:   "FungibleToken",
				Source: string(fungibleTokenCode),
			},
		},
	)
	assert.NoError(t, err)

	_, err = b.CommitBlock()
	assert.NoError(t, err)

	metayaUtilityCoinAccountKey, metayaUtilityCoinSigner := accountKeys.NewWithSigner()
	metayaUtilityCoinCode := LoadMetayaUtilityCoin(fungibleTokenAddress.String())

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

	// Simplify testing by having the contract address also be our initial Vault.
	SetupAccount(t, b, metayaUtilityCoinAddress, metayaUtilityCoinSigner, fungibleTokenAddress, metayaUtilityCoinAddress)

	return fungibleTokenAddress, metayaUtilityCoinAddress, metayaUtilityCoinSigner
}

func SetupAccount(
	t *testing.T,
	b *emulator.Blockchain,
	userAddress flow.Address,
	userSigner crypto.Signer,
	fungibleTokenAddress flow.Address,
	metayaUtilityCoinAddress flow.Address,
) {
	tx := flow.NewTransaction().
		SetScript(SetupAccountTransaction(fungibleTokenAddress.String(), metayaUtilityCoinAddress.String())).
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

func CreateAccount(
	t *testing.T,
	b *emulator.Blockchain,
	fungibleTokenAddress flow.Address,
	metayaUtilityCoinAddress flow.Address,
) (flow.Address, crypto.Signer) {
	userAddress, userSigner, _ := test.CreateAccount(t, b)
	SetupAccount(t, b, userAddress, userSigner, fungibleTokenAddress, metayaUtilityCoinAddress)
	return userAddress, userSigner
}

func MintTokens(
	t *testing.T,
	b *emulator.Blockchain,
	fungibleTokenAddress flow.Address,
	metayaUtilityCoinAddress flow.Address,
	metayaUtilityCoinSigner crypto.Signer,
	recipientAddress flow.Address,
	amount string,
	shouldRevert bool,
) {
	tx := flow.NewTransaction().
		SetScript(MintTokensTransaction(fungibleTokenAddress.String(), metayaUtilityCoinAddress.String())).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(metayaUtilityCoinAddress)

	_ = tx.AddArgument(cadence.NewAddress(recipientAddress))
	_ = tx.AddArgument(test.CadenceUFix64(amount))

	test.SignAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, metayaUtilityCoinAddress},
		[]crypto.Signer{b.ServiceKey().Signer(), metayaUtilityCoinSigner},
		shouldRevert,
	)
}

func replaceAddressPlaceholders(code, fungibleTokenAddress, metayaUtilityCoinAddress string) []byte {
	return []byte(test.ReplaceImports(
		code,
		map[string]*regexp.Regexp{
			fungibleTokenAddress:          test.FungibleTokenAddressPlaceholder,
			metayaUtilityCoinAddress:      test.MetayaUtilityCoinAddressPlaceHolder,
		},
	))
}

func LoadMetayaUtilityCoin(fungibleTokenAddress string) []byte {
	return []byte(test.ReplaceImports(
		string(test.ReadFile(metayaUtilityCoinContractPath)),
		map[string]*regexp.Regexp{
			fungibleTokenAddress: test.FungibleTokenAddressPlaceholder,
		},
	))
}

func SetupAccountTransaction(fungibleTokenAddress, metayaUtilityCoinAddress string) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaUtilityCoinSetupAccountPath)),
		fungibleTokenAddress,
		metayaUtilityCoinAddress,
	)
}

func TransferTokensTransaction(fungibleTokenAddress, metayaUtilityCoinAddress string) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaUtilityCoinTransferTokensPath)),
		fungibleTokenAddress,
		metayaUtilityCoinAddress,
	)
}

func TransferManyAccountTransaction(fungibleTokenAddress, metayaUtilityCoinAddress string) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaUtilityCoinTransferManyAccountsPath)),
		fungibleTokenAddress,
		metayaUtilityCoinAddress,
	)
}

func MintTokensTransaction(fungibleTokenAddress, metayaUtilityCoinAddress string) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaUtilityCoinMintTokensPath)),
		fungibleTokenAddress,
		metayaUtilityCoinAddress,
	)
}

func BurnTokensByAdminTransaction(fungibleTokenAddress, metayaUtilityCoinAddress string) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaUtilityCoinBurnTokensByAdminPath)),
		fungibleTokenAddress,
		metayaUtilityCoinAddress,
	)
}

func BurnTokensByUserTransaction(fungibleTokenAddress, metayaUtilityCoinAddress string) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaUtilityCoinBurnTokensByUserPath)),
		fungibleTokenAddress,
		metayaUtilityCoinAddress,
	)
}

func GetSupplyScript(fungibleTokenAddress, metayaUtilityCoinAddress string) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaUtilityCoinGetSupplyPath)),
		fungibleTokenAddress,
		metayaUtilityCoinAddress,
	)
}

func GetBalanceScript(fungibleTokenAddress, metayaUtilityCoinAddress string) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(metayaUtilityCoinGetBalancePath)),
		fungibleTokenAddress,
		metayaUtilityCoinAddress,
	)
}
