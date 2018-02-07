import { observable } from "mobx";

export interface IImageTag {
    ref_type: string;
    ref_test: string;
    name: string;
}

export class ImageTag implements IImageTag {
    @observable ref_type: string;
    @observable ref_test: string;
    @observable name: string;

    constructor(json: IImageTag) {
        this.ref_type = json.ref_type;
        this.ref_test = json.ref_test;
        this.name = json.name;
    }

    static create(): ImageTag {
        return new ImageTag({ ref_type: 'branch', name: '', ref_test: '' });
    }

    clone(): ImageTag {
        return new ImageTag({
            ref_type: this.ref_type,
            ref_test: this.ref_test,
            name: this.name,
        });
    }

    toJS(): IImageTag {
        const { ref_type, ref_test, name } = this;
        return { ref_type, ref_test, name };
    }
}
