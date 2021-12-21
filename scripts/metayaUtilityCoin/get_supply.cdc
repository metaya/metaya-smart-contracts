import MetayaUtilityCoin from "../../contracts/MetayaUtilityCoin.cdc"

pub fun main(): UFix64 {

    let supply = MetayaUtilityCoin.totalSupply

    log(supply)

    return supply
}