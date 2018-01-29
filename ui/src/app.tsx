import * as React from 'react';
import { Route, Switch, withRouter } from 'react-router-dom';

import SideBar from './components/layout/SideBar';

import ImagesScreen from './components/screens/Images';
import AddImageScreen from './components/screens/AddImage';
import EditImageScreen from './components/screens/EditImage';
import HistoryScreen from './components/screens/History';

@withRouter
export default class App extends React.Component {
    render () {
        return [
            <SideBar key="sidebar" />,
            <Switch key="routes">
                <Route path="/images/add" exact component={AddImageScreen} />
                <Route path="/images/:imageName/edit" exact component={EditImageScreen} />
                <Route path="/images" component={ImagesScreen} />
                <Route path="/history" component={HistoryScreen} />
            </Switch>
        ];
    }
}
