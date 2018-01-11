import * as React from 'react';
import { observable, action } from 'mobx';
import { observer, inject } from 'mobx-react';
import { withRouter, RouteComponentProps } from "react-router-dom";

import Screen from '../layout/Screen';
import Header from '../layout/Header';
import Icons from '../layout/Icons';

import Wizard from '../widgets/wizard/Wizard';
import ChooseProvider from './AddImageSteps/ChooseProvider';
import ChooseRepository from './AddImageSteps/ChooseRepository';
import ConfigureImage from './AddImageSteps/ConfigureImage';
import RootStore from '../../state/RootStore';
import Image from '../../state/models/Image';

enum Steps {
    ChooseProvider = 0,
    ChooseRepository = 1,
    ConfigureImage = 2
}

interface Props extends RouteComponentProps<{}> {
    store?: RootStore;
}

@withRouter
@inject('store')
@observer
export default class AddImageScreen extends React.Component<Props> {
    @observable step: Steps;
    @observable newImage: Image;

    constructor(props: Props) {
        super(props);
        this.step = Steps.ChooseProvider;
        this.newImage = Image.create();
    }

    @action showChooseRepository = () => {
        this.step = Steps.ChooseRepository;
    }
    
    @action showConfigureImage = () => {
        this.step = Steps.ConfigureImage;
    }

    @action onChangeDockerfilePath = (e: any) => {
        this.newImage.dockerfilePath = e.target.value;
    }

    cancel = () => {
        this.props.history.push('/images');
    }

    @action saveImage = () => {
        this.props.store.images.saveImage(this.newImage);
        this.newImage = Image.create();
        this.props.history.push('/images');
    }

    render () {
        return (
            <Screen className="screen-addimage">
                <Header title="ADD IMAGE" subTitle="Source Provider ..." />
                <Wizard step={ this.step }>
                    <ChooseProvider image={this.newImage} next={this.showChooseRepository} />
                    <ChooseRepository image={this.newImage} next={this.showConfigureImage} />
                    <ConfigureImage image={this.newImage} onFinish={this.saveImage} onCancel={this.cancel} />
                </Wizard>
            </Screen>
        );
    }
}
