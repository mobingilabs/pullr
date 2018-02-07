import { observable, computed } from "mobx";
import { IImageTag, ImageTag } from "./ImageBuild";
import { IRepository, Repository } from "./Repository";

export interface IImage {
    key?: string;
    name: string;
    repository: IRepository;
    tags: IImageTag[];
    dockerfile_path: string;
    created_at: Date;
    updated_at: Date;
}

export default class Image implements IImage {
    @observable key: string;
    @observable name: string;
    @observable repository: Repository;
    @observable dockerfile_path: string;
    @observable tags: ImageTag[];
    @observable created_at: Date;
    @observable updated_at: Date;

    constructor(json: IImage) {
        this.key = json.key;
        this.name = json.name;
        this.repository = new Repository(json.repository);
        this.tags = json.tags.map(t => new ImageTag(t));
        this.dockerfile_path = json.dockerfile_path;
    }

    addTag() {
        this.tags.push(ImageTag.create());
    }

    removeTag(buildIndex: number) {
        this.tags.splice(buildIndex, 1);
    }

    clone(): Image {
        return new Image({
            name: this.name,
            repository: this.repository.clone(),
            dockerfile_path: this.dockerfile_path,
            tags: this.tags.map(t => t.clone()),
            created_at: this.created_at,
            updated_at: this.updated_at
        });
    }

    toJS(): IImage {
        const { key, name, repository, dockerfile_path, tags, created_at, updated_at } = this;
        const jsTags = tags.map(tag => tag.toJS());
        const jsRepository = repository.toJS();

        return { key, name, repository: jsRepository, dockerfile_path, tags: jsTags, created_at, updated_at };
    }

    static create(): Image {
        return new Image({
            name: '',
            repository: Repository.create(),
            dockerfile_path: './Dockerfile',
            tags: [
                new ImageTag({ ref_type: 'branch', name: 'master', ref_test: 'latest' })
            ],
            created_at: new Date(),
            updated_at: new Date()
        });
    }
}
