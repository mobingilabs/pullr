import * as Promise from 'bluebird';
import { observable, action, computed, ObservableMap } from 'mobx';
import { History } from 'history';
import User from './models/User';
import AsyncCmd from '../libs/asyncCmd';
import AuthApi from '../libs/api/AuthApi';
import ApiError from '../libs/api/ApiError';
import OAuthApi from '../libs/api/OAuthApi';

export default class AuthStore {
    @observable user: User;
    @observable readonly loginUrls: ObservableMap<string>;
    @observable readonly loadProfile: AsyncCmd<void>;
    @observable readonly login: AsyncCmd<void, ApiError, string, string>;
    @observable readonly register: AsyncCmd<void, ApiError, string, string, string>;
    @observable readonly oauthInProgress: ObservableMap<boolean>;

    @observable private readonly loginUrlCmds: ObservableMap<AsyncCmd<string, ApiError>>;

    private authApi: AuthApi;
    private oauthApi: OAuthApi;

    constructor(authApi: AuthApi, oauthApi: OAuthApi) {
        this.authApi = authApi;
        this.oauthApi = oauthApi;

        this.loginUrls = new ObservableMap({});
        this.loadProfile = new AsyncCmd(this.loadProfileImpl);
        this.login = new AsyncCmd(this.loginImpl);
        this.register = new AsyncCmd(this.registerImpl);
        this.loginUrlCmds = new ObservableMap({});
        this.oauthInProgress = new ObservableMap({});
        // this.register = new AsyncCmd(this.registerImpl);

        window.addEventListener('message', this.handleMessage, false);
    }

    @computed
    get loggedIn(): boolean {
        return this.user && this.user.username && this.user.username.length > 0;
    }

    getLoginUrl(provider: string): AsyncCmd<string, ApiError> {
        if (this.loginUrlCmds.has(provider)) {
            return this.loginUrlCmds.get(provider);
        }

        const cmd = new AsyncCmd<string, ApiError>(() => this.getLoginUrlImpl(provider));
        this.loginUrlCmds.set(provider, cmd);
        return cmd;
    }

    @action.bound
    oauthStart(provider: string) {
        this.oauthInProgress.set(provider, true);
    }

    @action.bound
    private getLoginUrlImpl(provider: string): Promise<string> {
        return this.oauthApi.getLoginUrl(provider);
    }

    @action.bound
    private setLoginUrlForProvider(provider: string, url: string) {
        this.loginUrls.set(provider, url);
    }

    @action.bound
    private loadProfileImpl(): Promise<void> {
        return this.authApi.getProfile().tap(this.setUser).then(() => { });
    }

    @action.bound
    private loginImpl(username: string, password: string): Promise<void> {
        return this.authApi.login(username, password).tap(this.setUser).then(() => { });
    }

    @action.bound
    private registerImpl(username: string, email: string, password: string): Promise<void> {
        return this.authApi.register(username, email, password).tap(this.setUser).then(() => { });
    }

    @action.bound
    private handleMessage(e: MessageEvent) {
        if (e.origin !== location.origin) {
            return
        }

        if (e.data === 'OAUTH_SUCCESS') {
            this.loadProfile.run();
            this.oauthInProgress.clear();
        }
    }

    @action.bound
    setUser(user: User) {
        this.user = user;
    }

    @action.bound
    logout = (history?: History) => {
        this.authApi.logout();
        this.user = null;

        if (history != null) {
            history.replace('/login');
        }
    }
}
