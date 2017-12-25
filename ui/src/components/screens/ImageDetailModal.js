import React from 'react';
import PropTypes from 'prop-types';
import { withRouter } from 'react-router-dom';
import { connect } from 'react-redux';

import { Modal, ModalHeader, ModalContent, ModalActions } from '../layout/Modal';
import Button from '../layout/Button';
import DetailInfo from '../layout/DetailInfo';
import Icons from '../layout/Icons';

export class ImageDetailModal extends React.PureComponent {
    edit = () => {
        this.props.history.push(`/images/edit/${this.props.image.name}`);
    }

    showBuildHistory = () => {
        this.props.history.push(`/history/${this.props.image.name}`);
    }

    close = () => {
        this.props.history.goBack();
    }

    edit = () => {

    }

    showHistory = () => {

    }

    render() {
        const { image } = this.props;
        return (
            <Modal onClose={this.close}>
                <ModalHeader title={image.name} subTitle="Image Details" onClose={this.close} />
                <ModalContent>
                    <DetailInfo label="Source Provider:">{image.provider}</DetailInfo>
                    <DetailInfo label="Repository:">{image.organisation}/{image.repository}</DetailInfo>
                    <DetailInfo label="Dockerfile Path:">{image.dockerfilePath}</DetailInfo>
                    <DetailInfo label="Build Tags:">
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
                                    <tr key={tag.tag || tag.name}>
                                        <td>{tag.type}</td>
                                        <td>{tag.name}</td>
                                        <td>{tag.tag || 'Same as git tag'}</td>
                                    </tr>
                                )}
                            </tbody>
                        </table>
                    </DetailInfo>
                </ModalContent>
                <ModalActions>
                    <Button secondary text="Show History" icon={Icons.history} onClick={this.showHistory} />
                    <Button text="Edit" icon={Icons.edit} onClick={this.edit} />
                </ModalActions>
            </Modal>
        );
    }
}

function mapStateToProps(state, ownProps) {
    return {
        image: state.images.details[ownProps.match.params.imageName]
    };
}

export default withRouter(connect(mapStateToProps)(ImageDetailModal));
