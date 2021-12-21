import Metaya from "../../../contracts/Metaya.cdc"

// This script gets the setID associated with a moment
// in a collection by getting a reference to the moment
// and then looking up its setID

// Parameters:
//
// account: The Flow Address of the account whose moment data needs to be read
// id: The unique ID for the moment whose data needs to be read

// Returns: UInt32
// The setID associated with a moment with a specified ID

pub fun main(account: Address, id: UInt64): UInt32 {

    // Borrow a public reference to the owner's moment collection 
    let collectionRef = getAccount(account).getCapability(Metaya.CollectionPublicPath)
        .borrow<&{Metaya.MomentCollectionPublic}>()
        ?? panic("Could not get public moment collection reference")

    // Borrow a reference to the specified moment in the collection
    let token = collectionRef.borrowMoment(id: id)
        ?? panic("Could not borrow a reference to the specified moment")

    let data = token.data

    return data.setID
}