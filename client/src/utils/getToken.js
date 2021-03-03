import { useAuth0 } from "@auth0/auth0-react"

export const useGetToken = async ()=>{
    const {getIdTokenClaims} = useAuth0()
    const token = await getIdTokenClaims()
    return token?.__raw
}