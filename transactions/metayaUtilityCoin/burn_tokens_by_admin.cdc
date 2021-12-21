import FungibleToken from "../../contracts/FungibleToken.cdc"
import MetayaUtilityCoin from "../../contracts/MetayaUtilityCoin.cdc"

transaction(amount: UFix64) {

    let vault: @FungibleToken.Vault
    let admin: &MetayaUtilityCoin.Administrator

    prepare(tokenAdmin: AuthAccount) {

        self.vault <- tokenAdmin.borrow<&MetayaUtilityCoin.Vault>(from: MetayaUtilityCoin.VaultStoragePath)!
            .withdraw(amount: amount)

        self.admin = tokenAdmin.borrow<&MetayaUtilityCoin.Administrator>(from: MetayaUtilityCoin.AdminStoragePath)
            ?? panic("Could not borrow a reference to the admin resource")
    }

    execute {
        
        let burner <- self.admin.createNewBurner()
        burner.burnTokens(from: <-self.vault)

        destroy burner
    }
}