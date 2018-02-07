import * as Promise from 'bluebird';
import ApiClient from "./ApiClient";
import User from '../../state/models/User';
import ApiError from './ApiError';

export class CredentialsError extends Error { }
export class UsernameTakenError extends Error { }

export default class AuthApi {
    apiClient: ApiClient;

    constructor(apiClient: ApiClient) {
        this.apiClient = apiClient;
    }

    login = (username: string, password: string): Promise<User | ApiError> => {
        const body = JSON.stringify({ username, password });
        return this.apiClient.request('post', '/login', null, { body })
            .then(this.getProfile);
    }

    getProfile = (): Promise<User> => {
        return this.apiClient.json('get', '/profile')
            .then(user => new User(user));
    }

    logout = () => {
        this.apiClient.setTokens('', '');
        window.localStorage['authToken'] = '';
        window.localStorage['refreshToken'] = '';
    }
}
