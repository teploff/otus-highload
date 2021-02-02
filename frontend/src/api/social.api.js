import httpClient from "@/api/httpClient";

const SEARCH_BY_ANTHROPONYM_ENDPOINT = 'social/profile/search-by-anthroponym'
const CREATE_FRIENDSHIP_ENDPOINT = 'social/friendship/create'

const searchByAnthroponym = (payload) => httpClient.get(SEARCH_BY_ANTHROPONYM_ENDPOINT, {params: payload})
const createFriendship = (usersID) => httpClient.post(CREATE_FRIENDSHIP_ENDPOINT, {friends_id: usersID})

export {
    searchByAnthroponym,
    createFriendship
}