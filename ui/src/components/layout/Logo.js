import React from 'react';
import PropTypes from 'prop-types';

import './Logo.scss';

export default class Logo extends React.PureComponent {
    render () {
        const classes = ['logo'].concat([
            this.props.white ? 'white' : 'dark'
        ]).join(' ');

        return <div className={classes} />;
    }
}

Logo.propTypes = {
    white: PropTypes.bool
};

Logo.defaultProps = {
    white: true
};
