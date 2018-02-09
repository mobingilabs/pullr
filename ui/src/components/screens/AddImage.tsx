import * as React from 'react';
import { observable, action } from 'mobx';
import { observer, inject } from 'mobx-react';
import { withRouter, RouteComponentProps } from "react-router-dom";

import ImagesStore from '../../state/ImagesStore';

import Screen from '../layout/Screen';
import Header from '../layout/Header';
import Icons from '../layout/Icons';
import Wizard from '../widgets/wizard/Wizard';
import ChooseProvider from './AddImageSteps/ChooseProvider';
import ChooseRepository from './AddImageSteps/ChooseRepository';
import ConfigureImage from './AddImageSteps/ConfigureImage';
import Image from '../../state/models/Image';

enum Steps {
    ChooseProvider = 0,
    ChooseRepository = 1,
    ConfigureImage = 2
}

interface Props extends RouteComponentProps<{}> {
    images?: ImagesStore;
}

@withRouter
@inject('images')
@observer
export default class AddImageScreen extends React.Component<Props> {
    @observable step: Steps;
    @observable newImage: Image;

    constructor(props: Props) {
        super(props);
        this.step = Steps.ChooseProvider;
        this.newImage = Image.create();
    }

    @action.bound
    showChooseRepository() {
        this.step = Steps.ChooseRepository;
    }

    @action.bound
    showConfigureImage() {
        this.step = Steps.ConfigureImage;
    }

    @action.bound
    onChangeDockerfilePath(e: any) {
        this.newImage.dockerfile_path = e.target.value;
    }

    cancel = () => {
        this.props.history.push('/images');
    }

    @action.bound
    afterCreate() {
        this.newImage = Image.create();
        this.props.history.push('/images');
    }

    @action.bound
    handleBack() {
        if (this.step > 0) {
            this.step -= 1;
        } else {
            this.props.history.goBack();
        }
    }

    render() {
        return (
            <Screen className="screen-addimage">
                <Header back={true} onBack={this.handleBack} title="ADD IMAGE" subTitle="Source Provider ..." />
                <Wizard step={this.step}>
                    <ChooseProvider image={this.newImage} next={this.showChooseRepository} />
                    <ChooseRepository image={this.newImage} next={this.showConfigureImage} />
                    <ConfigureImage image={this.newImage} saveCmd={this.props.images.saveImage} onFinish={this.afterCreate} onCancel={this.cancel} />
                </Wizard>
            </Screen>
        );
    }
}
