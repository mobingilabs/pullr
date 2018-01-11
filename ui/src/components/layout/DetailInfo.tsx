import * as React from 'react';
import './DetailInfo.scss';

interface Props {
    label: string;
}

export default class DetailInfo extends React.PureComponent<Props> {
    render() {
        return (
            <div className="detail-field">
                <div className="detail-label">{ this.props.label }</div>
                <div className="detail-value">{ this.props.children }</div>
            </div>
        );
    }
}
