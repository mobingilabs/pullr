import * as React from 'react';
import * as PropTypes from 'prop-types';
import { observable, action, computed, IObservableArray, IReactionDisposer, reaction } from 'mobx';
import { observer, inject } from 'mobx-react';

import SourceStore from '../../../state/SourceStore';

import Button from '../../layout/Button';
import Icons from '../../layout/Icons';
import ImageEditor from '../../widgets/ImageEditor';
import Image from '../../../state/models/Image';
import WizardStep from '../../widgets/wizard/WizardStep';
import AsyncCmd from '../../../libs/asyncCmd';

interface Props {
    image: Image;
    sources?: SourceStore;
    onCancel: () => any;
    onFinish: () => any;
    saveCmd: AsyncCmd<void>;
}

@inject('sources')
@observer
export default class ConfigureImage extends React.Component<Props> {
    componentWillMount() {
        this.props.sources.selectProvider(this.props.image.repository.provider);
    }

    handleCreate = () => {
        if (this.props.saveCmd.inProgress) {
            return;
        }

        this.props.saveCmd.run().then(this.props.onFinish).done();
    }

    render() {
        const icon = this.props.saveCmd.inProgress ? Icons.loadingSpinner : null;
        return (
            <WizardStep title="Configure image">
                <div className="step-configure-image">
                    <ImageEditor image={this.props.image} sourceRepositories={this.props.sources.currentRepositories} sourceOwners={this.props.sources.currentOrganisations}>
                        <div className="form-actions">
                            <Button text="Cancel" size="big" secondary onClick={this.props.onCancel} />
                            <Button text="Create Image" icon={icon} size="big" onClick={this.props.onFinish} disabled={this.props.image.name.trim().length === 0} />
                        </div>
                    </ImageEditor>
                </div>
            </WizardStep>
        );
    }
}
