import ApiClient from "./ApiClient";

interface Owners {
    [provider: string]: Array<string>;
}

interface Repositories {
    [providerOwner: string]: Array<string>;
}

interface FetchingRepositories {
    [key: string]: Promise<Array<string>>;
}

interface FetchingOwners {
    [key: string]: Promise<Array<string>>;
}

export default class SourceApi {
    apiTokens = { github: 'tokenHere' };
    owners: Owners = { github: ['umurgdk', 'mobingilabs'] };
    repositories: Repositories = { 'github:umurgdk': ['soundlines', 'XPlayer'] };
    fetchingOwners: FetchingOwners = {};
    fetchingRepositories: FetchingRepositories = {};

    apiClient: ApiClient;

    constructor(apiClient: ApiClient) {
        this.apiClient = apiClient;
    }

    getOwners(provider: string): Promise<Array<string>> {
        if (this.owners[provider] && this.owners[provider].length > 0) {
            return Promise.resolve(this.owners[provider]);
        }
        
        return this.loadOwners(provider);
    }

    getRepositories(provider: string, owner: string): Promise<Array<string>> {
        let ownerKey = `${provider}:${owner}`;
        if (this.repositories[ownerKey] && this.repositories[ownerKey].length > 0) {
            return Promise.resolve(this.repositories[ownerKey]);
        }

        return this.loadRepositories(provider, owner);
    }

    loadOwners(provider: string) {
        // TODO: Implement fetching for SourceApi::fetchOwners
        if (!this.fetchingOwners[provider]) {
            this.fetchingOwners[provider] = Promise.resolve(['umurgdk', 'mobingilabs', 'mobingi'])
                .then((owners: Array<string>) => {
                    this.fetchingOwners[provider] = null;
                    this.owners[provider] = owners;
                    return this.owners[provider];
                });
        }

        return this.fetchingOwners[provider];
    }

    loadRepositories(provider: string, owner: string): Promise<Array<string>> {
        // TODO: Implement fetching for SourceApi::fetchRepositories
        const ownerKey = `${provider}:${owner}`;

        if (!this.fetchingRepositories[ownerKey]) {
            this.fetchingRepositories[ownerKey] = Promise.resolve(['test-repo'])
                .then((repositories: Array<string>) => {
                    this.fetchingRepositories[ownerKey] = null;
                    this.repositories[ownerKey] = (this.repositories[ownerKey] || []).concat(repositories);
                    return this.repositories[ownerKey];
                });
        }

        return this.fetchingRepositories[ownerKey];
    }
}
