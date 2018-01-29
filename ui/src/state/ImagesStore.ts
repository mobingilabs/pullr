import { observable, computed, action, IObservableArray } from "mobx";
import Image from "./models/Image";
import Pagination from "./models/Pagination";
import ImagesApi from '../libs/api/ImagesApi';

export default class ImageStore {
    @observable images: IObservableArray<Image>;
    @observable pagination: Pagination;
    @observable syncInProgress: boolean;
    imagesApi: ImagesApi;

    constructor(imagesApi: ImagesApi) {
        this.imagesApi = imagesApi;
        this.images = observable([]);
        this.pagination = new Pagination();
        this.syncInProgress = false;
    }

    findByName(imageName: string): Image {
        return this.images.find(image => image.name == imageName);
    }

    @action setImages(images: IObservableArray<Image>) {
        this.images = images;
        this.pagination.totalItems = Math.max(this.images.length, this.pagination.totalItems);
    }

    @action saveImage(image: Image) {
        this.images.unshift(image);
    }

    @action updateImage(imageName: string, image: Image) {
        const index = this.images.findIndex(image => image.name === imageName);
        if (index >= 0) {
            this.images[index] = image;
        }
    }

    @action loadPage(pageNumber: number) {

    }
}
