import * as Promise from 'bluebird';
import ApiClient from "./ApiClient";

export default class OAuthApi {
    apiClient: ApiClient;

    constructor(apiClient: ApiClient) {
        this.apiClient = apiClient;
    }

    getLoginUrl(provider: string): Promise<string> {
        const cbUri = `${location.origin}/oauth`;
        const url = `/oauth/${provider}/url?cb=${encodeURIComponent(cbUri)}`;
        return this.apiClient.json('get', url).then(({ login_url }) => login_url);
    }
}

