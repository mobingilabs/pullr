import * as React from 'react';
import { withRouter, RouteComponentProps } from 'react-router-dom';
import { observable, action, reaction, IObservableArray, IReactionDisposer, transaction } from 'mobx';
import { inject, observer } from 'mobx-react';

import SourceStore from '../../state/SourceStore';
import ImageStore from '../../state/ImagesStore';

import Screen from '../layout/Screen';
import Header from '../layout/Header';
import Button from '../layout/Button';
import ImageEditor from "../widgets/ImageEditor";
import Image from '../../state/models/Image';
import Alert from '../widgets/Alert';
import Icons from '../layout/Icons';

interface RouteParams {
    imageKey: string;
}

interface Props extends RouteComponentProps<RouteParams> {
    sources?: SourceStore;
    images?: ImageStore;
}

@withRouter
@inject('sources', 'images')
@observer
export default class EditImageScreen extends React.Component<Props> {
    @observable imageCopy: Image;
    @observable originalImage: Image;
    @observable showSuccessMsg: boolean;

    componentWillMount() {
        this.loadImage(this.props.match.params.imageKey);
    }

    loadImage(key: string) {
        this.props.images.findByKey.run(key).finally(this.findImageComplete).done();
    }

    findImageComplete = () => {
        if (this.props.images.findByKey.err) {
            requestAnimationFrame(() => {
                this.props.history.push('/images');
            });
        }

        const image = this.props.images.findByKey.value;
        if (image) {
            this.originalImage = image;
            this.imageCopy = image.clone();
            this.loadReposAndOwners();
        } else {
            requestAnimationFrame(() => {
                this.props.history.push('/images');
            });
        }
    }

    loadReposAndOwners() {
        transaction(() => {
            this.props.sources.selectProvider(this.imageCopy.repository.provider);
            this.props.sources.selectOrganisation(this.imageCopy.repository.owner);
            this.props.sources.loadOrganisations.run().done();
            this.props.sources.loadRepositories.run().done();
        });
    }

    saveChanges = () => {
        this.props.images.updateImage.run(this.originalImage.key, this.imageCopy).finally(this.afterSave).done();
    }

    @action.bound
    afterSave = () => {
        if (this.props.images.updateImage.err) return;

        const newImg = this.props.images.updateImage.value;
        this.loadImage(newImg.key);
        this.showSuccessMsg = true;
    }

    render() {
        if (!this.imageCopy) {
            return null;
        }

        return (
            <Screen>
                <Header title={this.originalImage.name} subTitle="Edit image properties..." back />
                {this.props.images.updateImage.err &&
                    <Alert message="An error happened while trying to update image record, please try again later." />
                }
                {this.showSuccessMsg && !this.props.images.updateImage.inProgress &&
                    <Alert type="success" message="Image updated successfully" />
                }
                <div className="content">
                    <ImageEditor image={this.imageCopy} sourceRepositories={this.props.sources.currentRepositories} sourceOwners={this.props.sources.currentOrganisations}>
                        <div className="form-actions">
                            <Button text="Save Changes" icon={this.props.images.updateImage.inProgress ? Icons.loadingSpinner : null} size="big" onClick={this.saveChanges} />
                        </div>
                    </ImageEditor>
                </div>
            </Screen>
        );
    }
}
