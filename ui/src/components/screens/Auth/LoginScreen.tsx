import * as React from 'react';
import { withRouter, RouteComponentProps } from 'react-router';
import { observable, computed, action } from 'mobx';
import { inject, observer } from 'mobx-react';

import { ERR_CREDENTIALS } from '../../../libs/api/ApiError';
import AuthStore from '../../../state/AuthStore';
import AuthScreen from './AuthScreen';
import Button from '../../layout/Button';

interface Props extends RouteComponentProps<{}> {
    auth?: AuthStore;
}

@withRouter
@inject('auth')
@observer
export default class LoginScreen extends React.Component<Props> {
    @observable username: string = '';
    @observable password: string = '';

    @computed
    get canSubmit(): boolean {
        return this.username !== '' && this.password !== '';
    }

    submit = (e: any) => {
        e.preventDefault();
        this.props.auth.login.run(this.username, this.password)
            .finally(this.afterLoginAttempt)
    }

    @action.bound
    updateValue(field: string, value: string) {
        (this as any)[field] = value;
    }

    bindInput(field: string): (e: any) => any {
        return (e: any) => {
            this.updateValue(field, e.target.value);
        };
    }

    afterLoginAttempt = () => {
        if (!this.props.auth.login.err) {
            this.props.history.replace('/images');
        }
    }

    gotoRegister = () => {
        this.props.history.push('/register');
    }

    renderLoginErr() {
        let msg;

        switch (this.props.auth.login.err.kind) {
            case ERR_CREDENTIALS:
                msg = `Username or password is wrong.`;
                break;
            default:
                msg = `Authentication failed unexpectedly.`;
        }

        return msg;
    }

    render() {
        return (
            <AuthScreen>
                <h1>Login</h1>
                <form className="form" onSubmit={this.submit}>
                    <div className="entry">
                        <label>Username or email:</label>
                        <input tabIndex={0} type="text" placeholder="johndoe@mobingi.com" onChange={this.bindInput('username')} />
                    </div>
                    <div className="entry">
                        <label>Password: <a tabIndex={1} href="#">Forgot your password?</a></label>
                        <input tabIndex={0} type="password" onChange={this.bindInput('password')} />
                    </div>
                    <div className="flex flex-h flex-grow" style={{ justifyContent: 'space-between', alignItems: 'center' }}>
                        <span className="err-msg">{this.props.auth.login.err && this.renderLoginErr()}</span>
                        <Button disabled={!this.canSubmit} text="Sign in" size="big" onClick={this.submit} />
                    </div>
                </form>
                <div className="alternative-section flex flex-h flex-grow" style={{ justifyContent: 'space-between', alignItems: 'center' }}>
                    <span>Don't you have an account yet?</span>
                    <Button outline text="Register" size="big" onClick={this.gotoRegister} />
                </div>
            </AuthScreen>
        );
    }
}
