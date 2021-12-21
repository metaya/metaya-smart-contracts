import MetayaBeneficiaryCut from "../../contracts/MetayaBeneficiaryCut.cdc"

pub fun main(name: String): Address {

    return MetayaBeneficiaryCut.getCopyrightOwnerCapability(name: name)!.address

}