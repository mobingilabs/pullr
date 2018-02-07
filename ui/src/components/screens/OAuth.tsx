import * as React from 'react';
import * as QS from 'query-string';
import Screen from '../layout/Screen';
import Header from '../layout/Header';
import Button from '../layout/Button';

export default class OAuthScreen extends React.PureComponent<{}> {
    params: { [key: string]: string };

    constructor() {
        super({});

        this.params = QS.parse(location.search);
    }

    renderErr(kind: string) {
        if (kind === 'ERR_INTERNAL') {
            return 'An unexpected error happened on our side, please try again later.';
        } else {
            return `The service you're trying to connect failed to grant access to us. Please connect with your service provider.`;
        }
    }

    renderSuccess() {
        return 'Authentication was successful, please wait for a moment...';
    }

    close = () => {
        window.close();
    }

    componentDidMount() {
        const { err_kind, err_status, provider } = this.params;
        const failed = err_kind || err_status;

        const opener = window.opener ? window.opener : window.dialogArguments;

        if (failed) {
            opener.postMessage('ERROR', location.origin);
        } else {
            opener.postMessage('OAUTH_SUCCESS', location.origin);
            setTimeout(() => window.close(), 1000);
        }
    }

    render() {
        const { err_kind, err_status, provider } = this.params;
        if (err_kind && err_status) {
            return this.renderErr(err_kind);
        }

        return (
            <Screen>
                <Header title={`Login with ${provider}`} />
                <div>
                    {err_kind ? this.renderErr(err_kind) : this.renderSuccess()}
                </div>
                <div>
                    <Button size="big" onClick={this.close} text="Close" />
                </div>
            </Screen>
        )
    }
}
