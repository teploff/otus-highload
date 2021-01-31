import httpClient from "@/api/httpClient";

const SIGN_UP_ENDPOINT = '/auth/sign-up'
const SIGN_IN_ENDPOINT = '/auth/sign-in'

const signUpUser = (signUpPayload) => httpClient.post(SIGN_UP_ENDPOINT, signUpPayload)
const signInUser = (email, password) => httpClient.post(SIGN_IN_ENDPOINT, {email, password})

export {
    signUpUser,
    signInUser,
}