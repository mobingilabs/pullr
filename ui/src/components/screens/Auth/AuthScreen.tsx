import * as React from 'react';
import Logo from '../../layout/Logo';

import './AuthScreen.scss';

export default class AuthScreen extends React.Component<{}> {
    render() {
        return (
            <div className="authScreen">
                <div className="authScreen-inner">
                    <Logo white={false} />
                    <div className="card">
                        {this.props.children}
                    </div>
                </div>
            </div>
        );
    }
}
