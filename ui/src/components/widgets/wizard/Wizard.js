import React from 'react';
import PropTypes from 'prop-types';

import './Wizard.scss';

export default class Wizard extends React.PureComponent {
    static propTypes = {
        step: PropTypes.number.isRequired,
    }

    render() {
        const steps = React.Children.toArray(this.props.children);
        const currentStep = steps[this.props.step];

        return (
            <div className="wizard">
                { currentStep }
            </div>
        )
    }
}
