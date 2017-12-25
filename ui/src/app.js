import React from 'react';
import { BrowserRouter as Router, Route, Switch } from 'react-router-dom';

import SideBar from './components/layout/SideBar';

import ImagesScreen from './components/screens/Images';
import AddImageScreen from './components/screens/AddImage';
import HistoryScreen from './components/screens/History';

export default class App extends React.Component {
    render () {
        return [
            <SideBar key="sidebar" />,
            <Switch key="routes">
                <Route path="/images/add" exact component={AddImageScreen} />
                <Route path="/images" component={ImagesScreen} />
                <Route path="/history" component={HistoryScreen} />
            </Switch>
        ];
    }
}
