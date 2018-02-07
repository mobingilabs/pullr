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
        this.props.store.init.run().finally(this.checkUser);
    }

    checkUser = () => {
        if (this.props.store.init.inProgress) {
            return;
        }

        if (!this.props.store.auth.loggedIn) {
            this.props.history.replace('/login');
        } else if (this.props.history.location.pathname === '/') {
            this.props.history.push('/images');
        }
    }

    render() {
        if (this.props.store.init.inProgress) {
            return <InitScreen />;
        }

        const showAuth = !this.props.store.auth.loggedIn || this.props.history.location.pathname === '/oauth';

        return [
            showAuth ? null : <SideBar key="sidebar" />,
            <Switch key="switch">
                <Route path="/oauth" exact component={OAuthScreen} />
                <Route path="/login" exact component={LoginScreen} />

                <Route path="/images/add" exact component={AddImageScreen} />
                <Route path="/images/:imageName/edit" exact component={EditImageScreen} />
                <Route path="/images" component={ImagesScreen} />
                <Route path="/history" component={HistoryScreen} />
            </Switch>
        ];
    }
}
