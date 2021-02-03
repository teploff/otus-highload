import httpClient from "@/api/httpClient";

const SEARCH_BY_ANTHROPONYM_ENDPOINT = 'social/profile/search-by-anthroponym'

const CREATE_FRIENDSHIP_ENDPOINT = 'social/friendship/create'
const CONFIRM_FRIENDSHIP_ENDPOINT = 'social/friendship/confirm'
const REJECT_FRIENDSHIP_ENDPOINT = 'social/friendship/reject'
const SPLIT_UP_FRIENDSHIP_ENDPOINT = 'social/friendship/split-up'
const GET_FRIENDS_FRIENDSHIP_ENDPOINT = 'social/friendship/get-friends'
const GET_FOLLOWERS_FRIENDSHIP_ENDPOINT = 'social/friendship/get-followers'

const searchByAnthroponym = (payload) => httpClient.get(SEARCH_BY_ANTHROPONYM_ENDPOINT, {params: payload})
const createFriendship = (usersID) => httpClient.post(CREATE_FRIENDSHIP_ENDPOINT, {friends_id: usersID})
const confirmFriendship = (usersID) => httpClient.post(CONFIRM_FRIENDSHIP_ENDPOINT, {friends_id: usersID})
const rejectFriendship = (usersID) => httpClient.post(REJECT_FRIENDSHIP_ENDPOINT, {friends_id: usersID})
const splitUpFriendship = (usersID) => httpClient.post(SPLIT_UP_FRIENDSHIP_ENDPOINT, {friends_id: usersID})
const getFriends = () => httpClient.get(GET_FRIENDS_FRIENDSHIP_ENDPOINT)
const getFollowers = () => httpClient.get(GET_FOLLOWERS_FRIENDSHIP_ENDPOINT)

export {
    searchByAnthroponym,
    createFriendship,
    confirmFriendship,
    rejectFriendship,
    splitUpFriendship,
    getFriends,
    getFollowers
}
