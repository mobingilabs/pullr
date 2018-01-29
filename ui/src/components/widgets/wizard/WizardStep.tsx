import * as React from 'react';
import * as PropTypes from 'prop-types';

import SectionTitle from '../../layout/SectionTitle';

interface Props {
    title: string;
}

export default class WizardStep extends React.PureComponent<Props> {
    render () {
        return (
            <div className="wizard-step">
                <SectionTitle>{ this.props.title }</SectionTitle>
                <div className="wizard-step-content">
                    { this.props.children }
                </div>
            </div>
        )
    }
}
