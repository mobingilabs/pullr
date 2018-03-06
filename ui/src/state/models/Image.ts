import { observable, computed, action } from 'mobx';
import { IImageTag, ImageTag } from './ImageBuild';
import { IRepository, Repository } from './Repository';
import Status, { IStatus } from './Status';

export const Statuses = {
    Ready: 'image:ready',
    Building: 'image:building',
}

export const Causes = {
    BuildStart: 'image:build:start',
    BuildFail: 'image:build:fail',
    BuildSuccess: 'image:build:success',
    Delete: 'image:delete',
}

export interface IImage {
    key?: string;
    name: string;
    repository: IRepository;
    tags: IImageTag[];
    dockerfile_path: string;
    created_at: Date;
    updated_at: Date;
    status?: IStatus;
}

export default class Image implements IImage {
    @observable key: string;
    @observable name: string;
    @observable repository: Repository;
    @observable dockerfile_path: string;
    @observable tags: ImageTag[];
    @observable created_at: Date;
    @observable updated_at: Date;
    @observable status?: Status;

    constructor(json: IImage) {
        this.key = json.key;
        this.name = json.name;
        this.repository = new Repository(json.repository);
        this.tags = json.tags.map(t => new ImageTag(t));
        this.dockerfile_path = json.dockerfile_path;
        this.status = json.status && new Status(json.status);
    }

    @action.bound
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
            updated_at: this.updated_at,
            status: this.status && this.status.clone()
        });
    }

    toJS(): IImage {
        const { key, name, repository, dockerfile_path, tags, created_at, updated_at, status } = this;
        const jsTags = tags.map(tag => tag.toJS());
        const jsRepository = repository.toJS();

        return { key, name, repository: jsRepository, dockerfile_path, tags: jsTags, created_at, updated_at, status };
    }

    static create(): Image {
        return new Image({
            name: '',
            repository: Repository.create(),
            dockerfile_path: './Dockerfile',
            tags: [
                new ImageTag({ ref_type: 'branch', ref_test: 'master', name: 'latest' })
            ],
            created_at: new Date(),
            updated_at: new Date()
        });
    }
}
