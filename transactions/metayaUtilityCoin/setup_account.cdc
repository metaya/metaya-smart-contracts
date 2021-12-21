import FungibleToken from "../../contracts/FungibleToken.cdc"
import MetayaUtilityCoin from "../../contracts/MetayaUtilityCoin.cdc"

transaction {

    prepare(signer: AuthAccount) {

        if signer.borrow<&MetayaUtilityCoin.Vault>(from: MetayaUtilityCoin.VaultStoragePath) != nil {
            return
        }

        signer.save(
            <-MetayaUtilityCoin.createEmptyVault(),
            to: MetayaUtilityCoin.VaultStoragePath
        )

        signer.link<&MetayaUtilityCoin.Vault{FungibleToken.Receiver}>(
            MetayaUtilityCoin.ReceiverPublicPath,
            target: MetayaUtilityCoin.VaultStoragePath
        )

        signer.link<&MetayaUtilityCoin.Vault{FungibleToken.Balance}>(
            MetayaUtilityCoin.BalancePublicPath,
            target: MetayaUtilityCoin.VaultStoragePath
        )
    }
}