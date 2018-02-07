import * as React from 'react';
import { observer, inject } from 'mobx-react';

import WizardStep from '../../widgets/wizard/WizardStep';
import Button from '../../layout/Button';
import Card from '../../widgets/Card';
import Icons from '../../layout/Icons';

import Image from '../../../state/models/Image';
import RootStore from '../../../state/RootStore';
import ApiError from '../../../libs/api/ApiError';

import './ChooseProvider.scss';

interface Props {
    image: Image;
    next: Function;
    store?: RootStore;
}

interface State {
    err: ApiError;
}

@inject('store')
@observer
export default class ChooseProvider extends React.Component<Props, State> {
    componentDidMount() {
        if (!this.props.store.auth.user.tokens.has('github')) {
            this.props.store.auth.getLoginUrl('github').run();
        }

        // Trigger other oauth providers here
    }

    selectProvider = (provider: string) => {
        this.props.image.repository.provider = provider;
        this.props.next();
    }

    renderProviderButton(provider: string) {
        if (this.props.store.auth.user.tokens.has(provider)) {
            return <Button text="Select" onClick={this.selectProvider.bind(null, provider)} />;
        }

        const loadUrlCmd = this.props.store.auth.getLoginUrl(provider);
        if (loadUrlCmd.inProgress) {
            return <Button disabled text="" onClick={() => { }} icon={Icons.loadingSpinner} />
        }

        if (loadUrlCmd.err) {
            return 'Not available for now';
        }

        const tokens = this.props.store.auth.user.tokens;
        return <Button aslink popup text="Link Account & Use" href={loadUrlCmd.value} onClick={() => this.props.store.auth.oauthStart(provider)} />;
    }

    render() {
        return (
            <WizardStep title="Select a source provider">
                <div className="step-choose-provider">
                    <Card icon="github" title="Github" background="#CF2A63" dark>
                        {this.renderProviderButton('github')}
                    </Card>
                    <Card icon="bitbucket" title="Bitbucket" disabled>
                        Coming soon...
                    </Card>
                    <Card icon="gitlab" title="Gitlab" disabled>
                        Coming soon...
                    </Card>
                </div>
            </WizardStep>
        )
    }
}
