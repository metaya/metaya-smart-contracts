import Metaya from "../../../contracts/Metaya.cdc"

transaction() {
    
    prepare(acct: AuthAccount) {

        let metadata: {String: String} = {"Title": "Artwork 000"}

        let newPlay = Metaya.Play(metadata: metadata)

    }
}