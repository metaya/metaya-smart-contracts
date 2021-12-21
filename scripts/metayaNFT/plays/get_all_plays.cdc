import Metaya from "../../../contracts/Metaya.cdc"

// This script returns an array of all the plays
// that have ever been created for Metaya

// Returns: [Metaya.Play]
// array of all plays created for Metaya

pub fun main(): [Metaya.Play] {

    return Metaya.getAllPlays()
}