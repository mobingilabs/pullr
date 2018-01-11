import * as React from 'react';
import * as PropTypes from 'prop-types';
import { observable, action, computed, IObservableArray, IReactionDisposer, reaction } from 'mobx';
import { observer, inject } from 'mobx-react';

import Button from '../../layout/Button';
import Icons from '../../layout/Icons';
import ImageEditor from '../../widgets/ImageEditor';
import Image from '../../../state/models/Image';
import WizardStep from '../../widgets/wizard/WizardStep';

import SourceApi from '../../../libs/api/SourceApi';

interface Props {
    image: Image;
    sourceApi?: SourceApi;
    onCancel: Function;
    onFinish: Function;
}

@inject('sourceApi')
@observer
export default class ConfigureImage extends React.Component<Props> {
    @observable sourceOwners: Array<string> = [];
    @observable sourceRepositories: Array<string> = [];

    private disposables: Array<IReactionDisposer> = [];

    componentWillMount() {
        this.props.sourceApi.getOwners(this.props.image.sourceProvider).then(this.setSourceOwners);
        this.props.sourceApi.getRepositories(this.props.image.sourceProvider, this.props.image.sourceOwner).then(this.setSourceRepositories);

        const sourceOwnerReaction = reaction(
            () => this.props.image.sourceOwner,
            (sourceOwner) => this.props.sourceApi.getRepositories(this.props.image.sourceProvider, sourceOwner).then(this.setSourceRepositories)
        );
    }

    componentWillUnmount() {
        this.disposables.forEach((disposer: IReactionDisposer) => disposer());
    }

    @action.bound
    setSourceOwners(sourceOwners: Array<string>) {
        this.sourceOwners = sourceOwners;
    }

    @action.bound
    setSourceRepositories(sourceRepositories: Array<string>) {
        this.sourceRepositories = sourceRepositories;
    }

    render() {
        return (
            <WizardStep title="Configure image">
                <div className="step-configure-image">
                    <ImageEditor image={this.props.image} sourceRepositories={this.sourceRepositories} sourceOwners={this.sourceOwners}>
                        <div className="form-actions">
                            <Button text="Cancel" size="big" secondary onClick={this.props.onCancel} />
                            <Button text="Create Image" size="big" onClick={this.props.onFinish} disabled={ this.props.image.name.trim().length === 0 } />
                        </div>
                    </ImageEditor>
                </div>
            </WizardStep>
        );
    }
}
