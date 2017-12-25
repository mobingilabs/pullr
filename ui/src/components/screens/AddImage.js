import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';

import Screen from '../layout/Screen';
import Header from '../layout/Header';
import Icons from '../layout/Icons';

import Wizard from '../widgets/wizard/Wizard';
import WizardStep from '../widgets/wizard/WizardStep';
import ChooseRepository from './AddImageSteps/ChooseRepository';
import ConfigureImage from './AddImageSteps/ConfigureImage';
import SelectSourceProvider from './AddImageSteps/SelectSourceProvider';

import * as actions from '../../state/addImage/actions';

export class AddImageScreen extends React.PureComponent {
    cancel = () => {
        this.props.onReset();
        this.props.history.push('/images');
    }

    createImage = () => {
        this.props.onCreateImage({
            provider: this.props.provider,
            organisation: this.props.organisation,
            repository: this.props.repository,
            dockerfilePath: this.props.config.dockerfilePath,
            tags: this.props.config.tags
        });
        this.props.history.push('/images');
    }

    render () {
        return (
            <Screen className="screen-addimage">
                <Header title="ADD IMAGE" subTitle="Source Provider ..." actions={[]} />
                <Wizard step={ this.props.step }>
                    <WizardStep 
                        title="Select a source provider" 
                        component={SelectSourceProvider} 
                        props={{ 
                            onSelectProvider: this.props.onSelectProvider 
                        }} />

                    <WizardStep 
                        title="Choose a repository from Github" 
                        component={ChooseRepository}
                        props={{
                            selectedOrganisation: this.props.selectedOrganisation,
                            organisations: this.props.organisations,
                            repositories: this.props.repositories,
                            onSelectOrganisation: this.props.onSelectOrganisation,
                            onSelectRepository: this.props.onSelectRepository
                        }}/>
                    <WizardStep 
                        title="Configure image" 
                        component={ConfigureImage}
                        props={{
                            selectedOrganisation: this.props.selectedOrganisation,
                            selectedRepository: this.props.selectedRepository,
                            organisations: this.props.organisations,
                            repositories: this.props.repositories,
                            config: this.props.config,
                            onCancel: this.cancel,
                            onSelectOrganisation: this.props.onSelectOrganisation,
                            onSelectRepository: this.props.onSelectRepository,
                            onCreateImage: this.createImage,
                            onUpdateConfig: this.props.onUpdateConfig
                        }}/>
                </Wizard>
            </Screen>
        );
    }
}

function mapStateToProps(state) {
    return {
        step: state.addImage.step,
        provider: state.addImage.provider,
        selectedOrganisation: state.addImage.organisation,
        selectedRepository: state.addImage.repository,
        organisations: state.addImage.organisations,
        repositories: state.addImage.repositories,
        config: state.addImage.config
    }
}

function mapDispatchToProps(dispatch) {
    return {
        onReset: () => dispatch(actions.reset()),
        onSelectProvider: (provider) => dispatch(actions.setProvider(provider)),
        onSelectOrganisation: (organisation) => dispatch(actions.setOrganisation(organisation)),
        onSelectRepository: (repository) => dispatch(actions.setRepository(repository)),
        onUpdateConfig: (config) => dispatch(actions.updateConfig(config)),
        onCreateImage: (image) => dispatch(actions.createImage())
    }
}

export default connect(mapStateToProps, mapDispatchToProps)(AddImageScreen);
