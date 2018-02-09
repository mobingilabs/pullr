import * as React from 'react';
import { observable } from 'mobx';
import { observer, inject } from 'mobx-react';
import { withRouter, RouteComponentProps } from 'react-router-dom';

import Image from '../../state/models/Image';
import ImageStore from '../../state/ImagesStore';

import { Modal, ModalHeader, ModalContent, ModalActions } from '../layout/Modal';
import Button from '../layout/Button';
import DetailInfo from '../layout/DetailInfo';
import Icons from '../layout/Icons';


interface RouteParams {
    imageKey: string;
}

interface Props extends RouteComponentProps<RouteParams> {
    images: ImageStore
}

@withRouter
@inject('images')
@observer
export default class ImageDetailModal extends React.Component<Props> {
    @observable image: Image;

    componentWillMount() {
        this.props.images.findByKey.run(this.props.match.params.imageKey).finally(this.afterImageFind).done();
    }

    close = () => {
        this.props.history.goBack();
    }

    afterImageFind = () => {
        if (this.props.images.findByKey.err) {
            console.warn(`Image not found by key: ${this.props.match.params.imageKey}`);
            this.props.history.replace('/images');
            return;
        }

        this.image = this.props.images.findByKey.value;
    }

    showBuildHistory = () => {
        this.props.history.push(`/history/${this.props.match.params.imageKey}`);
    }

    edit = () => {
        this.props.history.push(`/images/${this.props.match.params.imageKey}/edit`);
    }

    render() {
        if (!this.image) {
            return null;
        }

        return (
            <Modal onClose={this.close}>
                <ModalHeader title={this.image.name} subTitle="Image Details" onClose={this.close} />
                <ModalContent>
                    <DetailInfo label="Source Provider:">{this.image.repository.provider}</DetailInfo>
                    <DetailInfo label="Repository:">{this.image.repository.owner}/{this.image.repository.name}</DetailInfo>
                    <DetailInfo label="Dockerfile Path:">{this.image.dockerfile_path}</DetailInfo>
                    <DetailInfo label="Builds:">
                        <table className="table-inline">
                            <thead>
                                <tr>
                                    <th>Build By</th>
                                    <th>Matcher</th>
                                    <th>Docker Tag</th>
                                </tr>
                            </thead>
                            <tbody>
                                {this.image.tags.map(tag =>
                                    <tr key={tag.ref_test || tag.name}>
                                        <td>{tag.ref_type}</td>
                                        <td>{tag.ref_test}</td>
                                        <td>{tag.name || 'Same as git tag'}</td>
                                    </tr>
                                )}
                            </tbody>
                        </table>
                    </DetailInfo>
                </ModalContent>
                <ModalActions>
                    <Button secondary text="Show History" icon={Icons.history} onClick={this.showBuildHistory} />
                    <Button text="Edit" icon={Icons.edit} onClick={this.edit} />
                </ModalActions>
            </Modal>
        );
    }
}
