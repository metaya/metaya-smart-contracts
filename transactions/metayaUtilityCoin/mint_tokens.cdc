import FungibleToken from "../../contracts/FungibleToken.cdc"
import MetayaUtilityCoin from "../../contracts/MetayaUtilityCoin.cdc"

transaction(recipient: Address, amount: UFix64) {

    let tokenAdmin: &MetayaUtilityCoin.Administrator
    let tokenReceiver: &{FungibleToken.Receiver}

    prepare(signer: AuthAccount) {

        self.tokenAdmin = signer.borrow<&MetayaUtilityCoin.Administrator>(from: MetayaUtilityCoin.AdminStoragePath)
            ?? panic("Signer is not the UtilityCoin admin")

        self.tokenReceiver = getAccount(recipient)
            .getCapability(MetayaUtilityCoin.ReceiverPublicPath)!
            .borrow<&{FungibleToken.Receiver}>()
            ?? panic("Unable to borrow receiver reference")
    }

    execute {

        let minter <- self.tokenAdmin.createNewMinter(allowedAmount: amount)
        let mintedVault <- minter.mintTokens(amount: amount)

        self.tokenReceiver.deposit(from: <-mintedVault)

        destroy minter
    }
}