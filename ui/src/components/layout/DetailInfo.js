import React from 'react';
import PropTypes from 'prop-types';

import './DetailInfo.scss';

export default class DetailInfo extends React.PureComponent {
    static propTypes = {
        label: PropTypes.string.isRequired
    }

    render() {
        return (
            <div className="detail-field">
                <div className="detail-label">{ this.props.label }</div>
                <div className="detail-value">{ this.props.children }</div>
            </div>
        );
    }
}