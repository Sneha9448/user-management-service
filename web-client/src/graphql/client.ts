import { Client, cacheExchange, fetchExchange } from 'urql';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8081/graphql';

export const client = new Client({
    url: API_URL,
    exchanges: [cacheExchange, fetchExchange],
    fetchOptions: () => {
        const token = localStorage.getItem('token');
        return {
            headers: {
                Authorization: token ? `Bearer ${token}` : '',
            },
        };
    },
});
