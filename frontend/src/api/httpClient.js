import applyCaseMiddleware from 'axios-case-converter';
import axios from 'axios';

import router from '@/router';

const httpClient = applyCaseMiddleware(axios.create({
    baseURL: process.env.VUE_APP_BASE_URL,
    timeout: 5000,
    headers: {
        "Content-Type": "application/json",
    }
}));

const getAccessToken = () => localStorage.getItem('accessToken');
const getRefreshToken = () => localStorage.getItem('refreshToken');

const authInterceptor = (config) => {
    config.headers['Authorization'] = getAccessToken();
    return config;
}

httpClient.interceptors.request.use(authInterceptor);

async function expiredInterceptor(error) {
    if (error && error.response && error.response.status === 401) {
        const ax = axios.create({headers: {Authorization: `${getRefreshToken()}`}});

        const response = await ax.put('/auth/token').catch((e) => e);
        if (response.status === 200) {
            localStorage.setItem('accessToken', response.data.access_token);
            localStorage.setItem('refreshToken', response.data.refresh_token);
            return axios(error.config);
        } else {
            router.push({name: 'SignIn'}).catch(() => void 0);
            return new Error('Token refresh failed');

        }
    } else {
        return error;
    }
}

axios.interceptors.response.use(undefined, expiredInterceptor);

export default httpClient;