import * as Promise from 'bluebird';
import { observable, computed, action, runInAction, IObservableArray } from "mobx";
import Image, { IImage } from "./models/Image";
import Pagination from "./models/Pagination";
import ImagesApi from '../libs/api/ImagesApi';
import AsyncCmd from '../libs/asyncCmd';
import ApiError from '../libs/api/ApiError';

export default class ImageStore {
    imagesApi: ImagesApi;
    lastImageCreation: Date;

    @observable readonly images = observable<Image>([]);
    @observable readonly fetchImages: AsyncCmd<void>;
    @observable readonly saveImage: AsyncCmd<void>;

    constructor(imagesApi: ImagesApi) {
        this.imagesApi = imagesApi;
        this.images = observable([]);
        this.fetchImages = new AsyncCmd(this.fetchImagesImpl);
        this.saveImage = new AsyncCmd(this.saveImageImpl);
    }

    findByName(imageName: string): Image {
        return this.images.find(image => image.name == imageName);
    }

    @action.bound
    fetchImagesImpl(): Promise<void> {
        return this.imagesApi.getImages(this.lastImageCreation)
            .then(images => images.map((i: IImage) => new Image(i)))
            .tap(this.setImages)
            .then(() => { });
    }

    @action.bound
    setImages(images: Image[]) {
        if (this.lastImageCreation != null) {
            this.images.unshift(...images);
        } else {
            this.images.replace(images);
        }

        this.lastImageCreation = images[0].created_at;
    }

    @action.bound
    saveImageImpl(image: Image): Promise<void> {
        return this.imagesApi.create(image).then(() => { });
    }

    @action.bound
    updateImage(key: string, image: Image): Promise<void> {
        return this.imagesApi.update(key, image)
            .then(this.refreshImage.bind(key))
            .then(() => { });
    }

    @action.bound
    refreshImage(oldKey: string, key: string): Promise<Image> {
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
