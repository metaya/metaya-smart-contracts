import MetayaBeneficiaryCut from "../../contracts/MetayaBeneficiaryCut.cdc"
import FungibleToken from "../../contracts/FungibleToken.cdc"
import MetayaUtilityCoin from "../../contracts/MetayaUtilityCoin.cdc"

transaction(addr: Address) {

    let adminRef: &MetayaBeneficiaryCut.Admin

    prepare(acct: AuthAccount) {
        self.adminRef = acct.borrow<&MetayaBeneficiaryCut.Admin>(from: MetayaBeneficiaryCut.AdminStoragePath)
            ?? panic("No admin resource in storage")
    }

    execute {
        let account = getAccount(addr)
        let cap = account.getCapability<&{FungibleToken.Receiver}>(MetayaUtilityCoin.ReceiverPublicPath)

        self.adminRef.setMetayaCapability(capability: cap)
    }
}