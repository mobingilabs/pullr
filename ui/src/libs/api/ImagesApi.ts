import * as Promise from 'bluebird';
import ApiClient from "./ApiClient";
import Image, { IImage } from '../../state/models/Image';

export default class ImagesApi {
    apiClient: ApiClient;

    constructor(apiClient: ApiClient) {
        this.apiClient = apiClient;
    }

    get(key: string): Promise<IImage> {
        return this.apiClient.json('get', `/images/${key}`);
    }

    getImages(since?: Date): Promise<IImage[]> {
        const params: { [key: string]: string } = {};
        if (since) {
            params['since'] = (since.getTime() / 1000).toFixed(0)
        }

        return this.apiClient.json('get', '/images', params);
    }

    create(image: Image): Promise<Image> {
        const body = JSON.stringify(image.toJS());
        return this.apiClient.json('post', '/images', null, { body })
            .tap(({ key }) => image.key = key)
            .then(() => image);
    }

    update(key: string, image: Image): Promise<{ key: string }> {
        const body = JSON.stringify(image.toJS());
        return this.apiClient.json('post', `/images/${key}`, null, { body });
    }
}
