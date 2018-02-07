import React from 'react';
import ReactDOM from 'react-dom';
import { Provider } from 'mobx-react';
import { BrowserRouter as Router } from 'react-router-dom';
import DevTools from 'mobx-react-devtools';

import 'normalize.css/normalize.css';
import 'flexboxgrid/css/flexboxgrid.css';
import './styles/main.scss';

import App from './app';
import ModalBoundary from './components/layout/ModalBoundry';

import RootStore from './state/RootStore';
import ApiClient from './libs/api/ApiClient';
import SourceApi from './libs/api/SourceApi';
import ImagesApi from './libs/api/ImagesApi';
import AuthApi from './libs/api/AuthApi';
import OAuthApi from './libs/api/OAuthApi';

window.React = React;

const authToken = window.localStorage['authToken'];
const refreshToken = window.localStorage['refreshToken'];

const fetcher = window.fetch;
const apiClient = new ApiClient(fetcher, '/api/v1', authToken, refreshToken);
const sourceApi = new SourceApi(apiClient);
const imagesApi = new ImagesApi(apiClient);
const authApi = new AuthApi(apiClient);
const oauthApi = new OAuthApi(apiClient);

const rootStore = new RootStore(imagesApi, authApi, oauthApi, sourceApi);

const modalRoot = document.getElementById('modal-root');
ReactDOM.render(
    <Provider store={rootStore} sourceApi={sourceApi}>
        <Router>
            <ModalBoundary modalRoot={modalRoot}>
                <App />
                <DevTools />
            </ModalBoundary>
        </Router>
    </Provider>,
    document.getElementById('root')
);
