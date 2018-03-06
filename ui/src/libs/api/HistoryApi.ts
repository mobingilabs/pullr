import ApiClient from './ApiClient';
import ListOptions from '../../state/models/ListOptions';
import Image from '../../state/models/Image';
import Pagination from '../../state/models/Pagination';

interface ImageHistoryResponse {
    images: Image[];
    pagination: Pagination;
}

export default class HistoryApi {
    private apiClient: ApiClient;

    constructor(apiClient: ApiClient) {
        this.apiClient = apiClient;
    }

    historyOfImages(opts: ListOptions, since?: Date): Promise<{}> {

    }
}
