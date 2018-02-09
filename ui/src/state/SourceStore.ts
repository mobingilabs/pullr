import * as Promise from 'bluebird';
import { observable, ObservableMap, computed, action } from "mobx";
import * as fuzzysearch from 'fuzzysearch';
import SourceApi from '../libs/api/SourceApi';
import AsyncCmd from "../libs/asyncCmd";
import ApiError from "../libs/api/ApiError";

export default class SourceStore {
    @observable organisations: ObservableMap<string[]>;
    @observable repositories: ObservableMap<string[]>;
    @observable selectedProvider: string;
    @observable selectedOrganisation: string;
    @observable repositoryFilter: string;

    @observable loadOrganisations: AsyncCmd<void>;
    @observable loadRepositories: AsyncCmd<void>

    private sourceApi: SourceApi;

    constructor(sourceApi: SourceApi) {
        this.sourceApi = sourceApi;

        this.organisations = new ObservableMap({});
        this.repositories = new ObservableMap({});
        this.repositoryFilter = '';

        this.loadOrganisations = new AsyncCmd(this.loadOrganisationsImpl);
        this.loadRepositories = new AsyncCmd(this.loadRepositoriesImpl);
    }

    @computed
    get currentOrganisations(): string[] {
        if (this.organisations.has(this.selectedProvider)) {
            const organisations = this.organisations.get(this.selectedProvider);
            if (organisations && organisations.length > 0) {
                return organisations;
            }
        }

        return [];
    }

    @computed
    get currentRepositories(): string[] {
        const key = this.repoKey(this.selectedProvider, this.selectedOrganisation);
        if (!this.repositories.has(key)) {
            return [];
        }

        const repos = this.repositories.get(key);
        if (this.repositoryFilter.length > 0) {
            return repos.filter(r => fuzzysearch(this.repositoryFilter, r));
        }

        return repos;
    }

    private repoKey(provider: string, organisation: string): string {
        return `${provider}::${organisation}`;
    }

    @action.bound
    private loadOrganisationsImpl(): Promise<void> {
        return this.sourceApi.getOrganisations(this.selectedProvider)
            .tap(this.setOrganisations.bind(null, this.selectedProvider))
            .then(() => { });
    }

    @action.bound
    private loadRepositoriesImpl(): Promise<void> {
        return this.sourceApi.getRepositories(this.selectedProvider, this.selectedOrganisation)
            .tap(this.setRepositories.bind(null, this.selectedProvider, this.selectedOrganisation))
            .then(() => { });
    }

    @action.bound
    setOrganisations(provider: string, organisations: string[]) {
        this.organisations.set(provider, organisations);
    }

    @action.bound
    setRepositories(provider: string, organisation: string, repositories: string[]) {
        this.repositories.set(this.repoKey(provider, organisation), repositories);
    }

    @action.bound
    selectProvider(provider: string) {
        this.selectedProvider = provider;
    }

    @action.bound
    selectOrganisation(organisation: string) {
        this.selectedOrganisation = organisation;
    }

    @action.bound
    resetSelections() {
        this.selectedProvider = null;
        this.selectedOrganisation = null;
        this.repositoryFilter = '';
    }
}
