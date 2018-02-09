import * as React from 'react';

import './Alert.scss';

interface Props {
    message: string;
    type?: 'info' | 'error' | 'warning' | 'success';
}

export default class Alert extends React.Component<Props> {
    public static defaultProps: Partial<Props> = {
        type: 'error'
    };

    render() {
        const classes = `alert alert-${this.props.type}`;
        return <div className={classes}> {this.props.message} </div>;
    }
}
