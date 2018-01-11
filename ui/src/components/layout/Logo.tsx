import * as React from 'react';

import './Logo.scss';

interface Props {
    white?: boolean;
}

export default class Logo extends React.PureComponent<Props> {
    static defaultProps = {
        white: true
    }

    render () {
        const classes = ['logo'].concat([
            this.props.white ? 'white' : 'dark'
        ]).join(' ');

        return <div className={classes} />;
    }
}
