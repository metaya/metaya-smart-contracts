import FungibleToken from "../../contracts/FungibleToken.cdc"
import MetayaUtilityCoin from "../../contracts/MetayaUtilityCoin.cdc"

pub fun main(address: Address): UFix64 {

    let account = getAccount(address)
    
    let vaultRef = account.getCapability(MetayaUtilityCoin.BalancePublicPath)!
        .borrow<&MetayaUtilityCoin.Vault{FungibleToken.Balance}>()
        ?? panic("Could not borrow Balance reference to the Vault")

    return vaultRef.balance
}