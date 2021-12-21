import MetayaBeneficiaryCut from "../../contracts/MetayaBeneficiaryCut.cdc"

transaction(name: String) {

    let adminRef: &MetayaBeneficiaryCut.Admin

    prepare(acct: AuthAccount) {
        self.adminRef = acct.borrow<&MetayaBeneficiaryCut.Admin>(from: MetayaBeneficiaryCut.AdminStoragePath)
            ?? panic("No admin resource in storage")
    }

    execute {
        self.adminRef.setCopyrightOwnerCapability(name: name, capability: nil)
    }
}