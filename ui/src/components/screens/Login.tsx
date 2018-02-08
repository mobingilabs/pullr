import * as React from 'react';
import { inject, observer } from 'mobx-react';
import { withRouter, RouteComponentProps } from 'react-router-dom';
import AuthStore from '../../state/AuthStore';

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
            .then(this.redirectToImages)
            .catch(this.handleErr);
    }

    redirectToImages = () => {
        this.props.history.replace('/images');
    }

    handleErr = (err: any) => {
        console.error(err);
    }

    render() {
        return (
            <div className="login-screen">
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
