import Metaya from "../../contracts/Metaya.cdc"

// This script reads the current number of moments that have been minted
// from the Metaya contract and returns that number to the caller

// Returns: UInt64
// Number of moments minted from Metaya contract

pub fun main(): UInt64 {

    return Metaya.totalSupply
}