import * as React from 'react';
import { withRouter, RouteComponentProps } from 'react-router-dom';
import { observable, action, reaction, IObservableArray, IReactionDisposer } from 'mobx';
import { inject, observer } from 'mobx-react';

import Screen from '../layout/Screen';
import Header from '../layout/Header';
import Button from '../layout/Button';
import ImageEditor from "../widgets/ImageEditor";
import Notification from '../widgets/Notification';

import RootStore from '../../state/RootStore';
import SourceApi from '../../libs/api/SourceApi';
import Image from '../../state/models/Image';

interface RouteParams {
    imageName: string;
}

interface Props extends RouteComponentProps<RouteParams> {
    store?: RootStore;
    sourceApi?: SourceApi;
}

@withRouter
@inject('store')
@inject('sourceApi')
@observer
export default class EditImageScreen extends React.Component<Props> {
    @observable imageCopy: Image;
    @observable sourceOwners: Array<string> = [];
    @observable sourceRepositories: Array<string> = [];

    private disposables: Array<IReactionDisposer> = [];

    constructor(props: Props) {
        super(props);
        const image = this.props.store.images.findByName(this.props.match.params.imageName);

        if (image) {
            this.imageCopy = image.clone();

            this.props.sourceApi.getOwners(this.imageCopy.sourceProvider).then(this.setSourceOwners);
            this.props.sourceApi.getRepositories(this.imageCopy.sourceProvider, this.imageCopy.sourceOwner).then(this.setSourceRepositories);

            const sourceOwnerReaction = reaction(
                () => this.imageCopy.sourceOwner,
                (sourceOwner) => this.props.sourceApi.getRepositories(this.imageCopy.sourceProvider, sourceOwner).then(this.setSourceRepositories)
            );
        } else {
            requestAnimationFrame(() => {
                props.history.push('/images');
            });
        }
    }

    @action.bound
    setSourceOwners(sourceOwners: Array<string>) {
        this.sourceOwners = sourceOwners;
    }

    @action.bound
    setSourceRepositories(sourceRepositories: Array<string>) {
        this.sourceRepositories = sourceRepositories;
    }

    saveChanges = () => {
        this.props.store.images.updateImage(this.props.match.params.imageName, this.imageCopy);
    }

    render() {
        if (!this.imageCopy) {
            return null;
        }

        return (
            <Screen>
                <Header title={this.props.match.params.imageName} subTitle="Edit image properties..." back />
                <Notification id="images-update-image" />
                <div className="content">
                    <ImageEditor image={ this.imageCopy } sourceRepositories={this.sourceRepositories} sourceOwners={this.sourceOwners}>
                        <div className="form-actions">
                            <Button text="Save Changes" size="big" onClick={this.saveChanges} />
                        </div>
                    </ImageEditor>
                </div>
            </Screen>
        );
    }
}
