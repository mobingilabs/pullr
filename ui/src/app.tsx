import * as React from 'react';
import * as Promise from 'bluebird';
import { Route, Switch, withRouter, RouteComponentProps } from 'react-router-dom';
import { autorun, whyRun, IReactionDisposer } from 'mobx';
import { inject, observer } from 'mobx-react';

import RootStore from './state/RootStore';
import SideBar from './components/layout/SideBar';
import InitScreen from './components/screens/Init';
import ImagesScreen from './components/screens/Images';
import AddImageScreen from './components/screens/AddImage';
import EditImageScreen from './components/screens/EditImage';
import HistoryScreen from './components/screens/History';
import LoginScreen from './components/screens/Login';
import OAuthScreen from './components/screens/OAuth';

import ApiError from './libs/api/ApiError';

interface Props extends RouteComponentProps<{}> {
    store?: RootStore;
}

@withRouter
@inject('store')
@observer
export default class App extends React.Component<Props> {
    componentWillMount() {
        this.props.store.init.run().finally(this.checkUser).done();
    }

    checkUser = () => {
        const path = this.props.history.location.pathname;
        console.log(`PATH: ${path}`);
        if (!this.props.store.init.err && (path === '/login' || path === '/')) {
            this.props.history.push('/images');
        }
    }

    render() {
        if (this.props.store.init.inProgress) {
            return <InitScreen />;
        }

        if (!this.props.store.auth.loggedIn) {
            return (
                <Switch>
                    <Route path="/login" exact component={LoginScreen} />
                </Switch>
            );
        }

        return [
            <SideBar key="sidebar" />,
            <Switch key="switch">
                <Route path="/oauth" exact component={OAuthScreen} />
                <Route path="/images/add" exact component={AddImageScreen} />
                <Route path="/images/:imageKey/edit" exact component={EditImageScreen} />
                <Route path="/images" component={ImagesScreen} />
                <Route path="/history" component={HistoryScreen} />
            </Switch>
        ];
    }
}
