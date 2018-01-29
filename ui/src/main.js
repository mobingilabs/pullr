import React from 'react';
import ReactDOM from 'react-dom';
import { autorun } from 'mobx';
import { Provider } from 'mobx-react';
import { BrowserRouter as Router } from 'react-router-dom';
import DevTools from 'mobx-react-devtools';

import 'normalize.css/normalize.css';
import 'flexboxgrid/css/flexboxgrid.css';
import './styles/main.scss';

import App from './app';
import ModalBoundary from './components/layout/ModalBoundry';

import Image from './state/models/Image';
import RootStore from './state/RootStore';
import ApiClient from './libs/api/ApiClient';
import SourceApi from './libs/api/SourceApi';
import ImagesApi from './libs/api/ImagesApi';

window.React = React;

const fetcher = window.fetch;
const apiClient = new ApiClient(fetcher, '', '/');
const sourceApi = new SourceApi(apiClient);
const imagesApi = new ImagesApi(apiClient);

const rootStore = new RootStore(imagesApi);

const itemsCache = localStorage.getItem('images');
if (itemsCache) {
    try { rootStore.images.setImages(JSON.parse(itemsCache).map(image => new Image(image))); } 
    catch (e) { }
}
autorun(() => localStorage.setItem('images', JSON.stringify(rootStore.images.images)))

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
