import httpClient from "@/api/httpClient";

const SEARCH_BY_ANTHROPONYM_ENDPOINT = 'social/profile/search-by-anthroponym'

const searchByAnthroponym = (payload) => httpClient.get(SEARCH_BY_ANTHROPONYM_ENDPOINT, {params: payload})

export {
    searchByAnthroponym,
}