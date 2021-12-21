import MetayaBeneficiaryCut from "../../contracts/MetayaBeneficiaryCut.cdc"

pub fun main(): Address {

    return MetayaBeneficiaryCut.metayaCapability.address

}