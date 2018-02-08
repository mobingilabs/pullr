import * as Promise from 'bluebird';
import ApiClient from "./ApiClient";
import Image, { IImage } from '../../state/models/Image';
import Pagination from '../../state/models/Pagination';
import ListOptions from '../../state/models/ListOptions';

export default class ImagesApi {
    apiClient: ApiClient;

    constructor(apiClient: ApiClient) {
        this.apiClient = apiClient;
    }

    get(key: string): Promise<IImage> {
        return this.apiClient.json('get', `/images/${key}`);
    }

    getImages(opts: ListOptions, since?: Date): Promise<{ images: IImage[], pagination: Partial<Pagination> }> {
        let params: { [key: string]: string } = {};
        if (since) {
            params['since'] = (since.getTime() / 1000).toFixed(0)
        }

        params = Object.assign({}, params, opts.toJS())

        return this.apiClient.json('get', '/images', params)
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
