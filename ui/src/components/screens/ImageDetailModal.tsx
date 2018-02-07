import * as React from 'react';
import { observer, inject } from 'mobx-react';
import { withRouter, RouteComponentProps } from 'react-router-dom';

import { Modal, ModalHeader, ModalContent, ModalActions } from '../layout/Modal';
import Button from '../layout/Button';
import DetailInfo from '../layout/DetailInfo';
import Icons from '../layout/Icons';

import RootStore from '../../state/RootStore';
import Image from '../../state/models/Image';

interface RouteParams {
    imageName: string;
}

interface Props extends RouteComponentProps<RouteParams> {
    store: RootStore
}

@withRouter
@inject('store')
@observer
export default class ImageDetailModal extends React.Component<Props> {
    close = () => {
        this.props.history.goBack();
    }

    showBuildHistory = () => {
        this.props.history.push(`/history/${this.props.match.params.imageName}`);
    }

    edit = () => {
        this.props.history.push(`/images/${this.props.match.params.imageName}/edit`);
    }

    render() {
        const { store, match } = this.props;
        const image = store.images.findByName(this.props.match.params.imageName);
        if (!image) {
            this.props.history.push('/images');
            return null;
        }

        return (
            <Modal onClose={this.close}>
                <ModalHeader title={image.name} subTitle="Image Details" onClose={this.close} />
                <ModalContent>
                    <DetailInfo label="Source Provider:">{image.repository.provider}</DetailInfo>
                    <DetailInfo label="Repository:">{image.repository.owner}/{image.repository.name}</DetailInfo>
                    <DetailInfo label="Dockerfile Path:">{image.dockerfile_path}</DetailInfo>
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
                                {image.tags.map(tag =>
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
