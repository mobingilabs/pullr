type Fetcher = (url: string, options: RequestInit) => Promise<Response>;

export default class ApiClient {
    fetcher: Fetcher;
    authToken: string;
    baseUrl: string;

    constructor(fetcher: Fetcher, authToken: string, baseUrl: string) {
        this.fetcher = fetcher;
        this.authToken = authToken;
        this.baseUrl = baseUrl;
    }

    private prepareUrl(path: string, params: {[key: string]: string}): string {
        const encodedParams = Object.keys(params)
            .map(name => `${name}=${encodeURIComponent(params[name])}`)
            .join('&');

        return `${this.baseUrl}/${path}${encodedParams.length > 0 ? '?' : ''}${encodedParams}`;
    }

    doRequest(path: string, params: {[key: string]: string}, options: RequestInit): Promise<any> {
        const url = this.prepareUrl(path, params)
        // TODO: put token to options
        return fetch(url, options);
    }
}
