import * as Promise from 'bluebird';
import ApiClient from "./ApiClient";

export default class SourceApi {
    apiClient: ApiClient;

    constructor(apiClient: ApiClient) {
        this.apiClient = apiClient;
    }

    getOrganisations(provider: string): Promise<string[]> {
        return this.apiClient.json('get', `/vcs/${provider}/organisations`).then(({ organisations }) => organisations);
    }

    getRepositories(provider: string, organisation: string): Promise<string[]> {
        return this.apiClient.json('get', `/vcs/${provider}/${organisation}/repositories`).then(({ repositories }) => repositories);
    }
}
