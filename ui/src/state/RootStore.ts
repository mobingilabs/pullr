import * as Promise from 'bluebird';
import { observable, computed, action } from 'mobx';
import ImageStore from './ImagesStore';
import ApiClient from '../libs/api/ApiClient';
import ImagesApi from '../libs/api/ImagesApi';
import User from './models/User';
import AuthApi from '../libs/api/AuthApi';
import AuthStore from './AuthStore';
import OAuthApi from '../libs/api/OAuthApi';
import AsyncCmd from '../libs/asyncCmd';
import SourceApi from '../libs/api/SourceApi';
import SourceStore from './SourceStore';

export default class RootStore {
    images: ImageStore;
    auth: AuthStore;
    sources: SourceStore;

    @observable init: AsyncCmd<void>;
    @observable user: User;

    constructor(imagesApi: ImagesApi, authApi: AuthApi, oauthApi: OAuthApi, sourceApi: SourceApi) {
        this.images = new ImageStore(imagesApi);
        this.auth = new AuthStore(authApi, oauthApi);
        this.sources = new SourceStore(sourceApi);
        this.init = new AsyncCmd(this.initImpl);
    }

    initImpl = (): Promise<void> => {
        return this.auth.loadProfile.run();
    }
}
