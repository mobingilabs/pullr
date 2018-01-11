import ApiClient from "./ApiClient";
import Pagination from '../../state/models/Pagination';
import Image from '../../state/models/Image';

interface LoadImagesResponse {
    images: Array<Image>;
    pagination: Pagination;
}

export default class ImagesApi {
    apiClient: ApiClient;

    constructor(apiClient: ApiClient) {
        this.apiClient = apiClient;
    }

    loadImages(): Promise<LoadImagesResponse> {
        //TODO: ImagesApi::loadImages do real http request

        // FIXME: Delete me
        let images = [];
        try { images = JSON.parse(localStorage.getItem('images')) }
        catch (e) {}

        return Promise.resolve({
            images: [],
            pagination: new Pagination({
                itemsPerPage: 10,
                currentPage: 0,
                totalItems: 0
            })
        });
    }

    loadImagesPage(pageNumber: number): Promise<Array<Image>> {
        // TODO: ImagesApi::loadImagesPage do real http request

        // FIXME: Delete me
        let images = [];
        if (pageNumber > 0) {
            try { images = JSON.parse(localStorage.getItem('images')) }
            catch (e) { }
        }

        return Promise.resolve([]);
    }
}
