import * as React from 'react';
import { inject, observer } from 'mobx-react';
import { withRouter, RouteComponentProps } from 'react-router-dom';

import { ERR_CREDENTIALS } from '../../libs/api/ApiError';
import AuthStore from '../../state/AuthStore';

import Alert from '../widgets/Alert';

interface Props extends RouteComponentProps<{}> {
    auth?: AuthStore;
}

@withRouter
@inject('auth')
@observer
export default class LoginScreen extends React.Component<Props> {
    usernameIn: HTMLInputElement;
    passwordIn: HTMLInputElement;

    submit = (e: any) => {
        e.preventDefault();
        this.props.auth.login.run(this.usernameIn.value, this.passwordIn.value)
            .finally(this.afterLoginAttempt)
    }

    afterLoginAttempt = () => {
        if (!this.props.auth.login.err) {
            this.props.history.replace('/images');
        }
    }

    renderLoginErr() {
        let msg;

        switch (this.props.auth.login.err.kind) {
            case ERR_CREDENTIALS:
                msg = `Username or password is wrong. Please try again`;
                break;
            default:
                msg = `Failed to authenticate for an unknown reason, we're sorry. Please try again later`;
        }

        return (
            <Alert message={msg} />
        );
    }

    render() {
        return (
            <div className="login-screen">
                {this.props.auth.login.err && this.renderLoginErr()}
                <form onSubmit={this.submit}>
                    <div>
                        <label>Username:</label>
                        <input type="text" ref={e => this.usernameIn = e} />
                    </div>
                    <div>
                        <label>Password:</label>
                        <input type="password" ref={e => this.passwordIn = e} />
                    </div>
                    <button type="submit">Login</button>
                </form>
            </div>
        );
    }
}
