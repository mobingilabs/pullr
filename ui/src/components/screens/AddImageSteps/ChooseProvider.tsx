import * as React from 'react';
import { observer, inject } from 'mobx-react';

import WizardStep from '../../widgets/wizard/WizardStep';
import Button from '../../layout/Button';
import Card from '../../widgets/Card';

import Image from '../../../state/models/Image';
import SourceApi from '../../../libs/api/SourceApi';

import './ChooseProvider.scss';

interface Props {
    image: Image;
    next: Function;
    sourceApi?: SourceApi;
}

@inject('sourceApi')
@observer
export default class ChooseProvider extends React.Component<Props> {
    selectGithub = () => {
        this.props.image.sourceProvider = 'github';
        this.props.sourceApi.loadOwners('github');
        this.props.next();
    }

    render() {
        return (
            <WizardStep title="Select a source provider">
                <div className="step-choose-provider">
                    <Card icon="github" title="Github" background="#CF2A63" dark>
                        <Button text="Link Account & Use" onClick={ this.selectGithub } />
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
