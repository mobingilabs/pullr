import * as React from 'react';
import { withRouter, RouteComponentProps } from 'react-router';
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
export default class RegisterWaitScreen extends React.Component<Props> {
    render() {
        return (
            <AuthScreen>
                <h1>Awesome!</h1>
                <p>
                    We've send an email to you. Please confirm your email address by
                    clicking the activation link given in the mail.
                </p>
            </AuthScreen>
        );
    }
}
