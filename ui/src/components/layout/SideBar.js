import React from 'react';
import PropTypes from 'prop-types';

import Logo from './Logo';
import MenuItem from './MenuItem';
import Icons from './Icons';
import './SideBar.scss';

export default class SideBar extends React.PureComponent {
    static propTypes = {
        wide: PropTypes.bool
    }

    static defaultProps = {
        wide: true
    }

    handleLogout = () => {
        console.log('Logout');
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
                    <MenuItem icon={Icons.logout} text="Logout" onClick={this.handleLogout} />
                </ul>
            </div>
        );
    }
}
