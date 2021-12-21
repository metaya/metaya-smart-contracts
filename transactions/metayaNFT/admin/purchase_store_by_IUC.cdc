import FungibleToken from "../../../contracts/FungibleToken.cdc"
import MetayaUtilityCoin from "../../../contracts/MetayaUtilityCoin.cdc"
import MetayaBeneficiaryCut from "../../../contracts/MetayaBeneficiaryCut.cdc"

transaction(saleID: UInt32, purchaseAmount: UFix64, commonwealName: String) {

    // Local variable for the purchaser
    let payRef: &MetayaUtilityCoin.Vault

    prepare(acct: AuthAccount) {
        self.payRef = acct.borrow<&MetayaUtilityCoin.Vault>(from: MetayaUtilityCoin.VaultStoragePath)
            ?? panic("Could not borrow reference to the owner's Vault!")
    }

    execute {
        if (self.payRef.balance >= purchaseAmount) {

            let purchaseVault <- self.payRef.withdraw(amount: purchaseAmount) as! @MetayaUtilityCoin.Vault

            if (commonwealName != "null") {
                // Commonweal Cut
                let commonwealCutPercentage = MetayaBeneficiaryCut.getCommonwealCutPercentage(name: commonwealName)
                    ?? panic("Cannot find the commonweal cutPercentage by the name")
                let commonwealCutAmount = purchaseAmount * commonwealCutPercentage
                let commonwealCut <- purchaseVault.withdraw(amount: commonwealCutAmount)

                let commonwealCap = MetayaBeneficiaryCut.getCommonwealCapability(name: commonwealName)
                    ?? panic("Cannot find the commonweal by the name")
                let commonwealReceiverRef = commonwealCap.borrow()
                    ?? panic("Cannot find commonweal token receiver")
                commonwealReceiverRef.deposit(from: <-commonwealCut)
            }

            // Copyright owners Cut
            let tokenAmount = purchaseVault.balance
            for name in MetayaBeneficiaryCut.getStoreCopyrightOwnerNames(saleID: saleID)! {
                let copyrightOwnerCutPercentage = MetayaBeneficiaryCut.getStoreCutPercentage(saleID: saleID, name: name)
                    ?? panic("Cannot find the copyright owner cutPercentage by the name")
                let copyrightOwnerCutAmount = tokenAmount * copyrightOwnerCutPercentage
                let copyrightOwnerCut <- purchaseVault.withdraw(amount: copyrightOwnerCutAmount)

                let copyrightOwnerCap = MetayaBeneficiaryCut.getCopyrightOwnerCapability(name: name)
                    ?? panic("Cannot find the copyright owner by the name")
                let copyrightOwnerReceiverRef = copyrightOwnerCap.borrow()
                    ?? panic("Cannot find copyright owner token receiver")
                copyrightOwnerReceiverRef.deposit(from: <-copyrightOwnerCut)
            }
            destroy purchaseVault
        } else{
            panic("There's not enough MetayaUtilityCoin in the account")
        }
    }
}