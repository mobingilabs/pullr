import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { withRouter } from 'react-router-dom';

import Screen from '../layout/Screen';
import Header from '../layout/Header';
import Button from '../layout/Button';
import ImageEditor from "../widgets/ImageEditor";
import Notification from '../widgets/Notification';
import * as actions from '../../state/images/actions';

export class EditImageScreen extends React.Component {
    constructor(props) {
        super(props);
        this.state = { image: props.image, changed: false };
    }

    updateImage = (image) => {
        this.setState({ image });
    }

    cancel = () => {
        this.props.history.goBack();
    }

    saveChanges = () => {
        this.props.onSaveChanges(this.props.image.name, this.state.image);
    }

    render() {
        const { image, changed } = this.state;

        return (
            <Screen>
                <Header title={image.name} subTitle="Edit image properties..." back />
                <Notification id="images-update-image" />
                <div className="content">
                    <ImageEditor image={ image } onUpdateImage={this.updateImage}>
                        <div className="form-actions">
                            <Button text="Cancel" size="big" secondary onClick={this.cancel} />
                            <Button text="Save Changes" size="big" onClick={this.saveChanges} disabled={ this.state.image === this.props.image } />
                        </div>
                    </ImageEditor>
                </div>
            </Screen>
        );
    }
}

const mapStateToProps = (state, ownProps) => {
    const imageName = ownProps.match.params.imageName;
    return {
        image: state.images.data[imageName]
    };
};

const mapDispatchToProps = (dispatch) => ({
    onSaveChanges: (imageName, data) => dispatch(actions.saveChanges(imageName, data))
});

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(EditImageScreen));
