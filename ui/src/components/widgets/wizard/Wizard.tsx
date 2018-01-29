import * as React from 'react';
import './Wizard.scss';

interface Props {
    children?: React.ReactNode;
    step: number;
}

export default class Wizard extends React.PureComponent<Props> {
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
