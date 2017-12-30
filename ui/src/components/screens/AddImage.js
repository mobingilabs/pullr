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
        this.props.onCreateImage(this.props.image);
        this.props.onReset();
        this.props.history.push(`/images/${this.props.image.name}`);
    }

    render () {
        return (
            <Screen className="screen-addimage">
                <Header title="ADD IMAGE" subTitle="Source Provider ..." />
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
                            selectedOrganisation: this.props.image.organisation,
                            organisations: this.props.organisations,
                            repositories: this.props.repositories,
                            onSelectOrganisation: this.props.onSelectOrganisation,
                            onSelectRepository: this.props.onSelectRepository
                        }}/>
                    <WizardStep 
                        title="Configure image" 
                        component={ConfigureImage}
                        props={{
                            image: this.props.image,
                            organisations: this.props.organisations,
                            repositories: this.props.repositories[this.props.image.organisation],
                            onCancel: this.cancel,
                            onCreateImage: this.createImage,
                            onUpdateImage: this.props.onUpdateImage
                        }}/>
                </Wizard>
            </Screen>
        );
    }
}

function mapStateToProps(state) {
    return {
        step: state.addImage.step,
        organisations: state.addImage.organisations,
        repositories: state.addImage.repositories,
        image: state.addImage.image
    }
}

function mapDispatchToProps(dispatch) {
    return {
        onReset: () => dispatch(actions.reset()),
        onSelectProvider: (provider) => dispatch(actions.setProvider(provider)),
        onSelectOrganisation: (organisation) => dispatch(actions.setOrganisation(organisation)),
        onSelectRepository: (repository) => dispatch(actions.setRepository(repository)),
        onUpdateImage: (config) => dispatch(actions.updateImage(config)),
        onCreateImage: (image) => dispatch(actions.createImage(image))
    }
}

export default connect(mapStateToProps, mapDispatchToProps)(AddImageScreen);
