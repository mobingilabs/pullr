import * as React from 'react';

import './Alert.scss';

interface Props {
    message: string;
}

export default class Alert extends React.Component<Props> {
    render() {
        return <div className="alert"> {this.props.message} </div>;
    }
}
