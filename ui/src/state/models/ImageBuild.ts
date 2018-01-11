import { observable } from "mobx";

export interface IImageBuild {
    type: string;
    name: string;
    tag: string;
}

export class ImageBuild implements IImageBuild {
    @observable type: string;
    @observable name: string;
    @observable tag: string;

    constructor(json: IImageBuild) {
        this.type = json.type;
        this.name = json.name;
        this.tag = json.tag;
    }

    static create(): ImageBuild {
        return new ImageBuild({ type: 'branch', name: '', tag: ''});
    }

    clone(): ImageBuild {
        return new ImageBuild({
            type: this.type,
            name: this.name,
            tag: this.tag
        });
    }
}
