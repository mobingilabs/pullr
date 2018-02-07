import * as React from 'react';
import { withRouter, RouteComponentProps } from "react-router-dom";
import { inject } from 'mobx-react';

import Logo from './Logo';
import MenuItem from './MenuItem';
import Icons from './Icons';
import './SideBar.scss';
import RootStore from '../../state/RootStore';

interface Props extends RouteComponentProps<{}> {
    wide?: boolean;
    store?: RootStore;
}

@withRouter
@inject('store')
export default class SideBar extends React.Component<Props> {
    static defaultProps = {
        wide: true
    }

    render() {
        const classes = ['sidebar'].concat(
            this.props.wide ? ['wide'] : []
        ).join(" ");

        return (
            <div className={classes}>
                <Logo />
                <ul className="main-navigation">
                    <MenuItem icon={Icons.images} path="/images" text="Images" />
                    <MenuItem icon={Icons.history} path="/history" text="Build History" />
                </ul>
                <ul className="secondary-navigation">
                    <MenuItem icon={Icons.settings} path="/settings" text="Settings" />
                    <MenuItem icon={Icons.logout} text="Logout" onClick={this.props.store.auth.logout.bind(null, this.props.history)} />
                </ul>
            </div>
        );
    }
}
