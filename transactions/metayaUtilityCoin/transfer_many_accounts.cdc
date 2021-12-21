import FungibleToken from "../../contracts/FungibleToken.cdc"
import MetayaUtilityCoin from "../../contracts/MetayaUtilityCoin.cdc"

transaction(addressAmountMap: {Address: UFix64}) {

    let vaultRef: &MetayaUtilityCoin.Vault

    prepare(signer: AuthAccount) {

        self.vaultRef = signer.borrow<&MetayaUtilityCoin.Vault>(from: MetayaUtilityCoin.VaultStoragePath)
			?? panic("Could not borrow reference to the owner's Vault!")
    }

    execute {

        for address in addressAmountMap.keys {

            let sentVault <- self.vaultRef.withdraw(amount: addressAmountMap[address]!)

            let recipient = getAccount(address)

            let receiverRef = recipient.getCapability(MetayaUtilityCoin.ReceiverPublicPath)!
                .borrow<&{FungibleToken.Receiver}>()
                ?? panic("Could not borrow receiver reference to the recipient's Vault")

            receiverRef.deposit(from: <-sentVault)
        }
    }
}