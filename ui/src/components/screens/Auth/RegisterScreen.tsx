import * as React from 'react';
import { RouteComponentProps, withRouter } from 'react-router';
import { action } from 'mobx';
import { inject, observer } from 'mobx-react';

import { ERR_CREDENTIALS, ERR_USERNAMETAKEN, ERR_EMAILTAKEN } from '../../../libs/api/ApiError';
import AuthStore from '../../../state/AuthStore';
import AuthScreen from './AuthScreen';
import Button from '../../layout/Button';
import { observable, computed } from 'mobx';

interface Props extends RouteComponentProps<{}> {
    auth?: AuthStore;
}

@withRouter
@inject('auth')
@observer
export default class RegisterScreen extends React.Component<Props> {
    @observable username: string = '';
    @observable email: string = '';
    @observable password: string = '';

    @action.bound
    updateValue(field: string, value: string) {
        (this as any)[field] = value.trim();
    }

    @computed
    get canSubmit() {
        return this.username !== '' && this.email !== '' && this.password !== '';
    }


    bindInput(field: string): (e: any) => any {
        return (e: any) => {
            this.updateValue(field, e.target.value);
        };
    }

    submit = (e: any) => {
        e.preventDefault();
        this.props.auth.register.run(this.username, this.email, this.password)
            .finally(this.afterRegisterAttempt)
    }

    afterRegisterAttempt = () => {
        if (!this.props.auth.register.err) {
            this.props.history.replace('/register/waiting');
        }
    }

    gotoSignin = () => {
        this.props.history.push('/login');
    }

    renderRegisterErr() {
        let msg;

        switch (this.props.auth.register.err.kind) {
            case ERR_USERNAMETAKEN:
                msg = `This username is already taken.`;
                break;
            case ERR_EMAILTAKEN:
                msg = `This email is registered already`;
                break;
            default:
                msg = `Authentication failed unexpectedly.`;
        }

        return msg;
    }

    render() {
        return (
            <AuthScreen>
                <h1>Register</h1>
                <form className="form" onSubmit={this.submit}>
                    <div className="entry">
                        <label>Username:</label>
                        <input tabIndex={0} type="text" placeholder="johndoe" onChange={this.bindInput('username')} />
                    </div>
                    <div className="entry">
                        <label>Email:</label>
                        <input tabIndex={0} type="text" placeholder="johndoe@mobingi.com" onChange={this.bindInput('email')} />
                    </div>
                    <div className="entry">
                        <label>Password:</label>
                        <input tabIndex={0} type="password" onChange={this.bindInput('password')} />
                    </div>
                    <div className="flex flex-h flex-grow" style={{ justifyContent: 'space-between', alignItems: 'center' }}>
                        <span className="err-msg">
                            {
                                this.props.auth.register.err && this.canSubmit &&
                                this.renderRegisterErr()
                            }
                        </span>
                        <Button text="Register" size="big" disabled={!this.canSubmit} onClick={this.submit} />
                    </div>
                </form>
                <div className="alternative-section flex flex-h flex-grow" style={{ justifyContent: 'space-between', alignItems: 'center' }}>
                    <span>Already have an account?</span>
                    <Button outline text="Sign in" size="big" onClick={this.gotoSignin} />
                </div>
            </AuthScreen>
        );
    }
}
