import * as Promise from 'bluebird';
import { observable, observe, computed, action, runInAction, transaction, IObservableArray, IObjectChange } from "mobx";
import ImagesApi from '../libs/api/ImagesApi';
import AsyncCmd from '../libs/asyncCmd';
import ApiError from '../libs/api/ApiError';
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
    readonly updateImage: AsyncCmd<void, ApiError, string, Image>;
    readonly pagination: Pagination = new Pagination();

    readonly disposeListOptions: () => any;

    constructor(imagesApi: ImagesApi) {
        this.imagesApi = imagesApi;
        this.fetchImages = new AsyncCmd(this.fetchImagesImpl);
        this.saveImage = new AsyncCmd(this.saveImageImpl);
        this.updateImage = new AsyncCmd(this.updateImageImpl);
        this.listOptions = new ListOptions();

        this.disposeListOptions = observe(this.listOptions, this.handleListOptionsChange);
    }

    findByName(imageName: string): Image {
        return this.images.find(image => image.name == imageName);
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


    private fetchImagesImpl = (): Promise<void> => {
        return this.imagesApi.getImages(this.listOptions, this.lastFetchedAt)
            .then(({ images, pagination }) => {
                const newImages = images.map((i: IImage) => new Image(i))
                const newPagination = new Pagination(pagination)

                transaction(() => {
                    this.setImages(newImages);
                    this.setPagination(newPagination);
                });
            })
            .then(() => { });
    }


    private saveImageImpl(image: Image): Promise<void> {
        return this.imagesApi.create(image).then(() => { });
    }

    private updateImageImpl(key: string, image: Image): Promise<void> {
        return this.imagesApi.update(key, image)
            .then(this.refreshImage.bind(key))
            .then(() => { });
    }

    private refreshImage(oldKey: string, key: string): Promise<Image> {
        return this.imagesApi.get(key)
            .then(this.updateImageByKey.bind(null, oldKey))
    }

    @action.bound
    private updateImageByKey(oldKey: string, image: Image) {
        const imageIdx = this.images.findIndex(i => i.key === oldKey);
        if (imageIdx < 0) {
            console.warn(`Update image failed: image with key '${oldKey}' not found`);
            return image;
        }

        this.images[imageIdx] = image;
        return image;
    }
}
