import MetayaBeneficiaryCut from "../../contracts/MetayaBeneficiaryCut.cdc"

pub fun main(): [String] {

    return MetayaBeneficiaryCut.getAllCopyrightOwnerNames()

}