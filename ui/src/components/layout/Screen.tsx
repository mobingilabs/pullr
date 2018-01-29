import * as React from 'react';
import './Screen.scss';

interface Props {
    children: React.ReactNode,
    className?: string;
}

export default class Screen extends React.PureComponent<Props> {
    render () {
        const classes = ['screen'].concat(this.props.className ? [this.props.className] : []).join(' ');

        return (
            <div className={ classes }>{this.props.children}</div>
        );
    }
}
