import FungibleToken from "../../contracts/FungibleToken.cdc"
import MetayaUtilityCoin from "../../contracts/MetayaUtilityCoin.cdc"

transaction(amount: UFix64, to: Address) {

    let sentVault: @FungibleToken.Vault

    prepare(signer: AuthAccount) {

        let vaultRef = signer.borrow<&MetayaUtilityCoin.Vault>(from: MetayaUtilityCoin.VaultStoragePath)
			?? panic("Could not borrow reference to the owner's Vault!")

        self.sentVault <- vaultRef.withdraw(amount: amount)
    }

    execute {

        let recipient = getAccount(to)

        let receiverRef = recipient.getCapability(MetayaUtilityCoin.ReceiverPublicPath)!
            .borrow<&{FungibleToken.Receiver}>()
			?? panic("Could not borrow receiver reference to the recipient's Vault")

        receiverRef.deposit(from: <-self.sentVault)
    }
}