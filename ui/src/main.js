import 'babel-polyfill';
import React from 'react';
import ReactDOM from 'react-dom';
import { BrowserRouter as Router } from 'react-router-dom';
import { createStore } from 'redux';
import { Provider } from 'react-redux';

import 'normalize.css/normalize.css';
import 'flexboxgrid/css/flexboxgrid.css';
import './styles/main.scss';

import App from './app';
import ModalBoundary from './components/layout/ModalBoundry';
import InitialState from './state/initial';
import reducer from './state/reducer';

window.React = React;

const store = createStore(reducer, InitialState);
window.store = store;

if (module.hot) {
    // Enable Webpack hot module replacement for reducers
    module.hot.accept('./state/reducer', () => {
        const nextRootReducer = require('./state/reducer');
        store.replaceReducer(nextRootReducer);
    });
}

ReactDOM.render(
    <Provider store={store}>
        <Router>
            <ModalBoundary className="flex flex-h flex-grow">
                <App />
            </ModalBoundary>
        </Router>
    </Provider>,
    document.getElementById('root')
);
