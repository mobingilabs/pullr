import React from 'react';
import './Screen.scss';

export default class Screen extends React.PureComponent {
    render () {
        const classes = ['screen'].concat(this.props.className ? [this.props.className] : []).join(' ');

        return (
            <div className={ classes }>{this.props.children}</div>
        );
    }
}
