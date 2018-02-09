import * as Promise from 'bluebird';
import { observable, observe, computed, action, runInAction, transaction, IObservableArray, IObjectChange } from "mobx";

import ImagesApi from '../libs/api/ImagesApi';
import AsyncCmd from '../libs/asyncCmd';
import ApiError from '../libs/api/ApiError';
import debounce from '../libs/debounce';

import ListOptions from './models/ListOptions';
import Image, { IImage } from "./models/Image";
import Pagination, { PageOutOfRangeException } from "./models/Pagination";

export default class ImageStore {
    imagesApi: ImagesApi;
    lastFetchedAt: Date;

    @observable readonly images: IObservableArray<Image> = [] as IObservableArray<Image>;
    readonly listOptions: ListOptions;
    readonly fetchImages: AsyncCmd<void, ApiError>;
    readonly saveImage: AsyncCmd<void, ApiError, Image>;
    readonly updateImage: AsyncCmd<Image, ApiError, string, Image>;
    readonly findByKey: AsyncCmd<Image, ApiError, string>;
    readonly pagination: Pagination = new Pagination();

    readonly disposeListOptions: () => any;

    constructor(imagesApi: ImagesApi) {
        this.imagesApi = imagesApi;
        this.fetchImages = new AsyncCmd(this.fetchImagesImpl);
        this.saveImage = new AsyncCmd(this.saveImageImpl);
        this.updateImage = new AsyncCmd(this.updateImageImpl);
        this.findByKey = new AsyncCmd(this.findByKeyImpl);
        this.listOptions = new ListOptions();

        this.disposeListOptions = observe(this.listOptions, debounce(300, this.handleListOptionsChange));
    }

    @action.bound
    setImages(images: Image[]) {
        if (this.lastFetchedAt != null) {
            this.images.unshift(...images);
        } else {
            this.images.replace(images);
        }

        this.lastFetchedAt = new Date();
    }

    @action.bound
    setPagination(pagination: Pagination) {
        this.pagination.current = pagination.current;
        this.pagination.next = pagination.next;
        this.pagination.last = pagination.last;
        this.pagination.total = pagination.total;
        this.pagination.per_page = pagination.per_page;
    }

    @action.bound
    fetchImagesAtPage(page: number) {
        this.listOptions.page = page;
        this.fetchImages.run().done();
    }

    @action.bound
    private handleListOptionsChange(change: IObjectChange) {
        if (change.name == 'page') {
            return;
        }

        this.lastFetchedAt = null;
        this.fetchImages.run().done();
    }

    private findByKeyImpl = (key: string): Promise<Image> => {
        const image = this.images.find(image => image.key == key);
        if (image) {
            return Promise.resolve(image);
        }

        return this.imagesApi.get(key).then(img => new Image(img));
    }

    private fetchImagesImpl = (): Promise<void> => {
        return this.imagesApi.getImages(this.listOptions, this.lastFetchedAt)
            .tap(({ images, pagination }) => {
                const newImages = images.map((i: IImage) => new Image(i))
                const newPagination = new Pagination(pagination)

                transaction(() => {
                    this.setImages(newImages);
                    this.setPagination(newPagination);
                });
            })
            .then(() => { });
    }


    private saveImageImpl = (image: Image): Promise<void> => {
        return this.imagesApi.create(image).then(() => { });
    }

    private updateImageImpl = (key: string, image: Image): Promise<Image> => {
        return this.imagesApi.update(key, image)
            .then(this.refreshImage.bind(null, key))
    }

    private refreshImage = (oldKey: string, key: string): Promise<Image> => {
        return this.imagesApi.get(key)
            .then(img => new Image(img))
            .tap(this.updateImageByKey.bind(null, oldKey));
    }

    @action.bound
    private updateImageByKey(oldKey: string, image: Image) {
        const imageIdx = this.images.findIndex(i => i.key === oldKey);
        if (imageIdx < 0) {
            this.images.unshift(image);
        }

        this.images[imageIdx] = image;
    }
}
