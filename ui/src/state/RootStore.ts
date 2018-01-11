import { observable, computed } from 'mobx';
import ImageStore from './ImagesStore';
import ApiClient from '../libs/api/ApiClient';
import ImagesApi from '../libs/api/ImagesApi';

export default class RootStore {
    images: ImageStore;

    constructor(imagesApi: ImagesApi) {
        this.images = new ImageStore(imagesApi);
    }
}
