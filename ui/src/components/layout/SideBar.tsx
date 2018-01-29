import * as React from 'react';
import { withRouter } from "react-router-dom";

import Logo from './Logo';
import MenuItem from './MenuItem';
import Icons from './Icons';
import './SideBar.scss';

interface Props {
    wide?: boolean;
}

@withRouter
export default class SideBar extends React.PureComponent<Props> {
    static defaultProps = {
        wide: true
    }

    logout = () => {

    }

    render () {
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
                    <MenuItem icon={Icons.logout} text="Logout" onClick={ this.logout } />
                </ul>
            </div>
        );
    }
}
