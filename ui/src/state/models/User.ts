import { observable, ObservableMap, IObservable, IObservableMapInitialValues } from 'mobx';

export interface IUser {
    username: string;
    tokens: IObservableMapInitialValues<{ [provider: string]: string }>;
}

export default class User implements IUser {
    @observable username: string;
    @observable tokens: ObservableMap<{ [provider: string]: string }>;

    constructor(json: Partial<IUser> = {}) {
        this.username = json.username || '';
        this.tokens = new ObservableMap(json.tokens || {});
    }

    clone(): User {
        return new User({ username: this.username, tokens: this.tokens.toJS() });
    }
}
