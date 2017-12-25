import React from 'react';

import './SectionTitle.scss';

export default class SectionTitle extends React.PureComponent {
    render () {
        return <div className="sectiontitle">{ this.props.children }</div>;
    }
}
