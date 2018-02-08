import * as Promise from 'bluebird';
import ApiError from './ApiError';

type Fetcher = (url: string, options: RequestInit) => Promise<Response>;

export default class ApiClient {
    fetcher: Fetcher;
    baseUrl: string;
    baseReqOpts: Partial<RequestInit>;
    authToken: string;
    refreshToken: string;

    constructor(fetcher: Fetcher, baseUrl: string, authToken: string, refreshToken: string) {
        this.fetcher = fetcher;
        this.baseUrl = baseUrl;
        this.authToken = authToken;
        this.refreshToken = refreshToken;

        this.baseReqOpts = {
            headers: {
                'X-Requested-With': 'XMLHttpRequest',
                'Content-Type': 'application/json',
            },
        };
    }

    setTokens(authToken: string, refreshToken: string) {
        this.authToken = authToken;
        this.refreshToken = refreshToken;
    }

    json(method: string, path: string, params?: { [key: string]: string }, options?: RequestInit): Promise<any> {
        return this.request(method, path, params, options).then(res => res.json());
    }

    request(method: string, path: string, params?: { [key: string]: string }, options?: RequestInit): Promise<Response> {
        const url = this.prepareUrl(path, params);
        const newOpts = this.prepareReq(method, options);
        return Promise.resolve(fetch(url, newOpts)).then(this.checkStatus).then(this.afterReq);
    }

    encodeParams(params: { [key: string]: any }): string {
        return Object.keys(params || {})
            .map(name => `${name}=${encodeURIComponent(params[name].toString())}`)
            .join('&');
    }

    private hasTokens(): boolean {
        return this.authToken && this.authToken.length > 0 && this.refreshToken && this.refreshToken.length > 0;
    }

    private prepareReq(method: string, reqOpts: Partial<RequestInit>): Partial<RequestInit> {
        const newOpts = Object.assign({ method }, this.baseReqOpts, reqOpts || {});
        const headers = Object.assign({}, newOpts.headers);

        if (this.hasTokens()) {
            headers['Authorization'] = `Bearer ${this.authToken}`;
            headers['X-Refresh-Token'] = this.refreshToken;
        }

        newOpts.headers = headers;
        return newOpts;
    }

    private prepareUrl(path: string, params: { [key: string]: string }): string {
        const encodedParams = this.encodeParams(params);
        return `${this.baseUrl}${path}${encodedParams.length > 0 ? '?' : ''}${encodedParams}`;
    }

    private afterReq = (res: Response): Response => {
        const authToken = res.headers.get('X-Auth-Token') || '';
        const refreshToken = res.headers.get('X-Refresh-Token') || '';

        if (authToken.length > 0 && refreshToken.length > 0) {
            this.authToken = authToken;
            this.refreshToken = refreshToken;
            window.localStorage['authToken'] = this.authToken;
            window.localStorage['refreshToken'] = this.refreshToken;
        }

        return res;
    }

    private checkStatus = (res: Response): Promise<Response> => {
        if (res.status >= 400 && res.status < 600) {
            return Promise.resolve(res.json()).then((err: any) => {
                throw new ApiError(err.kind, res.status);
            });
        }

        return Promise.resolve(res);
    }

}
