import { observable } from 'mobx';

export interface IRepository {
    provider: string;
    owner: string;
    name: string;
}

export class Repository implements IRepository {
    @observable provider: string;
    @observable owner: string;
    @observable name: string;

    constructor(json: IRepository) {
        this.provider = json.provider;
        this.owner = json.owner;
        this.name = json.name;
    }

    static create(): Repository {
        return new Repository({
            name: '',
            owner: '',
            provider: '',
        });
    }

    clone(): Repository {
        const { provider, owner, name } = this;
        return new Repository({ provider, owner, name });
    }

    toJS(): IRepository {
        const { provider, owner, name } = this;
        return { provider, owner, name };
    }
}
