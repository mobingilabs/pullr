import 'babel-polyfill';
import React from 'react';
import ReactDOM from 'react-dom';
import { BrowserRouter as Router } from 'react-router-dom';
import { createStore, applyMiddleware, compose } from 'redux';
import { Provider } from 'react-redux';
import thunkMiddleware from "redux-thunk";

import 'normalize.css/normalize.css';
import 'flexboxgrid/css/flexboxgrid.css';
import './styles/main.scss';

import App from './app';
import ModalBoundary from './components/layout/ModalBoundry';
import InitialState from './state/initial';
import reducer from './state/reducer';

window.React = React;
const composeEnhancers = window.__REDUX_DEVTOOLS_EXTENSION_COMPOSE__ || compose;
const store = createStore(reducer, InitialState, composeEnhancers(applyMiddleware(thunkMiddleware)));
window.store = store;

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
