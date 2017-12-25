import React from 'react';
import PropTypes from 'prop-types';

import SectionTitle from '../../layout/SectionTitle';

export default class WizardStep extends React.PureComponent {
    static propTypes = {
        title: PropTypes.string.isRequired,
        component: PropTypes.func.isRequired,
        props: PropTypes.object
    }

    static defaultProps = {
        props: {}
    }

    render () {
        const Component = this.props.component;
        return (
            <div className="wizard-step">
                <SectionTitle>{ this.props.title }</SectionTitle>
                <div className="wizard-step-content">
                    <Component {...this.props.props} />
                </div>
            </div>
        )
    }
}
