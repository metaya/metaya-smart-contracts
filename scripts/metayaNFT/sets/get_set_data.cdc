import Metaya from "../../../contracts/Metaya.cdc"

// This script returns all the metadata about the specified set

// Parameters:
//
// setID: The unique ID for the set whose data needs to be read

// Returns: Metaya.QuerySetData

pub fun main(setID: UInt32): Metaya.QuerySetData {

    let data = Metaya.getSetData(setID: setID)
        ?? panic("Could not get data for the specified set ID")

    return data
}