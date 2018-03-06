import { observable } from 'mobx';


export interface IStatus {
    account: string;
    kind: string;
    id: string;
    time: string | Date;
    name: string;
    metadata: string;
    cause: string;
}

export default class Status implements IStatus {
    @observable account: string;
    @observable kind: string;
    @observable id: string;
    @observable time: Date;
    @observable name: string;
    @observable metadata: string;
    @observable cause: string;

    constructor(json: IStatus) {
        this.account = json.account;
        this.kind = json.kind;
        this.id = json.id;
        this.time = new Date(json.time as string);
        this.name = json.name;
        this.metadata = json.metadata;
        this.cause = json.cause;
    }

    clone(): Status {
        return new Status({
            account: this.account,
            kind: this.kind,
            id: this.id,
            time: this.time,
            name: this.name,
            metadata: this.metadata,
            cause: this.cause,
        })
    }
}
